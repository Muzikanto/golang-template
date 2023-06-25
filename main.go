package main

import "go-backend-template/src"

func main() {
	var app = src.CreateApp()

	app.Init()
	app.Start(app.Config.HTTP().Address())
}
