NAME=citybike

all: build download run

build:
	go build -o ${NAME} cmd/main.go

download:
	./${NAME} -download

run:
	./${NAME}

clean: 
	rm ${NAME}

test:
	go test ./... -v

docker_build:
	docker build -t bengissimo/citybike .

docker_run:
	docker run --rm -p 8000:8000 bengissimo/citybike

docker: docker_build docker_run
