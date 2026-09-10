package tests

import (
	CPU6502 "nesem/6502"
	Bus "nesem/bus"
	"testing"
)

// TAX
func TestTransferAXCanTransferANonNegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xAA // TAX

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x01
	cpu.A = 0x42

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x42 {
		t.Errorf("cpu.X = %#02x, want 0x42", cpu.X)
	}

	if cpu.A != 0x42 {
		t.Errorf("cpu.A = %#02x, want 0x42", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TAX
func TestTransferAXCanTransferANonNegativeZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xAA // TAX

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x01
	cpu.A = 0x00

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x00 {
		t.Errorf("cpu.X = %#02x, want 0x00", cpu.X)
	}

	if cpu.A != 0x00 {
		t.Errorf("cpu.A = %#02x, want 0x00", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TAX
func TestTransferAXCanTransferANegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xAA // TAX

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x01
	cpu.A = 0x8B

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x8B {
		t.Errorf("cpu.X = %#02x, want 0x8B", cpu.X)
	}

	if cpu.A != 0x8B {
		t.Errorf("cpu.A = %#02x, want 0x8B", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TAY
func TestTransferAYCanTransferANonNegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xA8 // TAY

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x01
	cpu.A = 0x42

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.Y != 0x42 {
		t.Errorf("cpu.Y = %#02x, want 0x42", cpu.Y)
	}

	if cpu.A != 0x42 {
		t.Errorf("cpu.A = %#02x, want 0x42", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TAY
func TestTransferAYCanTransferANonNegativeZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xA8 // TAY

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x01
	cpu.A = 0x00

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.Y != 0x00 {
		t.Errorf("cpu.Y = %#02x, want 0x00", cpu.Y)
	}

	if cpu.A != 0x00 {
		t.Errorf("cpu.A = %#02x, want 0x00", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TAY
func TestTransferAYCanTransferANegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0xA8 // TAY

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x01
	cpu.A = 0x8B

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.Y != 0x8B {
		t.Errorf("cpu.Y = %#02x, want 0x8B", cpu.Y)
	}

	if cpu.A != 0x8B {
		t.Errorf("cpu.A = %#02x, want 0x8B", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TXA
func TestTransferXACanTransferANonNegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x8A // TXA

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x42
	cpu.A = 0x01

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x42 {
		t.Errorf("cpu.X = %#02x, want 0x42", cpu.X)
	}

	if cpu.A != 0x42 {
		t.Errorf("cpu.A = %#02x, want 0x42", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TXA
func TestTransferXACanTransferANonNegativeZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x8A // TXA

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x00
	cpu.A = 0x01

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x00 {
		t.Errorf("cpu.Y = %#02x, want 0x00", cpu.Y)
	}

	if cpu.A != 0x00 {
		t.Errorf("cpu.A = %#02x, want 0x00", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TXA
func TestTransferXACanTransferANegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x8A // TXA

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.X = 0x8B
	cpu.A = 0x01

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.X != 0x8B {
		t.Errorf("cpu.X = %#02x, want 0x8B", cpu.X)
	}

	if cpu.A != 0x8B {
		t.Errorf("cpu.A = %#02x, want 0x8B", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TYA
func TestTransferYACanTransferANonNegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x98 // TYA

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x42
	cpu.A = 0x01

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.Y != 0x42 {
		t.Errorf("cpu.Y = %#02x, want 0x42", cpu.Y)
	}

	if cpu.A != 0x42 {
		t.Errorf("cpu.A = %#02x, want 0x42", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TXY
func TestTransferYACanTransferANonNegativeZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x98 // TYA

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x00
	cpu.A = 0x01

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.Y != 0x00 {
		t.Errorf("cpu.Y = %#02x, want 0x00", cpu.Y)
	}

	if cpu.A != 0x00 {
		t.Errorf("cpu.A = %#02x, want 0x00", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&CPU6502.FlagN != 0 {
		t.Error("N flag should be clear")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// TYA
func TestTransferYACanTransferANegativeNonZeroValue(t *testing.T) {
	// given:
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00
	bus.RAM[0x0000] = 0x98 // TYA

	cpu.Reset()
	statusBefore := cpu.Status
	cpu.Y = 0x8B
	cpu.A = 0x01

	// when:
	for i := 0; i < 10; i++ {
		cpu.Clock()
	}

	// then:
	if cpu.Y != 0x8B {
		t.Errorf("cpu.Y = %#02x, want 0x8B", cpu.Y)
	}

	if cpu.A != 0x8B {
		t.Errorf("cpu.A = %#02x, want 0x8B", cpu.A)
	}

	if cpu.Status&CPU6502.FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&CPU6502.FlagN == 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}
