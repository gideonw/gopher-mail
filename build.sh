#! /usr/bin/env /bin/bash

echo "[INFO] Building postmaster..."

cd postmaster
env GOOS=linux go build -ldflags="-s -w" -o ../bin/postmaster main.go sns.go

echo "[INFO] Building mailman..."
cd ../mailman
env GOOS=linux go build -ldflags="-s -w" -o ../bin/mailman main.go

echo "[INFO] Building mailtruck..."
cd ../mailtruck
env GOOS=linux go build -ldflags="-s -w" -o ../bin/mailtruck main.go
cd ../

echo "[INFO] Archiving postmaster..."
zip ./bin/postmaster.zip ./bin/postmaster
echo "[INFO] Archiving mailman..."
zip ./bin/mailman.zip ./bin/mailman 
echo "[INFO] Archiving mailtruck..."
zip ./bin/mailtruck.zip ./bin/mailtruck 
