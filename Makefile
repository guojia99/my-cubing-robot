all: go

go:
	go run -v main.go robot

buildx:
	go build -o mycube-robot main.go


local:
	go run -v main.go robot --config etc/config_local.yml