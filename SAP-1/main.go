package main

import "fmt"

type CPU struct {
	A          uint8     //Accumulator
	B          uint8     // B regiser
	PC         uint8     // Program counter
	MAR        uint8     // Memory address regiser
	IR         uint8     // Instruction regiser
	OUT        uint8     // output register
	Z          bool      // zero flag
	clockCycle int       // cpu cycle
	memory     [16]uint8 // memory address bus
	C          bool
}

type Instruction struct {
	Opcode       uint8
	MemoryAdress uint8
}

func splitInstruction(Input uint8) Instruction {
	instruction := Instruction{}
	instruction.Opcode = (Input & 0xF0) >> 4
	instruction.MemoryAdress = Input & 0x0F
	return instruction
}

func step(c *CPU, input uint8) {
	switch c.clockCycle {
	case 0:
		c.MAR = c.PC
	case 1:
		c.PC++
	case 2:
		c.IR = c.memory[c.MAR]
	case 3, 4, 5:
		execute(c, input)
	}
	c.clockCycle = (c.clockCycle + 1) % 6

}

func execute(c *CPU, i *Instruction) {

	switch i.Opcode {
	case 0x0: //LDA
		c.memory[i.MemoryAdress] = c.A
	case 0x1: //ADD
		c.A += c.memory[i.MemoryAdress]
	case 0x2: //SUB
		c.A -= c.memory[i.MemoryAdress]
	case 0x3: //STA
		c.memory[i.MemoryAdress] = c.A
	case 0x4: // LDI
		c.A = i.MemoryAdress
	case 0x5: // JMP
		c.PC = i.MemoryAdress
	case 0x6: // Jump if Zero
		if c.Z {
			c.PC = i.MemoryAdress
		}
	case 0x7: // Jump if carry
		if c.C {
			c.PC = i.MemoryAdress
		}
	case 0x8: // Jump if zero
		c.OUT = c.PC
	case 0x9: // OUT
		return

	}
}

func main() {
	fmt.Printf("Hello world")
}
