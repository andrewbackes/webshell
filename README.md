# Webshell
Shell interaction over the web. Users can run a command from a terminal in a web browser on the backend system running the API. It is meant to be added as an admin portal to a backend service. Of course you should be careful with it since you are allowing users to execute whatever they want over the web.

## How does it work?
Webshell uses websockets to communicate between the browser and the backend. It also uses GitHub OAuth2 for authentication and authorization and keeps track of user sessions.

Users can send commands through a web ui that looks similar to a terminal window. The commands run on the backend where the API lives.

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

## Auth
You can turn on OAuth2 with GitHub as the provider. To do so you will need to do three things:
1. [Add an OAuth app your GitHub account](https://developer.github.com/apps/building-oauth-apps/creating-an-oauth-app/).
2. Set the `AUTH_ENABLED=true` environment variable before running.
3. Create the file `./auth_config.json` using `./auth_config_schema.json` as a template.

### Example
An example of setting up auth:
```
$ cat <<EOF >auth_config.json
{
    "ClientID": "abc",
    "ClientSecret": "123",
    "RedirectURL": "http://localhost:8080",
    "AuthorizedUsers": ["andrewbackes"]
}
EOF
$ export AUTH_ENABLED=true
$ go run cmd/webshell.go
```