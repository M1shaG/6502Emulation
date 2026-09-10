package tests

import (
	CPU6502 "nesem/6502"
	Bus "nesem/bus"
	"testing"
)

func expectUnaffectedFlagsOnStore(t *testing.T, before, after byte) {
	t.Helper()
	const unaffectedMask = CPU6502.FlagC | CPU6502.FlagZ | CPU6502.FlagI | CPU6502.FlagD | CPU6502.FlagB | CPU6502.FlagV | CPU6502.FlagN
	if before&unaffectedMask != after&unaffectedMask {
		t.Errorf("unaffected flags changed: before=%08b after=%08b", before, after)
	}
}

// Store Register
// Zero Page
func StoreRegisterZeroPage(t *testing.T, opcode byte, setReg func(*CPU6502.CPU, byte)) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x0080] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	setReg(cpu, 0x84)

	// when:
	for i := 0; i < 11; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x0080] != 0x84 {
		t.Errorf("bus.RAM[0x0080] = %#02x, want 0x84", bus.RAM[0x0080])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// STA Zero Page
func TestSTAZeroPageCanStoreTheARegisterIntoMemory(t *testing.T) {
	StoreRegisterZeroPage(t, 0x85, func(c *CPU6502.CPU, value byte) {
		c.A = value
	})
}

// STX Zero Page
func TestSTXZeroPageCanStoreTheXRegisterIntoMemory(t *testing.T) {
	StoreRegisterZeroPage(t, 0x86, func(c *CPU6502.CPU, value byte) {
		c.X = value
	})
}

// STY Zero Page
func TestSTYZeroPageCanStoreTheYRegisterIntoMemory(t *testing.T) {
	StoreRegisterZeroPage(t, 0x84, func(c *CPU6502.CPU, value byte) {
		c.Y = value
	})
}

// Zero Page X
func StoreRegisterZeroPageX(t *testing.T, opcode byte, setReg func(*CPU6502.CPU, byte)) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x008F] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x0F
	setReg(cpu, 0x84)

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x008F] != 0x84 {
		t.Errorf("bus.RAM[0x008F] = %#02x, want 0x84", bus.RAM[0x008F])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// STA Zero Page X
func TestSTAZeroPageXCanStoreTheARegisterIntoMemory(t *testing.T) {
	StoreRegisterZeroPageX(t, 0x95, func(c *CPU6502.CPU, value byte) {
		c.A = value
	})
}

// STY Zero Page X
func TestSTYZeroPageXCanStoreTheYRegisterIntoMemory(t *testing.T) {
	StoreRegisterZeroPageX(t, 0x94, func(c *CPU6502.CPU, value byte) {
		c.Y = value
	})
}

// Zero Page Y
func StoreRegisterZeroPageY(t *testing.T, opcode byte, setReg func(*CPU6502.CPU, byte)) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x80
	bus.RAM[0x008F] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x0F
	setReg(cpu, 0x84)

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x008F] != 0x84 {
		t.Errorf("bus.RAM[0x008F] = %#02x, want 0x84", bus.RAM[0x008F])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// STX Zero Page Y
func TestSTXZeroPageYCanStoreTheYRegisterIntoMemory(t *testing.T) {
	StoreRegisterZeroPageY(t, 0x96, func(c *CPU6502.CPU, value byte) {
		c.X = value
	})
}

// Absolute
func StoreRegisterAbsolute(t *testing.T, opcode byte, setReg func(*CPU6502.CPU, byte)) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = opcode
	bus.RAM[0x0001] = 0x00
	bus.RAM[0x0002] = 0x80
	bus.RAM[0x8000] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	setReg(cpu, 0x84)

	// when:
	for i := 0; i < 12; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x8000] != 0x84 {
		t.Errorf("bus.RAM[0x8000] = %#02x, want 0x84", bus.RAM[0x8000])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// STA Absolute
func TestSTAAbsoluteCanStoreTheARegisterIntoMemory(t *testing.T) {
	StoreRegisterAbsolute(t, 0x8D, func(c *CPU6502.CPU, value byte) {
		c.A = value
	})
}

// STX Absolute
func TestSTXAbsoluteCanStoreTheXRegisterIntoMemory(t *testing.T) {
	StoreRegisterAbsolute(t, 0x8E, func(c *CPU6502.CPU, value byte) {
		c.X = value
	})
}

// STY Absolute
func TestSTYAbsoluteCanStoreTheYRegisterIntoMemory(t *testing.T) {
	StoreRegisterAbsolute(t, 0x8C, func(c *CPU6502.CPU, value byte) {
		c.Y = value
	})
}

// Absolute X
// STA Absolute X
func TestSTAAbsoluteXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x9D // STA Absolute X
	bus.RAM[0x0001] = 0x00
	bus.RAM[0x0002] = 0x80
	bus.RAM[0x8001] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x01
	cpu.A = 0x84

	// when:
	for i := 0; i < 13; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x8001] != 0x84 {
		t.Errorf("bus.RAM[0x8001] = %#02x, want 0x84", bus.RAM[0x8001])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// Absolute Y
// STA Absolute Y
func TestSTAAbsoluteYCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x99 // STA Absolute Y
	bus.RAM[0x0001] = 0x00
	bus.RAM[0x0002] = 0x80
	bus.RAM[0x8001] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x01
	cpu.A = 0x84

	// when:
	for i := 0; i < 13; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x8001] != 0x84 {
		t.Errorf("bus.RAM[0x8001] = %#02x, want 0x84", bus.RAM[0x8001])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// Indirect X
// STA Indirect X
func TestSTAAIndirectXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x81 // STA Indirect X
	bus.RAM[0x0001] = 0x20
	bus.RAM[0x002F] = 0x00
	bus.RAM[0x0030] = 0x80
	bus.RAM[0x8000] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x0F
	cpu.A = 0x84

	// when:
	for i := 0; i < 14; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x8000] != 0x84 {
		t.Errorf("bus.RAM[0x8001] = %#02x, want 0x84", bus.RAM[0x8000])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}

// Indirect Y
// STA Indirect Y
func TestSTAAIndirectYCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x91 // STA Indirect Y
	bus.RAM[0x0001] = 0x20
	bus.RAM[0x0020] = 0x00
	bus.RAM[0x0021] = 0x80
	bus.RAM[0x800F] = 0x00

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x0F
	cpu.A = 0x84

	// when:
	for i := 0; i < 14; i++ {
		cpu.Clock()
	}

	// then:
	if bus.RAM[0x800F] != 0x84 {
		t.Errorf("bus.RAM[0x8001] = %#02x, want 0x84", bus.RAM[0x800F])
	}

	expectUnaffectedFlagsOnStore(t, statusBefore, cpu.Status)
}
