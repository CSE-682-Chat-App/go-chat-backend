# Setup

Clone the repository to `$GOPATH/src`

### Example
```bash
> cd $GOPATH/src
> git clone git@github.com:CSE-682-Chat-App/go-chat-backend.git
```

# Bundling and Running

Build the docker container

### Example
```bash
> docker build -t go-chat-backend .
> docker run -it -p 9090:9090 go-chat-backend
INFO[0000] Starting Server on port 9090
```

You can then open your browser to: `http://localhost:9090` and you should see a default response.
