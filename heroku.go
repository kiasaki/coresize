// +build heroku
package coresize

func main() {
	s := NewServer()
	s.ParseFlags()
	s.SetupAndRun()
}
