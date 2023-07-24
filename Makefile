run: build
	./sidb
build:
	gcc sidb.c -o sidb
clean:
	rm -f sidb
