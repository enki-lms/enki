package main

import "github.com/enki/daemon/internal/server"

func main() {
	err := server.StartServer()
	if err != nil {
		panic(err)
	}

}
