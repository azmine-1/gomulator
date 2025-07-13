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
		if c.PC >= 16 {
			c.PC = 0
		}
	case 2:
		if c.MAR >= 16 {
			c.MAR = 0
		}
		c.IR = c.memory[c.MAR]
	case 3, 4, 5:
		instruction := splitInstruction(c.IR)
		switch instruction.Opcode {
		case 0x0: // LDA
			if c.clockCycle == 3 {
				c.MAR = instruction.MemoryAddress
			} else if c.clockCycle == 4 {
				c.A = c.memory[c.MAR]
				c.Z = (c.A == 0)
			}

		case 0x1: // ADD
			if c.clockCycle == 3 {
				c.MAR = instruction.MemoryAddress
			} else if c.clockCycle == 4 {
				c.B = c.memory[c.MAR]
			} else if c.clockCycle == 5 {
				result := uint16(c.A) + uint16(c.B)
				c.C = (result > 255)
				c.A = uint8(result)
				c.Z = (c.A == 0)
			}

		case 0x2: // SUB
			if c.clockCycle == 3 {
				c.MAR = instruction.MemoryAddress
			} else if c.clockCycle == 4 {
				c.B = c.memory[c.MAR]
			} else if c.clockCycle == 5 {
				c.C = (c.A < c.B)
				c.A = c.A - c.B
				c.Z = (c.A == 0)
			}

		case 0x3: // STA
			if c.clockCycle == 3 {
				c.MAR = instruction.MemoryAddress
			} else if c.clockCycle == 4 {
				c.memory[c.MAR] = c.A
			}

		case 0x4: // LDI
			if c.clockCycle == 3 {
				c.A = instruction.MemoryAddress
				c.Z = (c.A == 0)
			}

		case 0x5: // JMP
			if c.clockCycle == 3 {
				c.PC = instruction.MemoryAddress
			}

		case 0x6: // JZ (Jump if Zero)
			if c.clockCycle == 3 && c.Z {
				c.PC = instruction.MemoryAddress
			}

		case 0x7: // JC (Jump if Carry)
			if c.clockCycle == 3 && c.C {
				c.PC = instruction.MemoryAddress
			}

		case 0x8: // OUT
			if c.clockCycle == 3 {
				c.OUT = c.A
			}

		case 0x9: // HALT
			if c.clockCycle == 3 {
				c.Halted = true
			}
		}
	}
	c.clockCycle = (c.clockCycle + 1) % 6

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
	for i := 0; i < maxCycles && !c.Halted; i++ {
		step(c)
	}
}

func printCPUStateDetailed(c *CPU) {
	fmt.Printf("Cycle:%d A:%02X B:%02X PC:%02X MAR:%02X IR:%02X OUT:%02X Z:%t C:%t Halted:%t\n",
		c.clockCycle, c.A, c.B, c.PC, c.MAR, c.IR, c.OUT, c.Z, c.C, c.Halted)
}

func loadROM(c *CPU, rom []uint8) {
	copy(c.memory[:], rom)
}
func stepDebug(c *CPU) {
	printCPUStateDetailed(c)
	step(c)
}
func rom1_hello() []uint8 {
	return []uint8{
		0x42, // LDI 2 (load immediate value 2 into A)
		0x80, // OUT (output A to OUT register)
		0x90, // HALT
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}
}
func rom_add_5_and_3() []uint8 {
	return []uint8{
		0x45, // LDI 5 (load immediate value 5 into A)
		0x1F, // ADD F (add value at memory address F to A)
		0x80, // OUT (output A to OUT register)
		0x90, // HALT
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x03, // Memory address F (15) contains the value 3
	}
}

func main() {
	fmt.Println("========= SAP-1 EMULATOR =========")
	cpu := &CPU{}
	// Second program - should print 8
	reset(cpu)
	loadROM(cpu, rom_add_5_and_3())
	fmt.Println("Running second program...")
	run(cpu, 100000000000)
	fmt.Printf("Program output: %d\n", cpu.OUT)
}
