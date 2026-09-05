package main

import "testing"

// STA
func STAZeroPageCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ZP
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 3
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STAAZeroPageXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	cpu.X = 0x0F
	mem.Data[0xFFFC] = INS_STA_ZPX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x008F] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x008F] != 0x2F {
		t.Errorf("mem.Data[0x008F] = %#02x, want 0x2F", mem.Data[0x008F])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STAAbsoluteCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ABS
	mem.Data[0xFFFD] = 0x00
	mem.Data[0xFFEE] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STAAbsoluteXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ABSX
	mem.Data[0xFFFD] = 0x00
	mem.Data[0xFFEE] = 0x80
	mem.Data[0x008F] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x008F] != 0x2F {
		t.Errorf("mem.Data[0x008F] = %#02x, want 0x2F", mem.Data[0x008F])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STAAbsoluteYCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ABSY
	mem.Data[0xFFFD] = 0x00
	mem.Data[0xFFEE] = 0x80
	mem.Data[0x008F] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x008F] != 0x2F {
		t.Errorf("mem.Data[0x008F] = %#02x, want 0x2F", mem.Data[0x008F])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STAAIndirectXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_INDX
	mem.Data[0xFFFD] = 0x20
	mem.Data[0x002F] = 0x30
	mem.Data[0x0030] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 6
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STAAIndirectYCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_INDY
	mem.Data[0xFFFD] = 0x20
	mem.Data[0x0020] = 0x30
	mem.Data[0x003F] = 0x00
	var NUM_CYCLES int32 = 6
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x003F] != 0x2F {
		t.Errorf("mem.Data[0x003F] = %#02x, want 0x2F", mem.Data[0x003F])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// STX
func STXZeroPageCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x2F
	mem.Data[0xFFFC] = INS_STX_ZP
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 3
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STXZeroPageYCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	cpu.Y = 0x0F
	mem.Data[0xFFFC] = INS_STX_ZPY
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x008F] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x008F] != 0x2F {
		t.Errorf("mem.Data[0x008F] = %#02x, want 0x2F", mem.Data[0x008F])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STXAbsoluteCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STX_ABS
	mem.Data[0xFFFD] = 0x00
	mem.Data[0xFFEE] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// STY
func STYZeroPageCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x2F
	mem.Data[0xFFFC] = INS_STY_ZP
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 3
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STYZeroPageXCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	cpu.X = 0x0F
	mem.Data[0xFFFC] = INS_STY_ZPX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x008F] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x008F] != 0x2F {
		t.Errorf("mem.Data[0x008F] = %#02x, want 0x2F", mem.Data[0x008F])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func STYAbsoluteCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STY_ABS
	mem.Data[0xFFFD] = 0x00
	mem.Data[0xFFEE] = 0x80
	mem.Data[0x0080] = 0x00
	var NUM_CYCLES int32 = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x0080] != 0x2F {
		t.Errorf("mem.Data[0x0080] = %#02x, want 0x2F", mem.Data[0x0080])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}
