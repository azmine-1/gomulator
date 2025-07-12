package main

import "fmt"

type CPU struct {
	A   uint8
	B   uint8
	PC  uint8
	MAR uint8
	IR  uint8
	OUT uint8
}

// struct I represents instructions
// named I for convenience
type I struct {
	O uint8 // opcode
	M uint8 // memory address
}

func main() {
	fmt.Printf("Hello world")
}
