
# run make/*.go

build: webshellmake
	make/web-shell-make build

down: webshellmake
	make/web-shell-make down

winpty: webshellmake
	make/web-shell-make winpty

gen: webshellmake
	make/web-shell-make gen

debug: webshellmake
	make/web-shell-make debug

run: webshellmake
	make/web-shell-make run

clean: webshellmake
	make/web-shell-make clean

help: webshellmake
	make/web-shell-make help


webshellmake:
	cd make && go build
