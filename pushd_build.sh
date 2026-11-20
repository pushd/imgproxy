#!/bin/bash -xe
# This script is for building and pushing to pushd ECR docker repo.

# # get aws account id from sts if not set
if [[ "${AWS_ACCOUNT_ID}x" == "x" ]]; then
  AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
fi

docker buildx build --platform linux/amd64 -t $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/pushd/imgproxy:$(git rev-parse HEAD)  -f docker/Dockerfile .
docker push $AWS_ACCOUNT_ID.dkr.ecr.us-east-1.amazonaws.com/pushd/imgproxy:$(git rev-parse HEAD)
