#! /usr/bin/env /bin/bash

echo "[INFO] Building postmaster..."

cd postmaster
rm ../bin/postmaster 2> /dev/null
env GOOS=linux GOARCH=amd64 go build -o ../bin/postmaster main.go sns.go email.go

echo "[INFO] Building mailman..."
cd ../mailman
rm ../bin/mailman 2> /dev/null
env GOOS=linux GOARCH=amd64 go build -o ../bin/mailman main.go s3.go auth.go

echo "[INFO] Building mailtruck..."
cd ../mailtruck
rm ../bin/mailtruck 2> /dev/null
env GOOS=linux GOARCH=amd64 go build -o ../bin/mailtruck main.go
cd ../

CMD=zip
if [[ "$OSTYPE" == "msys" ]]; then
    command -v build-lambda-zip >/dev/null 2>&1 || { echo >&2 "Error: build-lambda-zip is required to ensure the zip is properly formated. https://github.com/aws/aws-lambda-go#for-developers-on-windows"; exit 1; }
    CMD="build-lambda-zip --output"
fi

echo "[INFO] Using \"$CMD\" to create lambda deployment zip."

cd ./bin
echo "[INFO] Archiving postmaster..."
rm postmaster.zip 2> /dev/null
chmod +x postmaster 
$CMD postmaster.zip postmaster 

echo "[INFO] Archiving mailman..."
rm mailman.zip 2> /dev/null
chmod +x mailman 
$CMD mailman.zip mailman 

echo "[INFO] Archiving mailtruck..."
rm mailtruck.zip
chmod +x mailtruck 2> /dev/null
$CMD mailtruck.zip mailtruck 
cd ../
