package main

import "log"

// routes, handlers

// db connections

func main() {
	if err := StartServer(); err != nil {
		log.Fatal("Error during server start: ", err.Error())
	}
}
