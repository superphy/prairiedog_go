package main

import (
	"runtime"

	"github.com/superphy/prairiedog/cmd"
)

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	cmd.Execute()

}
