#!/bin/bash

project=crawlerapi
env=arc
AWS_REGION=ap-southeast-1

for functionName in $(aws lambda list-functions --region $AWS_REGION --query 'Functions[?starts_with(FunctionName, `'$env'-'$project'-`) == `true`].FunctionName' --output json) ; 
do 
    if [ "$functionName" != "[" ] && [ "$functionName" != "]" ]
    then
        name=$(echo $functionName | sed 's/"//g' | sed 's/,//g') 
        echo "getting config for $name"
        aws lambda get-function-configuration --function-name $name > configuration-$name.json
    fi
done