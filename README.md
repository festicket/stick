# Stick

## Dev dependencies

This project requires Go to develop. All dependencies are are listed in the Brewfile, and can be installed using Homebrew:

```
brew bundle
```

## Build

You have two options to build the binary:

1. If you have Go installed locally and you checkouted the repository withint your $GOPATH then just use:
Run:

```bash
make build
```

2. You also can checkout the repo anywhere and build the binary using Docker:

```bash
make build_docker
```

By default `build_docker` command will build the binary for MacOS (amd64) but you can override this via `GOOS` and `GOARCH` paramaters e.g.

```bash
make build_docker GOOS=windows
```

* although currently it is tested only for MacOS.

Either way you will have a binary called `stick` in the current directory. 


In order of it to work - need to create a config file:

```bash
cp config.json.example config.json
```

Update `config.json` 

Where:

`command`: a command to run the server and make it ready to accept HTTP requests to the endpoint specified as `url`. For example you might want to checkout a specific branch, reset the cache and start the server. 

To do that, you could create a shell script to run your Django-like project:

```bash
cd /foo/bar/project/ &&
git reset HEAD --hard && git checkout feature &&
source venv/bin/activate &&
python manage.py runserver
```


`pre_request_command`: the same idea as for `command` but it will be executed before earch HTTP request. 
