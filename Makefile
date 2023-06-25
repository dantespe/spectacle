TEST=go test -v -coverprofile

COVERS=\
	db/cover.out\
	operation/cover.out\
	dataset/cover.out\
	header/cover.out\
	record/cover.out\
	cell/cover.out\

export SPECTACLE_DIR=${CURDIR}

.PHONY: sqlite all

all: sqlite
	go run server.go 

sqlite:
	sqlite3 ${SPECTACLE_DIR}/db/dev.db < ${SPECTACLE_DIR}/db/schema.sql
	sqlite3 ${SPECTACLE_DIR}/db/test.db < ${SPECTACLE_DIR}/db/schema.sql

test: db_test operation_test dataset_test header_test record_test

db_test: sqlite
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
	rm db/dev.db db/test.db $(COVERS)

format:
	go fmt ./...