# Background
Endeavouring to scratch (pun intended) an itch, by running an AWS Lambda function, written in Go, in a scratch (so OS-less) container. This repo contains a couple of files 

> 💡 This has been created by an idiot who should rarely be trusted with anything sharper than a butter knife, proceed with caution...

# Files 💾
### main.go
Rather trivial Go [source code](./main.go "main.go") that renders some log messages when the Lambda function is invoked. To compile, the following might be reasonably assumed to behave as required:
```bash
go env -w GOOS=linux
go env -w GOARCH=amd64 
go mod download
go build main.go
```

### Dockerfile
The [`Dockerfile`](./Dockerfile "Dockerfile") is deliberately as short and simple as possible. Nuff said?

# Local Explorification 🧑‍💻
All this was done on macOS...

 #### 🛠 setup
```bash
export ACCOUNT={12 digit AWS account ID}
export REGION=eu-west-2
export ECR_REPO=go
export FUNCTION_NAME=demo
export FUNCTION_OUTPUT=/tmp/logs-command
aws ecr get-login-password --region ${REGION} | docker login --username AWS --password-stdin ${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com
```

#### 🤞 build
Assumes docker and [dive](https://github.com/wagoodman/dive "dive") are installed.
```bash
GOOS=linux GOARCH=amd64 go build -o main main.go && echo -n "Build success, provide an image tag version: " && read VERSION || echo "Build failed! 🧨"

docker build --tag ${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com/${ECR_REPO}:${VERSION} .
docker push ${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com/${ECR_REPO}:${VERSION}

dive ${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com/${ECR_REPO}:${VERSION}

aws lambda update-function-code --function-name ${FUNCTION_NAME} --image-uri ${ACCOUNT}.dkr.ecr.${REGION}.amazonaws.com/${ECR_REPO}:${VERSION} --no-cli-pager && aws lambda wait function-updated-v2 --function-name ${FUNCTION_NAME} && echo "Function ${FUNCTION_NAME} updated 👌" || echo "Failed to update ${FUNCTION_NAME} 😱"
```

#### 🏃 run
Assumes [jq](https://stedolan.github.io/jq/ "jq") is installed.
```bash
aws lambda invoke --function-name ${FUNCTION_NAME} --payload '{"wibble":"wobble","plop":["plip"],"true":false,"emoji":"🤓"}' --cli-binary-format raw-in-base64-out --no-cli-pager ${FUNCTION_OUTPUT} && eval $(sleep 3; cat ${FUNCTION_OUTPUT} | cut -d\" -f2) | jq '.events[].message' -r | sed -e '/^$/d'; rm ${FUNCTION_OUTPUT}
```

# ToDo
- create a multi-stage builder Dockerfile

# Possibly Handy Links
- https://docs.aws.amazon.com/lambda/latest/dg/lambda-golang.html
- https://github.com/aws/aws-sdk-go
- https://github.com/aws/aws-lambda-go
- https://docs.aws.amazon.com/lambda/latest/dg/go-image.html
- https://go.dev/doc/install/source#environment
- https://go.dev/learn/
