package main

import "fmt"

// 6502 Documentation
// https://6502.org/users/obelisk/6502/

// Memory

// Fill memory with zeros
func (memory *Mem) Init() {
	for i := 0; i < MAX_MEM; i++ {
		memory.Data[i] = 0
	}
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

// CPU

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

// Addressing modes
// Immediate
// Immediate addressing allows specify an 8 bit constant within the instruction.
// Example: LDA #10 -> Load 10 into the accumulator
// We already have FetchByte() func which implement this mode

// Zero Page
// This addressing mode can references only to first 256bytes of memory.
// And use only the least significant byte of the address. The most significant byte is always zero.
// For example , LDY $10 -> Load from memory address $10 value into Y register
func (c *CPU) AddrZeroPage(Cycles *int32, memory *Mem) byte {
	ZeroPageAddr := c.FetchByte(Cycles, memory)
	return ZeroPageAddr
}

// Zero Page X
// Same as Zero page, but we add to address value from X register.
// Wraps the byte if exceed $008F. For example, $80 + $FF = $8F
//
//	1000 0000 ($80)
//
// + 1111 1111 ($FF)
// 1 0111 1111 ($17F) 9 bit drop last bit/
//
//	0111 1111 ($7F)
//
// For example , LDX $10, Y -> Load from memory address $10 + Y value into X register
func (c *CPU) AddrZeroPageX(Cycles *int32, memory *Mem) byte {
	ZeroPageAddr := c.FetchByte(Cycles, memory)
	ZeroPageAddr += c.X
	*Cycles--
	return ZeroPageAddr
}

// Zero Page Y
func (c *CPU) AddrZeroPageY(Cycles *int32, memory *Mem) byte {
	ZeroPageAddr := c.FetchByte(Cycles, memory)
	ZeroPageAddr += c.Y
	*Cycles--
	return ZeroPageAddr
}

// Absolute
// Can references to full 16bit address\
// For examle, LDA $1000 -> Load value into accumulator from address $1000
func (c *CPU) AddrAbsolute(Cycles *int32, memory *Mem) word {
	AbsAddress := c.FetchWord(Cycles, memory)
	return AbsAddress
}

// Absolute X
// Can references to full 16bit address, but we add to Address X value
// For examle, LDA $1000, X -> Load value into accumulator from address $1000 + X
func (c *CPU) AddrAbsoluteX(Cycles *int32, memory *Mem) word {
	AbsAddress := c.FetchWord(Cycles, memory)
	AbsAddressX := AbsAddress + word(c.X)
	if AbsAddressX-AbsAddress >= 0xFF {
		*Cycles--
	}
	return AbsAddressX
}

// Absolute Y
// Can references to full 16bit address, but we add to Address Y value
// For examle, LDA $1000, Y -> Load value into accumulator from address $1000 + Y
func (c *CPU) AddrAbsoluteY(Cycles *int32, memory *Mem) word {
	AbsAddress := c.FetchWord(Cycles, memory)
	AbsAddressY := AbsAddress + word(c.Y)
	if AbsAddressY-AbsAddress >= 0xFF {
		*Cycles--
	}
	return AbsAddressY
}

// Indexed Indirect (INDX)
// We have memory address, for example $10, then to that address we add
// value from X register and from that value we read byte as address and
// go to that address
func (c *CPU) AddrIndirectX(Cycles *int32, memory *Mem) word {
	ZPAddress := c.FetchByte(Cycles, memory)
	ZPAddress += c.X
	*Cycles--
	EffectiveAddr := c.ReadWord(Cycles, word(ZPAddress), memory)
	return EffectiveAddr
}

// TODO: COMMENTS
// Indirect Indexed (INDY)
func (c *CPU) AddrIndirectY(Cycles *int32, memory *Mem) word {
	ZPAddress := c.FetchByte(Cycles, memory)
	EffectiveAddr := c.ReadWord(Cycles, word(ZPAddress), memory)
	EffectiveAddrY := EffectiveAddr + word(c.Y)
	if EffectiveAddrY-EffectiveAddr >= 0xFF {
		*Cycles--
	}
	return EffectiveAddrY
}

// LDA, LDX, LDY helpful function
// Set correct status when instruction executed
func (c *CPU) LoadRegisterSetStatus(Register byte) {
	if Register == 0 {
		c.SetFlag(FlagZ)
	}

	if (Register & 0b10000000) > 0 {
		c.SetFlag(FlagN)
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
	// LDX
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
	// STA
	INS_STA_ZP   byte = 0x85
	INS_STA_ZPX  byte = 0x95
	INS_STA_ABS  byte = 0x8B
	INS_STA_ABSX byte = 0x9D
	INS_STA_ABSY byte = 0x99
	INS_STA_INDX byte = 0x81
	INS_STA_INDY byte = 0x91
	// STX
	INS_STX_ZP  byte = 0x86
	INS_STX_ZPY byte = 0x96
	INS_STX_ABS byte = 0x8E
	// STY
	INS_STY_ZP  byte = 0x84
	INS_STY_ZPX byte = 0x94
	INS_STY_ABS byte = 0x8C

	INS_JSR byte = 0x20
)

func (c *CPU) Execute(Cycles int32, memory *Mem) int32 {
	CyclesRequested := Cycles

executeLoop:
	for Cycles > 0 {
		Ins := c.FetchByte(&Cycles, memory)

		switch Ins {
		//LDA Immediate
		case INS_LDA_IM:
			Value := c.FetchByte(&Cycles, memory)
			c.A = Value
			c.LoadRegisterSetStatus(c.A)
		//LDX Immediate
		case INS_LDX_IM:
			Value := c.FetchByte(&Cycles, memory)
			c.X = Value
			c.LoadRegisterSetStatus(c.X)
		//LDY Immediate
		case INS_LDY_IM:
			Value := c.FetchByte(&Cycles, memory)
			c.Y = Value
			c.LoadRegisterSetStatus(c.Y)
		//LDA Zero Page
		case INS_LDA_ZP:
			Address := c.AddrZeroPage(&Cycles, memory)
			c.A = c.ReadByteFromByte(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		//LDX Zero Page
		case INS_LDX_ZP:
			Address := c.AddrZeroPage(&Cycles, memory)
			c.X = c.ReadByteFromByte(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.X)
		//LDY Zero Page
		case INS_LDY_ZP:
			Address := c.AddrZeroPage(&Cycles, memory)
			c.Y = c.ReadByteFromByte(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.Y)
		//LDA Zero Page X
		case INS_LDA_ZPX:
			Address := c.AddrZeroPageX(&Cycles, memory)
			c.A = c.ReadByteFromByte(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		//LDY Zero Page X
		case INS_LDY_ZPX:
			Address := c.AddrZeroPageX(&Cycles, memory)
			c.Y = c.ReadByteFromByte(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.Y)
		//LDX Zero Page Y
		case INS_LDX_ZPY:
			Address := c.AddrZeroPageY(&Cycles, memory)
			c.X = c.ReadByteFromByte(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.X)
		//LDA Absolute
		case INS_LDA_ABS:
			Address := c.AddrAbsolute(&Cycles, memory)
			c.A = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		//LDY Absolute
		case INS_LDY_ABS:
			Address := c.AddrAbsolute(&Cycles, memory)
			c.Y = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.Y)
		//LDX Absolute
		case INS_LDX_ABS:
			Address := c.AddrAbsolute(&Cycles, memory)
			c.X = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.X)
		//LDA Absolute X
		case INS_LDA_ABSX:
			Address := c.AddrAbsoluteX(&Cycles, memory)
			c.A = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		//LDX Absolute Y
		case INS_LDX_ABSY:
			Address := c.AddrAbsoluteY(&Cycles, memory)
			c.X = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.X)
		//LDY Absolute X
		case INS_LDY_ABSX:
			Address := c.AddrAbsoluteX(&Cycles, memory)
			c.Y = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.Y)
		//LDA Absolute Y
		case INS_LDA_ABSY:
			Address := c.AddrAbsoluteY(&Cycles, memory)
			c.A = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		//LDA Indexed Indirect
		case INS_LDA_INDX:
			Address := c.AddrIndirectX(&Cycles, memory)
			c.A = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		//LDA Indirect Indexed
		case INS_LDA_INDY:
			Address := c.AddrIndirectY(&Cycles, memory)
			c.A = c.ReadByteFromWord(&Cycles, Address, memory)
			c.LoadRegisterSetStatus(c.A)
		// Sometime i will remember u
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
