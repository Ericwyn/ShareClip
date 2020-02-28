echo 'build windows application'

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ShareClipClient src/main/ShareClipClient.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ShareClipServer src/main/ShareClipServer.go

echo 'build linux application'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ShareClipClient.exe src/main/ShareClipClient.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ShareClipServer.exe src/main/ShareClipServer.go