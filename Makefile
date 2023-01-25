NAME=citybike
DB=citybike.db

all: run

${NAME}:
	go build -o ${NAME} cmd/main.go

${DB}:
	./${NAME} -download

download: ${DB}

run: ${NAME} ${DB}
	./${NAME}

clean:
	rm ${NAME}

cleanall: 
	rm ${NAME} ${DB}

test:
	go test ./... -v

docker_build:
	docker build -t bengissimo/citybike .

docker_run: docker_build
	docker run --rm -p 8000:8000 bengissimo/citybike

docker: docker_build docker_run
