package main

import "github.com/sodibus/sodibus"

func main() {
	s := sodibus.NewServer("0.0.0.0:7788")
	err := s.Run()
	println("Run Error: ", err)
}
