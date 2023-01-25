# specify the base image to  be used for the application, alpine or ubuntu
FROM golang:1.19 as builderImage

# create a working directory inside the image
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
# 
RUN CGO_ENABLED=1 go build --ldflags '-linkmode external -extldflags=-static' -o citybike cmd/main.go
RUN ./citybike -download

FROM scratch
WORKDIR /build

COPY --from=builderImage /build/citybike /build/citybike
COPY --from=builderImage /build/citybike.db /build/citybike.db

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8000

# command to be used to execute when the image is used to start a container
CMD ["./citybike"]