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

echo "[INFO] Uploading mailman to $LAMBDA_ARCHIVE_BUCKET/mailman/"
aws s3 cp ./bin/mailman.zip s3://$LAMBDA_ARCHIVE_BUCKET/mailman/mailman.zip

echo "[INFO] Uploading mailtruck to $LAMBDA_ARCHIVE_BUCKET/mailtruck/"
aws s3 cp ./bin/mailtruck.zip s3://$LAMBDA_ARCHIVE_BUCKET/mailtruck/mailtruck.zip