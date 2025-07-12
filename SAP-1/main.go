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
	Halted     bool
}

type Instruction struct {
	Opcode        uint8
	MemoryAddress uint8
}

func splitInstruction(Input uint8) Instruction {
	instruction := Instruction{}
	instruction.Opcode = (Input & 0xF0) >> 4
	instruction.MemoryAddress = Input & 0x0F
	return instruction
}

func step(c *CPU) {
	switch c.clockCycle {
	case 0:
		c.MAR = c.PC
	case 1:
		c.PC++
	case 2:
		c.IR = c.memory[c.MAR]
	case 3, 4, 5:
		instruction := splitInstruction(c.IR)
		execute(c, &instruction)
	}
	c.clockCycle = (c.clockCycle + 1) % 6

}

func execute(c *CPU, i *Instruction) {
	if c.Halted {
		return
	}

	switch i.Opcode {
	case 0x0: //LDA
		c.A = c.memory[i.MemoryAddress]
	case 0x1: //ADD
		result := uint16(c.A) + uint16(c.memory[i.MemoryAddress])
		if result > 255 {
			c.C = true
		} else {
			c.C = false
		}
		c.A = uint8(result)
		if c.A == 0 {
			c.Z = true
		} else {
			c.Z = false
		}
	case 0x2: //SUB
		if c.A < c.memory[i.MemoryAddress] {
			c.C = true
		} else {
			c.C = false
		}
		c.A -= c.memory[i.MemoryAddress]
		if c.A == 0 {
			c.Z = true
		} else {
			c.Z = false
		}
	case 0x3: //STA
		c.memory[i.MemoryAddress] = c.A
	case 0x4: // LDI
		c.A = i.MemoryAddress
	case 0x5: // JMP
		c.PC = i.MemoryAddress
	case 0x6: // Jump if Zero
		if c.Z {
			c.PC = i.MemoryAddress
		}
	case 0x7: // Jump if carry
		if c.C {
			c.PC = i.MemoryAddress
		}
	case 0x8: // OUT
		c.OUT = c.A
	case 0x9: //HALT
		c.Halted = true

	}
}
func reset(c *CPU) {
	c.A = 0
	c.B = 0
	c.PC = 0
	c.MAR = 0
	c.IR = 0
	c.OUT = 0
	c.Z = false
	c.C = false
	c.clockCycle = 0

}

func run(c *CPU, maxCycles int) {
	for i := 0; i < maxCycles; i++ {
		step(c)
	}
}
func printCPUState(c *CPU) {
	fmt.Printf("A:%02X B:%02X PC:%02X OUT:%02X Z:%t C:%t Halted:%t\n",
		c.A, c.B, c.PC, c.OUT, c.Z, c.C, c.Halted)
}

func loadROM(c *CPU, rom []uint8) {
	copy(c.memory[:], rom)
}
func rom1_hello() []uint8 {
	return []uint8{
		0x42, // LDI 2 (load immediate value 2 into A)
		0x80, // OUT (output A to OUT register)
		0x90, // HALT
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}
func main() {
	fmt.Printf("=========SAP-1 EMULATOR =========")
	cpu := &CPU{}

	reset(cpu)
	loadROM(cpu, rom1_hello())
	run(cpu, 100)
	printCPUState(cpu)
}
