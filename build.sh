CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o build/ares .
docker build . -f Dockerfile -t ares
