// +build heroku
package coresize

func main() {
	s := coresize.NewServer()
	s.ParseFlags()
	s.SetupAndRun()
}
