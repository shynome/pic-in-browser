build:
	CGO_ENABLED=0 go build -o pic-in-browser
docker: build
	docker build . -t shynome/pic-in-browser:$$(git describe --tags --always --dirty)
run: docker
	docker run --rm -ti -p 7070:7070 shynome/pic-in-browser:$$(git describe --tags --always --dirty)
