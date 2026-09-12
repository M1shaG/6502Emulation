package CPU6502

import (
	Bus "nesem/bus"
)

// CPU STRUCT
type CPU struct {
	A, X, Y byte
	SP      byte
	PC      Bus.Word
	Status  byte

	bus *Bus.Bus

	addrAbs Bus.Word
	addrRel Bus.Word
	fetched byte
	temp    Bus.Word

	cycles       byte
	opcode       byte
	currentInstr Instruction
}

// FLAGS
const (
	FlagC byte = 1 << iota
	FlagZ
	FlagI
	FlagD
	FlagB
	FlagU
	FlagV
	FlagN
)

func (c *CPU) GetFlag(f byte) bool {
	return c.Status&f > 0
}

func (c *CPU) SetFlag(f byte, v bool) {
	if v {
		c.Status |= f
	} else {
		c.Status &^= f
	}
}

func (c *CPU) ConnectBus(b *Bus.Bus) {
	c.bus = b
}

func (c *CPU) Read(addr Bus.Word) byte {
	return c.bus.Read(addr)
}

func (c *CPU) Write(addr Bus.Word, data byte) {
	c.bus.Write(addr, data)
}

type Instruction struct {
	Name      string
	Operate   func(*CPU) byte
	AddrMode  func(*CPU) byte
	isImplied bool
	Cycles    byte
}

func (c *CPU) fetch() byte {
	if !c.currentInstr.isImplied {
		c.fetched = c.Read(c.addrAbs)
	}
	return c.fetched

}

func (c *CPU) Reset() {
	c.addrAbs = 0xFFFC
	lo := c.Read(c.addrAbs + 0)
	hi := c.Read(c.addrAbs + 1)

	c.PC = (Bus.Word(hi) << 8) | Bus.Word(lo)
	c.A = 0
	c.X = 0
	c.Y = 0
	c.SP = 0
	c.Status = 0x00

	c.addrRel = 0x0000
	c.addrAbs = 0x0000
	c.fetched = 0x00
	c.cycles = 8
}

func (c *CPU) Clock() {
	if c.cycles == 0 {
		c.opcode = c.Read(c.PC)
		c.PC++
		c.currentInstr = opcodeTable[c.opcode]
		c.cycles = opcodeTable[c.opcode].Cycles

		additional_cycle1 := c.currentInstr.AddrMode(c)
		additional_cycle2 := c.currentInstr.Operate(c)

		c.cycles += (additional_cycle1 & additional_cycle2)
	}
	c.cycles--
}

// Addresses
func (c *CPU) IMP() byte {
	c.fetched = c.A
	return 0
}

func (c *CPU) IMM() byte {
	c.addrAbs = c.PC
	c.PC++
	return 0
}

func (c *CPU) ZP0() byte {
	c.addrAbs = Bus.Word(c.Read(c.PC))
	c.PC++
	return 0
}

func (c *CPU) ZPX() byte {
	//fmt.Printf("%#02x, %#02x = %#02x", c.Read(c.PC), c.X, c.Read(c.PC)+c.X)
	c.addrAbs = Bus.Word(c.Read(c.PC) + c.X)
	c.PC++
	c.addrAbs &= 0x00FF
	return 0
}

func (c *CPU) ZPY() byte {
	c.addrAbs = Bus.Word(c.Read(c.PC) + c.Y)
	c.PC++
	return 0
}

func (c *CPU) REL() byte {
	c.addrRel = Bus.Word(c.Read(c.PC))
	c.PC++
	if c.addrRel&0x80 != 0 {
		c.addrRel |= 0xFF00
	}
	return 0
}

func (c *CPU) ABS() byte {
	lo := c.Read(c.PC)
	c.PC++
	hi := c.Read(c.PC)
	c.PC++

	c.addrAbs = Bus.Word(hi)<<8 | Bus.Word(lo)
	return 0
}

func (c *CPU) ABX() byte {
	lo := c.Read(c.PC)
	c.PC++
	hi := c.Read(c.PC)
	c.PC++

	c.addrAbs = Bus.Word(hi)<<8 | Bus.Word(lo)
	c.addrAbs += Bus.Word(c.X)

	if (c.addrAbs & 0xFF00) != (Bus.Word(hi) << 8) {
		return 1
	} else {
		return 0
	}
}

func (c *CPU) ABY() byte {
	lo := c.Read(c.PC)
	c.PC++
	hi := c.Read(c.PC)
	c.PC++

	c.addrAbs = Bus.Word(hi)<<8 | Bus.Word(lo)
	c.addrAbs += Bus.Word(c.Y)

	if (c.addrAbs & 0xFF00) != (Bus.Word(hi) << 8) {
		return 1
	} else {
		return 0
	}
}

func (c *CPU) IND() byte {
	ptr_lo := c.Read(c.PC)
	c.PC++
	ptr_hi := c.Read(c.PC)
	c.PC++

	ptr := Bus.Word(ptr_hi)<<8 | Bus.Word(ptr_lo)
	if ptr_lo == 0x00FF {
		c.addrAbs = Bus.Word(c.Read(ptr&0xFF00))<<8 | Bus.Word(ptr+0)
	} else {
		c.addrAbs = Bus.Word(c.Read(ptr+1&0xFF00))<<8 | Bus.Word(ptr+0)
	}
	return 0
}

