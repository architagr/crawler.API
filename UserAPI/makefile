lambda-build:
	make clear
	GOOS=linux GOARCH=amd64 go build -o UserAPI main.go 
	zip main.zip UserAPI
	rm -f UserAPI

clear:
	rm -f UserAPI main.zip
