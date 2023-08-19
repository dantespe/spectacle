TEST=go test -v -coverprofile

COVERS=\
	db/cover.out\
	operation/cover.out\
	dataset/cover.out\
	header/cover.out\
	record/cover.out\
	cell/cover.out\

DATABASES=\
	prod\
	dev\
	test\
	tmp\

PSQL_SCHEMA=db/schema.sql
CREATE_DATABASES=db/databases.sql
DOCKER_CONTAINER=postgres
DOCKER_SLEEP=5
POSTGRES_USER=postgres
# Store the password in the environment variable on the machine
POSTGRES_PASSWORD=$$PGPASSWORD
POSTGRES_DATABASE_NAME=dev
DOCKER_IMAGE=postgres:15.3

export SPECTACLE_DIR=${CURDIR}
export SPECTACLE_DATA_DIR=$$HOME/spectacle/data

.PHONY: postgres all

all: docker_start server
	
server:
	go run server.go 

postgres:
	mkdir -p ${SPECTACLE_DATA_DIR}
	docker run -d \
		-p 5432:5432\
		--name $(DOCKER_CONTAINER) \
		-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
		-v ${SPECTACLE_DATA_DIR}:/var/lib/postgresql/data \
		$(DOCKER_IMAGE) && sleep $(DOCKER_SLEEP)
	cat $(PSQL_SCHEMA) | docker exec -i $(DOCKER_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE_NAME)
	cat $(PSQL_SCHEMA) | docker exec -i $(DOCKER_CONTAINER) psql -U $(POSTGRES_USER) -d test

docker_start:
	for i in $(DATABASES); do \
		file=$$(echo "$$i"); \
		echo "CREATE DATABASE $$file" | docker exec -i $(DOCKER_CONTAINER) psql -U $(POSTGRES_USER); \
		cat $(PSQL_SCHEMA) | docker exec -i $(DOCKER_CONTAINER) psql -U $(POSTGRES_USER) -d $$file ; \
	done

docker_exec:
	docker exec -it $(DOCKER_CONTAINER) psql -U $(POSTGRES_USER) -d dev

docker_stop:
	docker stop $(DOCKER_CONTAINER)

docker_create:
	cat $(PSQL_SCHEMA) | docker exec -i $(DOCKER_CONTAINER) psql -U $(POSTGRES_USER) -d $(POSTGRES_DATABASE_NAME)

docker_clean: docker_stop
	rm -rf ${SPECTACLE_DATA_DIR}
	mkdir ${SPECTACLE_DATA_DIR}

test: docker_start db_test operation_test dataset_test header_test record_test cell_test

db_test:
	$(TEST) db/cover.out ./db

operation_test: operation/operation.*go
	$(TEST) operation/cover.out ./operation

dataset_test: dataset/dataset.*go
	$(TEST) dataset/cover.out ./dataset

header_test: header/header.*go
	$(TEST) header/cover.out ./header

record_test: record/record.*go
	$(TEST) record/cover.out ./record

cell_test: cell/cell.*go
	$(TEST) cell/cover.out ./cell

clean:
	rm $(COVERS)

format:
	go fmt ./...
