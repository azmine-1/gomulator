package main

import "fmt"

type CPU struct {
	A          uint8
	B          uint8
	PC         uint8
	MAR        uint8
	IR         uint8
	OUT        uint8
	Z          bool
	clockCycle int
	memory     [16]uint8
}

type Instruction struct {
	Opcode       uint8
	MemoryAdress uint8
}

func splitInstruction(RI uint8) Instruction {
	instruction := Instruction{}
	instruction.Opcode = (RI & 0xF0) >> 4
	instruction.MemoryAdress = RI & 0x0F
	return instruction
}

func step(c *CPU, i *Instruction) {
	switch c.clockCycle {
	case 0:
		c.MAR = c.PC
	case 1:
		c.PC++
	case 2:
		c.IR = c.memory[c.MAR]
	case 3:
		execute(c, i)
	}

}

func execute(c *CPU, i *Instruction) {
	return (nil)
}

func main() {
	fmt.Printf("Hello world")
}
