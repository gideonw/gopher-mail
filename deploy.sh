#! /bin/bash

########################################################################
# Values
cd terraform

LAMBDA_ARCHIVE_BUCKET="$(terraform output lambda_archive_bucket)"
echo "[INFO] Uploading lambda archives to $LAMBDA_ARCHIVE_BUCKET"

cd ../

########################################################################
# Deploy

echo "[INFO] Uploading postmaster to $LAMBDA_ARCHIVE_BUCKET/postmaster/"

aws s3 cp ./bin/postmaster.zip s3://$LAMBDA_ARCHIVE_BUCKET/postmaster/postmaster.zip