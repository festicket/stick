# Stick

## Build

Run:

```bash
make build
```

It will create a binary called `stick` in the current directory. In order of it to work - need to create a config file:

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
