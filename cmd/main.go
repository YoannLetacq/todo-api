package main

import (
	"YoannLetacq/todo-api.git/config"
	"log"
)

func main() {
	// charge les variables d'environnement
	config.InitEnv()

	// Initialise la BDD (set testing to false)
	config.InitDB(false)

	log.Println("Initialisation Successfull !")
}
