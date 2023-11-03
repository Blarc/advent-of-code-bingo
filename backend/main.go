package main

import "github.com/gin-gonic/gin"

func main() {
	app := App{
		router: gin.Default(),
	}
	app.start()
}
