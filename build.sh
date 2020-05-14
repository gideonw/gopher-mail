cd mail-ingress
env GOOS=linux go build -ldflags="-s -w" -o ../bin/mail-ingress main.go sns.go
cd ../
