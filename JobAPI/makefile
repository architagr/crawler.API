lambda-build:
	make clear
	GOOS=linux GOARCH=amd64 go build -o JobAPI main.go 
	zip main.zip JobAPI
	rm -f JobAPI

clear:
	rm -f JobAPI main.zip
