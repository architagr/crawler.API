accountId := $(shell aws sts get-caller-identity --query Account --output text --profile default)
region := $(AWS_REGION)
project := crawlerapi
env := $(ENV)
ifndef region
	region = ap-southeast-1
endif

ifndef env
	env = dev
endif

build:
	make -C ../JobAPI lambda-build
	make -C ../LoginAPI lambda-build
clear:
	make -C ../JobAPI clear
	make -C ../LoginAPI clear
	
bootstrap:
	cdk bootstrap

deploy:
	make build
	rm -rf cdk.out
	clear
	cdk deploy --all --require-approval never \
	--context ACCOUNT_ID=$(accountId) \
	--context REGION=$(region) \
	--context PROJECT=$(project) \
	--context ENV=$(env)
	make clear

destroy:
	cdk destroy --all --require-approval never