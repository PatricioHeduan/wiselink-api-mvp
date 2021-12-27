package main

import (
	"log"
	"wiselink/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	serv, err := server.New("8080")
	if err != nil {
		log.Fatal(err)
	}
	//Ejecutamos de forma concurrente el servidor
	serv.Start()

}
