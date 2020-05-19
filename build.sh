#! /usr/bin/env /bin/bash

echo "[INFO] Building postmaster..."

cd postmaster
env GOOS=linux GOARCH=amd64 go build -o ../bin/postmaster main.go sns.go

echo "[INFO] Building mailman..."
cd ../mailman
env GOOS=linux GOARCH=amd64 go build -o ../bin/mailman main.go

echo "[INFO] Building mailtruck..."
cd ../mailtruck
env GOOS=linux GOARCH=amd64 go build -o ../bin/mailtruck main.go
cd ../

# TODO: On windows require build-lambda-zip be installed.

cd ./bin
echo "[INFO] Archiving postmaster..."
rm postmaster.zip
chmod +x postmaster 
# zip -j postmaster.zip postmaster 
build-lambda-zip --output postmaster.zip postmaster 

echo "[INFO] Archiving mailman..."
rm mailman.zip
chmod +x mailman 
# zip -j mailman.zip mailman 
build-lambda-zip --output mailman.zip mailman 

echo "[INFO] Archiving mailtruck..."
rm mailtruck.zip
chmod +x mailtruck
# zip -j mailtruck.zip mailtruck 
build-lambda-zip --output mailtruck.zip mailtruck 
cd ../
