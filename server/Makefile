# First time working with a Makefile, This will be the last focus on the project.

server_build:
	go build -o ./runServer server.go crypto.go auth.go

server_run:
	go run server.go crypto.go auth.go

server_clean:
	[ -f ./priv_key ] && rm ./priv_key
	[ -f ./pub_key ] && rm ./pub_key

server_clean_run:
	[ -f ./priv_key ] && rm ./priv_key
	[ -f ./pub_key ] && rm ./pub_key
	go run server.go crypto.go	auth.go
