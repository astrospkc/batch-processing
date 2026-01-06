package main

import (
	"batch-processing/src/async"
	"batch-processing/src/connect"
	"fmt"
	"log"
)


func main(){
	fmt.Print("hello there")
	connect.InitRedisConnect()
	connect.MongoConnect()
	server := async.NewServer()
	mux:= async.NewMux()
	log.Println("🚀 Asynq worker started")
	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}

}


