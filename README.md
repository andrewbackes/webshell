# Webshell
Shell interaction over the web. Users can run a command from a terminal in a web browser on the backend system running the API. It is meant to be added as an admin portal to a backend service. Of course you should be careful with it since you are allowing users to execute whatever they want over the web.

## How does it work?
There are two pieces. The frontend webapp is done with React and found in `webapp/`. The backend app is written in Go and can be found in `cmd/`.

## How do I run it?
### Directly

Make sure you have yarn and go installed.

Then you need to build the webapp by:
```
$ cd webapp
# yarn install
$ yarn build
```
Finally, you need to run the service:
```
$ cd ...
$ go run cmd/webshell.go
```

### With Docker

```
docker build -t webshell .
docker run --rm -p 8080:8080 webshell
```