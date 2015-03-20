package main

import (
	"github.com/kiasaki/coresize"
)

func main() {
	s := coresize.NewServer()
	s.ParseFlags()
	s.SetupAndRun()
}
