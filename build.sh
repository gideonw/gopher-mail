cd postmaster
env GOOS=linux go build -ldflags="-s -w" -o ../bin/postmaster main.go sns.go
cd ../

cd mailman
env GOOS=linux go build -ldflags="-s -w" -o ../bin/mailman main.go sns.go
cd ../

cd mailtruck
env GOOS=linux go build -ldflags="-s -w" -o ../bin/mailtruck main.go sns.go
cd ../

