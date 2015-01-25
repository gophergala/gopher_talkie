package main

import (
	_ "fmt"
	"github.com/gophergala/gopher_talkie/src/cli/app"
)

func main() {
	app := app.NewApp()
	app.Run()
}
