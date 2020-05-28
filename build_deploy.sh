#! /usr/bin/env /bin/bash
./build.sh
./deploy.sh
cd terraform
terraform apply -target aws_lambda_function.postmaster -target aws_lambda_function.mailman --auto-approve
cd ../