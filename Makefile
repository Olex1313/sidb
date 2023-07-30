run: test
	./sidb
build:
	gcc sidb.c -o sidb
clean:
	rm -f sidb
test: build
	go test sidb_test.go
