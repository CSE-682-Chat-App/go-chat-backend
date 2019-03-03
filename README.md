# Install and Run

### Checkout
Check the file out to the directory: `$GOPATH/src/github.com/CSE-682-Chat-App/go-chat-backend`

### Install
```
go install ./...
```

### Run
```
go-chat-backend
```


# Container


### Build
```
docker build -t go-chat-backend .
```

### Run

```
docker run -p 8080:8080 go-chat-backend
```
