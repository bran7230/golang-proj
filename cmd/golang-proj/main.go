package main

import "log"

func main() {
	if err := StartServer(); err != nil {
		log.Fatal("Error during server start: ", err.Error())
	}
}