func (c *CPU) IZX() byte {
	t := c.Read(c.PC)
	c.PC++

	lo := c.Read(Bus.Word(t+c.X) & 0x00FF)
	hi := c.Read(Bus.Word(t+c.X+1) & 0x00FF)

	c.addrAbs = Bus.Word(hi)<<8 | Bus.Word(lo)

	return 0
}

func (c *CPU) IZY() byte {
	t := c.Read(c.PC)
	c.PC++

	lo := c.Read(Bus.Word(t) & 0x00FF)
	hi := c.Read(Bus.Word(t+1) & 0x00FF)

	c.addrAbs = Bus.Word(hi)<<8 | Bus.Word(lo)
	c.addrAbs += Bus.Word(c.Y)

	if (c.addrAbs & 0xFF00) != (Bus.Word(hi) << 8) {
		return 1
	} else {
		return 0
	}
}

// Load/Store Operations
// Load
func (c *CPU) LDA() byte {
	c.fetch()
	c.A = c.fetched
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 1
}

func (c *CPU) LDX() byte {
	c.fetch()
	c.X = c.fetched
	c.SetFlag(FlagZ, c.X == 0x00)
	c.SetFlag(FlagN, c.X&0x80 != 0)
	return 1
}

func (c *CPU) LDY() byte {
	c.fetch()
	c.Y = c.fetched
	c.SetFlag(FlagZ, c.Y == 0x00)
	c.SetFlag(FlagN, c.Y&0x80 != 0)
	return 1
}

// Store
func (c *CPU) STA() byte {
	c.Write(c.addrAbs, c.A)
	return 0
}

func (c *CPU) STX() byte {
	c.Write(c.addrAbs, c.X)
	return 0
}

func (c *CPU) STY() byte {
	c.Write(c.addrAbs, c.Y)
	return 0
}

// Register Transfers
func (c *CPU) TAX() byte {
	c.X = c.A
	c.SetFlag(FlagZ, c.X == 0x00)
	c.SetFlag(FlagN, c.X&0x80 != 0)
	return 0
}

func (c *CPU) TAY() byte {
	c.Y = c.A
	c.SetFlag(FlagZ, c.Y == 0x00)
	c.SetFlag(FlagN, c.Y&0x80 != 0)
	return 0
}

func (c *CPU) TXA() byte {
	c.A = c.X
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 0
}

func (c *CPU) TYA() byte {
	c.A = c.Y
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 0
}

// Stack Operations
func (c *CPU) TSX() byte {
	c.X = c.SP
	c.SetFlag(FlagZ, c.X == 0x00)
	c.SetFlag(FlagN, c.X&0x80 != 0)
	return 0
}

func (c *CPU) TXS() byte {
	c.SP = c.X
	return 0
}

func (c *CPU) PHA() byte {
	c.Write(0x0100+Bus.Word(c.SP), c.A)
	c.SP--
	return 0
}

func (c *CPU) PHP() byte {
	c.Write(0x0100+Bus.Word(c.SP), c.Status|FlagB|FlagU)
	c.SetFlag(FlagB, false)
	c.SetFlag(FlagU, false)
	c.SP--
	return 0
}

func (c *CPU) PLA() byte {
	c.SP++
	c.Read(0x0100 + Bus.Word(c.SP))
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 0
}

func (c *CPU) PLP() byte {
	c.SP++
	c.Status = c.Read(0x0100 + Bus.Word(c.SP))
	c.SetFlag(FlagU, true)
	return 0
}

// Logical

func (c *CPU) AND() byte {
	c.fetch()
	c.A = c.A & c.fetched
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 1
}

func (c *CPU) EOR() byte {
	c.fetch()
	c.A = c.A ^ c.fetched
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 1
}

func (c *CPU) ORA() byte {
	c.fetch()
	c.A = c.A | c.fetched
	c.SetFlag(FlagZ, c.A == 0x00)
	c.SetFlag(FlagN, c.A&0x80 != 0)
	return 1
}

func (c *CPU) BIT() byte {
	c.fetch()
	c.temp = Bus.Word(c.A) & Bus.Word(c.fetched)
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x00)
	c.SetFlag(FlagN, c.fetched&(1<<7) != 0)
	c.SetFlag(FlagV, c.fetched&(1<<6) != 0)
	return 0
}

// Arithmetic
func (c *CPU) ADC() byte {
	c.fetch()
	c.temp = Bus.Word(c.A) + Bus.Word(c.fetched)

	if c.GetFlag(FlagC) {
		c.temp += 1
	}

	c.SetFlag(FlagC, c.temp > 255)
	c.SetFlag(FlagV, Bus.Word(c.A)&Bus.Word(c.fetched)&Bus.Word(c.A)^Bus.Word(c.temp) != 0)
	c.SetFlag(FlagN, c.temp&0x80 != 0)
	c.A = byte(c.temp & 0x00FF)
	return 1
}

