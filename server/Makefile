# First time working with a Makefile, This will be the last focus on the project.

server_build:
	go build -o ./runServer server.go crypto.go auth.go routeHandler.go
    docker pull redis

server_run:
	go run server.go crypto.go auth.go routeHandler.go
    docker run --name redis -d redis

server_clean:
	[ -f ./priv_key ] && rm ./priv_key
	[ -f ./pub_key ] && rm ./pub_key
    docker stop redis

server_clean_run:
	[ -f ./priv_key ] && rm ./priv_key
	[ -f ./pub_key ] && rm ./pub_key
	go run server.go crypto.go	auth.go routeHandler.go
