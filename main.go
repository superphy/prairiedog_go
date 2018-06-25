package main

import "github.com/superphy/prairiedog/cmd"

func main() {
	app := cmd.New()
	app.RunAndExitOnError()
}
