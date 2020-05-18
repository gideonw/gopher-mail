#! /bin/bash

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