func (c *CPU) SBC() byte {
	c.fetch()

	value := Bus.Word(c.fetched) ^ 0x00FF
	c.temp = Bus.Word(c.A) + value
	if c.GetFlag(FlagC) {
		c.temp += 1
	}
	c.SetFlag(FlagC, c.temp&0xFF00 != 0)
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0)
	c.SetFlag(FlagV, (c.temp^Bus.Word(c.A))&(c.temp&value) != 0)
	c.SetFlag(FlagN, c.temp&0x0080 != 0)
	c.A = byte(c.temp & 0x00FF)
	return 1
}

func (c *CPU) CMP() byte {
	c.fetch()
	c.temp = Bus.Word(c.A) - Bus.Word(c.fetched)
	c.SetFlag(FlagC, c.A >= c.fetched)
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x000)
	c.SetFlag(FlagN, c.temp&0x0080 != 0)
	return 1
}

func (c *CPU) CPX() byte {
	c.fetch()
	c.temp = Bus.Word(c.X) - Bus.Word(c.fetched)
	c.SetFlag(FlagC, c.X >= c.fetched)
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x0000)
	c.SetFlag(FlagN, c.temp&0x0080 != 0)
	return 0
}

func (c *CPU) CPY() byte {
	c.fetch()
	c.temp = Bus.Word(c.Y) - Bus.Word(c.fetched)
	c.SetFlag(FlagC, c.Y >= c.fetched)
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x0000)
	c.SetFlag(FlagN, c.temp&0x0080 != 0)
	return 0
}

// Increments & Decrements
func (c *CPU) INC() byte {
	c.fetch()
	c.temp = Bus.Word(c.fetched) + 1
	c.Write(c.addrAbs, byte(c.temp&0x00FF))
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x0000)
	c.SetFlag(FlagN, c.temp&0x0080 != 0)
	return 0
}

func (c *CPU) INX() byte {
	c.X++
	c.SetFlag(FlagZ, c.X == 0x00)
	c.SetFlag(FlagN, c.X&0x80 != 0)
	return 0
}

func (c *CPU) INY() byte {
	c.Y++
	c.SetFlag(FlagZ, c.Y == 0x00)
	c.SetFlag(FlagN, c.Y&0x80 != 0)
	return 0
}

func (c *CPU) DEC() byte {
	c.fetch()
	c.temp = Bus.Word(c.fetched) - 1
	c.Write(c.addrAbs, byte(c.temp&0x00FF))
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x0000)
	c.SetFlag(FlagN, c.temp&0x0080 != 0)
	return 0
}

func (c *CPU) DEX() byte {
	c.X--
	c.SetFlag(FlagZ, c.X == 0x00)
	c.SetFlag(FlagN, c.X&0x80 != 0)
	return 0
}

func (c *CPU) DEY() byte {
	c.Y--
	c.SetFlag(FlagZ, c.Y == 0x00)
	c.SetFlag(FlagN, c.Y&0x80 != 0)
	return 0
}

// Shifts
func (c *CPU) ASL() byte {
	c.fetch()
	c.temp = Bus.Word(c.fetched) << 1
	c.SetFlag(FlagC, (c.temp&0xFF00) > 0)
	c.SetFlag(FlagZ, (c.temp&0x00FF) == 0x00)
	c.SetFlag(FlagN, c.temp&0x80 != 0)
	// TODO: IMP Addr
	c.Write(c.addrAbs, byte(c.temp&0x00FF))
	return 0
}

func (c *CPU) LSR() byte {
	c.fetch()

	return 0
}

func (c *CPU) BCC() byte {
	return 0
}

func (c *CPU) BCS() byte {
	return 0
}

func (c *CPU) BEQ() byte {
	return 0
}

func (c *CPU) BMI() byte {
	return 0
}

func (c *CPU) BNE() byte {
	return 0
}

func (c *CPU) BPL() byte {
	return 0
}

func (c *CPU) BRK() byte {
	return 0
}

func (c *CPU) BVC() byte {
	return 0
}

func (c *CPU) BVS() byte {
	return 0
}

func (c *CPU) CLC() byte {
	return 0
}

func (c *CPU) CLD() byte {
	return 0
}

func (c *CPU) CLI() byte {
	return 0
}

func (c *CPU) CLV() byte {
	return 0
}

func (c *CPU) JMP() byte {
	return 0
}

func (c *CPU) JSR() byte {
	return 0
}

func (c *CPU) NOP() byte {
	return 0
}

func (c *CPU) ROL() byte {
	return 0
}

func (c *CPU) ROR() byte {
	return 0
}

func (c *CPU) RTI() byte {
	return 0
}

func (c *CPU) RTS() byte {
	return 0
}

func (c *CPU) SEC() byte {
	return 0
}

func (c *CPU) SED() byte {
	return 0
}

func (c *CPU) SEI() byte {
	return 0
}

func (c *CPU) XXX() byte {
	return 0
}
