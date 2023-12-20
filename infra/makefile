accountId := $(shell aws sts get-caller-identity --query Account --output text --profile default)
region := $(AWS_REGION)
project := crawlerapi
env := $(ENV)
dbConnectionString := $(DbConnectionString)
databaseName := $(DatabaseName)
hostedZoneId := $(HostedZoneId)
certificateArn := $(CertificateArn)
ifndef region
	region = ap-southeast-1
endif

ifndef env
	env = dev
endif

build:
	make -C ../JobAPI lambda-build
	make -C ../LoginAPI lambda-build
	make -C ../UserAPI lambda-build
	make -C ../EmployerAPI lambda-build

clear:
	make -C ../JobAPI clear
	make -C ../LoginAPI clear
	make -C ../UserAPI clear
	make -C ../EmployerAPI clear

bootstrap:
	cdk bootstrap

deploy:
	make build
	rm -rf cdk.out
	clear
	cdk deploy --all --require-approval never \
	--context ACCOUNT_ID="$(accountId)" \
	--context REGION="$(region)" \
	--context PROJECT="$(project)" \
	--context ENV="$(env)" \
	--context DbConnectionString="$(dbConnectionString)" \
	--context DatabaseName="$(databaseName)" \
	--context HostedZoneId="$(hostedZoneId)" \
	--context CertificateArn="$(certificateArn)"
	make clear

destroy:
	cdk destroy --all --require-approval never

set-env:
	./get-env.sh
	go run get-env.go

tidy:
	make -C ../JobAPI tidy
	make -C ../LoginAPI tidy
	make -C ../UserAPI tidy
	make -C ../EmployerAPI tidy
	go mod tidy