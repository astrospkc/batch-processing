package main

import (
	"batch-processing/src/async"
	"fmt"
	"log"
)


func main(){
	fmt.Print("hello there")
	server := async.NewServer()
	mux:= async.NewMux()
	log.Println("🚀 Asynq worker started")
	if err := server.Run(mux); err != nil {
		log.Fatal(err)
	}

}


