package main

import "time"
import "math/rand"
import "github.com/sodibus/sodibus"

func main() {
	rand.Seed(time.Now().UnixNano())
	s := sodibus.NewNode("0.0.0.0:7788")
	err := s.Run()
	println("Run Error: ", err)
}
