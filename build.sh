echo 'build linux application'
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ShareClip src/main/*.go

echo 'build windows application'
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ShareClip.exe src/main/*.go
