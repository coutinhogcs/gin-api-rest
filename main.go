package main

import (
	"github.com/coutinhogcs/api-go-gin/database"
	"github.com/coutinhogcs/api-go-gin/routes"
)

func main() {
	database.DbConection()
	routes.HandleRequest()
}
