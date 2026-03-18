# To do list example

Shows a fully loaded server that by default listens on a unix, http and https socket.

You can run just the http listener for quick testing:

```shellsession
go run ./cmd/todo-list-server/main.go --scheme http
```

## Run full server

To run the full server you need to generate TLS certificates, build the binary and run it with sudo enabled.

### Generate TLS certificates

From the repository root:

```shellsession
cd hack/tools
go run . gen-certs
```

This generates a self-signed CA, server certificate (`mycert1.crt` / `mycert1.key`), and client
certificate (`myclient.crt` / `myclient.key`) in the `todo-list-errors/` directory.

### Start the server

```shellsession
go build ./cmd/todo-list-server
sudo ./todo-list-server --tls-certificate mycert1.crt --tls-key mycert1.key
```

## Generate code

```shellsession
swagger generate server -A TodoList -f ./swagger.yml -e
```
