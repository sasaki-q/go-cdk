build: 
	GOOS=linux GOARCH=amd64 go build -o handler/bin/$(target)/$(target) handler/$(target)/$(target).go