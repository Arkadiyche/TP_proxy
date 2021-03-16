package main

import (
	"github.com/Arkadiyche/TP_proxy/database"
	"github.com/Arkadiyche/TP_proxy/models"
	server "github.com/Arkadiyche/TP_proxy/server"
	"log"
)

func main() {
	db := database.InitDatabase()
	server := server.NewServer(models.ServerConfig.Port, db)
	defer db.Close()
	//fmt.Println(models.Params)
	log.Fatal(server.ListenAndServe())
}
