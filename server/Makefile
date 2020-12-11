cc := 

server_build:
	go build -o ./runServer server.go crypto.go

server_run:
	go run server.go crypto.go

server_clean:
	[ -f ./priv_key ] && rm ./priv_key
	[ -f ./pub_key ] && rm ./pub_key

server_clean_run:
	server_clean
	server_run
