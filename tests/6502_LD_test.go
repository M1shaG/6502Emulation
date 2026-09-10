package tests

import (
	CPU6502 "nesem/6502"
	Bus "nesem/bus"
	"testing"
)

func expectUnaffectedFlags(t *testing.T, before, after byte) {
	t.Helper()
	const unaffectedMask = CPU6502.FlagC | CPU6502.FlagI | CPU6502.FlagD | CPU6502.FlagB | CPU6502.FlagV
	if before&unaffectedMask != after&unaffectedMask {
		t.Errorf("unaffected flags changed: before=%08b after=%08b", before, after)
	}
}

// Load Register
// Immediate
func LoadRegisterImmediate(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Immediate
func TestLDAImmediateCanLoadAValueIntoTheARegister(t *testing.T) {
	LoadRegisterImmediate(t, 0xA9, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDX Immediate
func TestLDXImmediateCanLoadAValueIntoTheXRegister(t *testing.T) {
	LoadRegisterImmediate(t, 0xA2, func(c *CPU6502.CPU) byte {
		return c.X
	})
}

// LDY Immediate
func TestLDYImmediateCanLoadAValueIntoTheYRegister(t *testing.T) {
	LoadRegisterImmediate(t, 0xA0, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// Zero Page
func LoadRegisterZeroPage(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x37
	bus.RAM[0x0037] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status

	// when:
	for i := 0; i < 11; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Zero Page
func TestLDAZeroPageCanLoadAValueIntoTheARegister(t *testing.T) {
	LoadRegisterZeroPage(t, 0xA5, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDX Zero Page
func TestLDXZeroPageCanLoadAValueIntoTheXRegister(t *testing.T) {
	LoadRegisterZeroPage(t, 0xA6, func(c *CPU6502.CPU) byte {
		return c.X
	})
}

// LDY Zero Page
func TestLDYZeroPageCanLoadAValueIntoTheYRegister(t *testing.T) {
	LoadRegisterZeroPage(t, 0xA4, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// Zero Page Offset
func LoadRegisterZeroPageX(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x37
	bus.RAM[0x003C] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x5

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Zero Page X
func TestLDAZeroPageXCanLoadAValueIntoTheARegister(t *testing.T) {
	LoadRegisterZeroPageX(t, 0xB5, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDY Zero Page X
func TestLDYZeroPageXCanLoadAValueIntoTheYRegister(t *testing.T) {
	LoadRegisterZeroPageX(t, 0xB4, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// LDX Zero Page Y
func TestLDXZeroPageYCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xB6
	bus.RAM[0x0001] = 0x37
	bus.RAM[0x003C] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x5

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x84 {
		t.Errorf("register = %#02x, want 0x84", cpu.X)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// Wraps Zero Page
func LoadRegisterZeroPageXWhenItWraps(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x007F] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0xFF

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Zero Page X
func TestLDAZeroPageXCanLoadAValueIntoTheARegisterWhenItWraps(t *testing.T) {
	LoadRegisterZeroPageXWhenItWraps(t, 0xB5, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDY Zero Page X
func TestLDYZeroPageXCanLoadAValueIntoTheYRegisterWhenItWraps(t *testing.T) {
	LoadRegisterZeroPageXWhenItWraps(t, 0xB4, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// LDX Zero Page Y
func TestLDXZeroPageYCanLoadAValueIntoTheXRegisterWhenItWraps(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xB6
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x007F] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0xFF

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x84 {
		t.Errorf("register = %#02x, want 0x84", cpu.X)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// Absolute
func LoadRegisterAbsolute(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x0002] = 0x44
	bus.RAM[0x4480] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Absolute
func TestLDAAbsoluteCanLoadAValueIntoTheARegister(t *testing.T) {
	LoadRegisterAbsolute(t, 0xAD, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDX Absolute
func TestLDXZAbsoluteCanLoadAValueIntoTheXRegister(t *testing.T) {
	LoadRegisterAbsolute(t, 0xAE, func(c *CPU6502.CPU) byte {
		return c.X
	})
}

// LDY Absolute
func TestLDYAbsoluteCanLoadAValueIntoTheYRegister(t *testing.T) {
	LoadRegisterAbsolute(t, 0xAC, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// Absolute X
func LoadRegisterAbsoluteX(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x0002] = 0x44
	bus.RAM[0x4481] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x01

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Absolute X
func TestLDAAbsoluteXCanLoadAValueIntoTheARegister(t *testing.T) {
	LoadRegisterAbsoluteX(t, 0xBD, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDY Absolute X
func TestLDYAbsoluteXCanLoadAValueIntoTheYRegister(t *testing.T) {
	LoadRegisterAbsoluteX(t, 0xBC, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// Absolute X Page Boundary
func LoadRegisterAbsoluteXPageBoundary(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0xFF
	bus.RAM[0x0002] = 0x44
	bus.RAM[0x4500] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x01

	// when:
	for i := 0; i < 13; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Absolute X Page Boundary
func TestLDAAbsoluteXCanLoadAValueIntoTheARegisterWhenCrossingPage(t *testing.T) {
	LoadRegisterAbsoluteXPageBoundary(t, 0xBD, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDY Absolute X Page Boundary
func TestLDYAbsoluteXCanLoadAValueIntoTheYRegisterWhenCrossingPage(t *testing.T) {
	LoadRegisterAbsoluteXPageBoundary(t, 0xBC, func(c *CPU6502.CPU) byte {
		return c.Y
	})
}

// Absolute Y
func LoadRegisterAbsoluteY(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x0002] = 0x44
	bus.RAM[0x4481] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x01

	// when:
	for i := 0; i < 11; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Absolute Y
func TestLDAAbsoluteYCanLoadAValueIntoTheARegister(t *testing.T) {
	LoadRegisterAbsoluteY(t, 0xB9, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDY Absolute Y
func TestLDXAbsoluteYCanLoadAValueIntoTheYRegister(t *testing.T) {
	LoadRegisterAbsoluteY(t, 0xBE, func(c *CPU6502.CPU) byte {
		return c.X
	})
}

// Absolute Y Page Boundary
func LoadRegisterAbsoluteYPageBoundary(t *testing.T, opcode byte, getReg func(*CPU6502.CPU) byte) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0xFF
	bus.RAM[0x0002] = 0x44
	bus.RAM[0x4500] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x01

	// when:
	for i := 0; i < 11; i++ {
		cpu.Clock()
	}

	// then:
	reg := getReg(cpu)
	if reg != 0x84 {
		t.Errorf("register = %#02x, want 0x84", reg)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Absolute Y
func TestLDAAbsoluteYCanLoadAValueIntoTheARegisterWhenCrossingPage(t *testing.T) {
	LoadRegisterAbsoluteYPageBoundary(t, 0xB9, func(c *CPU6502.CPU) byte {
		return c.A
	})
}

// LDY Absolute Y
func TestLDXAbsoluteYCanLoadAValueIntoTheYRegisterWhenCrossingPage(t *testing.T) {
	LoadRegisterAbsoluteYPageBoundary(t, 0xBE, func(c *CPU6502.CPU) byte {
		return c.X
	})
}

// LDA Indirect Y
func TestLDAIndirectXCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xA1
	bus.RAM[0x0001] = 0x02
	bus.RAM[0x0006] = 0x00
	bus.RAM[0x0007] = 0x80
	bus.RAM[0x8000] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x04

	// when:
	for i := 0; i < 14; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.A != 0x84 {
		t.Errorf("register = %#02x, want 0x84", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDA Indirect Y
func TestLDAIndirectYCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xB1
	bus.RAM[0x0001] = 0x10
	bus.RAM[0x0010] = 0x00
	bus.RAM[0x0011] = 0x80
	bus.RAM[0x8004] = 0x84

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x04

	// when:
	for i := 0; i < 13; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.A != 0x84 {
		t.Errorf("register = %#02x, want 0x84", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}
