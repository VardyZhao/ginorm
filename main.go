package main

import (
	"ginorm/app"
)

func main() {
	r := app.Init()
	app.Run(r)
}
