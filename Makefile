build:
	go build -o pic-in-browser
docker: build
	docker build . -t shynome/pic-in-browser
run: docker
	docker run --rm -ti -p 7070:7070 shynome/pic-in-browser
