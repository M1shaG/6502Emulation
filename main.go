package main

import "fmt"

// TODO:
// Methods for bit manipulation
// WriteWord function

type word uint16 // handy name for 16bit

// 64 KB
const MAX_MEM = 1024 * 64

type Mem struct {
	Data [MAX_MEM]byte
}

// 6502 Documentation
// https://6502.org/users/obelisk/6502/

// Manipulate with bits
//
// cpu.Status |= FlagX -> Set X status to 1
// cpu.Status &^= FlagX -> Set X status to 0
// cpu.Status ^= FlagX -> test whether X is set
// cpu.Status ^= FlagX -> toggle X, e.g X was 0 became 1, was 1 became 0
const (
	FlagC byte = 1 << iota
	FlagZ
	FlagI
	FlagD
	FlagB
	_
	FlagV
	FlagN
)

type CPU struct {
	PC word // Program Counter
	SP byte // Stack Pointer

	A, X, Y byte //  Accumulator, Index Register X, Index Register Y

	// Processor Status
	// #7 #6 #5 #4 #3 #2 #1 #0 - List of bits
	// #0 (C) - Carry Flag
	// #1 (Z) - Zero Flag
	// #2 (I) - Interrupt Disable Flag
	// #3 (D) - Decimal Mode Flag
	// #4 (B) - Break Flag
	// #5 - unused
	// #6 (V) - Overflow Flag
	// #7 (N) - Negative Flag
	// https://www.middle-engine.com/images/2020-06-23-programming-the-nes-the-6502-in-detail/processor-status-register-2x.png
	Status byte
}

type Memory interface {
	Init()
}

// Fill memory with zeros
func (memory *Mem) Init() {
	for i := 0; i < MAX_MEM; i++ {
		memory.Data[i] = 0
	}
}

type Processor interface {
	Reset(memory *Mem)
	FetchByte(Cycles *uint32, memory *Mem) byte
	ReadByteFromByte(Cycles *uint32, Address byte, memory *Mem) byte
	ReadByteFromWord(Cycles *uint32, Address word, memory *Mem) byte
	ReadWord(Cycles *uint32, Address word, memory *Mem) byte
	FetchWord(Cycles *uint32, memory *Mem) word
	Execute(memory *Mem)
}

// Reset CPU (Initialize and start program execution)
func (c *CPU) Reset(memory *Mem) {
	c.PC = 0xFFFC
	c.SP = 0xFD
	c.Status &^= FlagC
	c.Status &^= FlagZ
	c.Status &^= FlagI
	c.Status &^= FlagD
	c.Status &^= FlagB
	c.Status &^= FlagV
	c.Status &^= FlagN
	c.A, c.X, c.Y = 0, 0, 0
	memory.Init()
}

// Reading Functions
func (c *CPU) FetchByte(Cycles *int32, memory *Mem) byte {
	Data := memory.Data[c.PC]
	c.PC++
	*Cycles--
	return Data
}

func (c *CPU) ReadByteFromByte(Cycles *int32, Address byte, memory *Mem) byte {
	Data := memory.Data[Address]
	*Cycles--
	return Data
}

func (c *CPU) ReadByteFromWord(Cycles *int32, Address word, memory *Mem) byte {
	Data := memory.Data[Address]
	*Cycles--
	return Data
}

func (c *CPU) ReadWord(Cycles *int32, Address word, memory *Mem) word {
	LoByte := c.ReadByteFromWord(Cycles, Address, memory)
	HiByte := c.ReadByteFromWord(Cycles, Address+1, memory)
	return word(LoByte) | (word(HiByte) << 8)
}

func (c *CPU) FetchWord(Cycles *int32, memory *Mem) word {
	var Data word = word(memory.Data[c.PC])
	c.PC++

	Data |= (word(memory.Data[c.PC]) << 8)
	c.PC++

	*Cycles -= 2

	// TODO: handle big endian
	return Data

}

func (c *CPU) LDASetStatus() {
	if c.A == 0 {
		c.Status |= FlagZ
	}

	if (c.A & 0b10000000) > 0 {
		c.Status |= FlagN
	}
}

