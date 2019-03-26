help:
	@echo "build - build from sources"
	@echo "run - build and run"


build:
	go build -o stick *.go


run: build
	./stick
