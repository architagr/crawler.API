lambda-build:
	make clear
	GOOS=linux GOARCH=amd64 go build -o LoginAPI main.go 
	zip main.zip LoginAPI
	rm -f LoginAPI

clear:
	rm -f LoginAPI main.zip