// opcodes
// INS_XXX_AM
// INS - Instuction
// XXX - Name of Instruction
// AM - Addressing Mode
const (
	// LDA
	INS_LDA_IM   byte = 0xA9
	INS_LDA_ZP   byte = 0xA5
	INS_LDA_ZPX  byte = 0xB5
	INS_LDA_ABS  byte = 0xAD
	INS_LDA_ABSX byte = 0xBD
	INS_LDA_ABSY byte = 0xB9
	INS_LDA_INDX byte = 0xA1
	INS_LDA_INDY byte = 0xB1
	// lDX
	INS_LDX_IM   byte = 0xA2
	INS_LDX_ZP   byte = 0xA6
	INS_LDX_ZPY  byte = 0xB6
	INS_LDX_ABS  byte = 0xAE
	INS_LDX_ABSY byte = 0xBE
	// LDY
	INS_LDY_IM   byte = 0xA0
	INS_LDY_ZP   byte = 0xA4
	INS_LDY_ZPX  byte = 0xB4
	INS_LDY_ABS  byte = 0xAC
	INS_LDY_ABSX byte = 0xBC

	INS_JSR byte = 0x20
)

func (c *CPU) Execute(Cycles int32, memory *Mem) int32 {

	CyclesRequested := Cycles

executeLoop:
	for Cycles > 0 {
		Ins := c.FetchByte(&Cycles, memory)

		switch Ins {
		case INS_LDA_IM:
			Value := c.FetchByte(&Cycles, memory)
			c.A = Value
			c.LDASetStatus()
		case INS_LDA_ZP:
			ZeroPageAddress := c.FetchByte(&Cycles, memory)
			c.A = c.ReadByteFromByte(&Cycles, ZeroPageAddress, memory)
			c.LDASetStatus()
		case INS_LDA_ZPX:
			ZeroPageAddress := c.FetchByte(&Cycles, memory)
			ZeroPageAddress += c.X
			Cycles--
			c.A = c.ReadByteFromByte(&Cycles, ZeroPageAddress, memory)
			c.LDASetStatus()
		case INS_LDA_ABS:
			AbsAddress := c.FetchWord(&Cycles, memory)
			c.A = c.ReadByteFromWord(&Cycles, AbsAddress, memory)
		case INS_LDA_ABSX:
			AbsAddress := c.FetchWord(&Cycles, memory)
			AbsAddressX := AbsAddress + word(c.X)
			c.A = c.ReadByteFromWord(&Cycles, AbsAddressX, memory)
			if AbsAddressX-AbsAddress >= 0xFF {
				Cycles--
			}
		case INS_LDA_ABSY:
			AbsAddress := c.FetchWord(&Cycles, memory)
			AbsAddressY := AbsAddress + word(c.Y)
			c.A = c.ReadByteFromWord(&Cycles, AbsAddressY, memory)
			if AbsAddressY-AbsAddress >= 0xFF {
				Cycles--
			}
		case INS_LDA_INDX:
			ZPAddress := c.FetchByte(&Cycles, memory)
			ZPAddress += c.X
			Cycles--
			EffectiveAddr := c.ReadWord(&Cycles, word(ZPAddress), memory)
			c.A = c.ReadByteFromWord(&Cycles, EffectiveAddr, memory)
		case INS_LDA_INDY:
			ZPAddress := c.FetchByte(&Cycles, memory)
			EffectiveAddr := c.ReadWord(&Cycles, word(ZPAddress), memory)
			EffectiveAddrY := EffectiveAddr + word(c.Y)
			c.A = c.ReadByteFromWord(&Cycles, EffectiveAddrY, memory)
			if EffectiveAddrY-EffectiveAddr >= 0xFF {
				Cycles--
			}
		case INS_JSR:
			SubAddr := c.FetchWord(&Cycles, memory)

			c.SP--
			memory.Data[word(c.SP)] = byte((c.PC - 1) >> 8)
			Cycles--

			c.SP--
			memory.Data[word(c.SP)] = byte(c.PC - 1)
			Cycles--

			c.PC = SubAddr
			Cycles--
		default:
			fmt.Printf("Instruction not handled %d\n", Ins)
			break executeLoop
		}
	}
	return CyclesRequested - Cycles
}

func main() {
	cpu := &CPU{}
	mem := &Mem{}
	cpu.Reset(mem)
	cpu.Y = 0xFF
	mem.Data[0xFFFC] = INS_LDA_ABSY
	mem.Data[0xFFFD] = 0x02
	mem.Data[0xFFFE] = 0x44 // 0x4402
	mem.Data[0x4501] = 0x37 // 0x4402 + 0xFF cross page boundary!
	cpu.Execute(5, mem)
}
