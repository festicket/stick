CURRENT_DIR = $(shell pwd)

# When build is done inside of Docker - we need to compile the binary
# for the host machine. User can specify approprite flags. 
# Default is MacOS
# https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
ifeq ($(GOOS),)
	GOOS := "darwin"
endif

ifeq ($(GOARCH),)
	GOARCH := "amd64"
endif


help:
	@echo "build - build from sources"
	@echo "build_docker - build using Docker"
	@echo "run - build and run"


build:
	go build -o stick *.go


build_docker:
	docker build . -t stick:latest
	docker run -v $(CURRENT_DIR):/src/ -e GOOS=$(GOOS) -e GOARCH=$(GOARCH) stick:latest make build

run: build
	./stick
