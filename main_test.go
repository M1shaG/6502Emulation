package main

import "testing"

func newTestCPU() (*CPU, *Mem) {
	cpu := &CPU{}
	mem := &Mem{}
	cpu.Reset(mem)
	return cpu, mem
}

func TestResetSetsCorrectInitialState(t *testing.T) {
	// given:
	cpu, _ := newTestCPU()

	// when:

	// then:
	if cpu.SP != 0xFD {
		t.Errorf("expected SP = 0xFD, got 0x%02x", cpu.SP)
	}

	if (cpu.Status&FlagC != 0) || (cpu.Status&FlagZ != 0) || (cpu.Status&FlagI != 0) || (cpu.Status&FlagD != 0) ||
		(cpu.Status&FlagB != 0) || (cpu.Status&FlagV != 0) || (cpu.Status&FlagN != 0) {
		t.Errorf("expected all Flags to be set zero after reset")
	}

	if cpu.A != 0 || cpu.X != 0 || cpu.Y != 0 {
		t.Errorf("expected A, X, Y = 0, got A=%d X=%d Y=%d", cpu.A, cpu.X, cpu.Y)
	}
}

func TestTheCPUDoesNothingWhenWeExecuteZeroCycles(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	const NUM_CYCLES = 0

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cyclesUsed != 0 {
		t.Errorf("cyclesUsed = %d, want 0", cyclesUsed)
	}
}

func TestCPUCanExecuteMoreCyclesThanRequestedIfRequiredByTheInstruction(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDA_IM
	mem.Data[0xFFFD] = 0x84
	const NUM_CYCLES = 1

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}
}

func expectUnaffectedFlags(t *testing.T, before, after byte) {
	t.Helper()
	const unaffectedMask = FlagC | FlagI | FlagD | FlagB | FlagV
	if before&unaffectedMask != after&unaffectedMask {
		t.Errorf("unaffected flags changed: before=%08b after=%08b", before, after)
	}
}

// LDA
func TestLDAImmediateCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDA_IM
	mem.Data[0xFFFD] = 0x84
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(2, mem)

	// then:
	if cpu.A != 0x84 {
		t.Errorf("A = %#02x, want 0x84", cpu.A)
	}

	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN == 0 {
		t.Error("N flag should be set")
	}
	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAImmediateCanAffectTheZeroFlag(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x44
	mem.Data[0xFFFC] = INS_LDA_IM
	mem.Data[0xFFFD] = 0x0
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(2, mem)

	// then:
	if cpu.A != 0x0 {
		t.Errorf("A = %#02x, want 0x84", cpu.A)
	}

	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be clear")
	}
	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAZeroPageCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()

	mem.Data[0xFFFC] = INS_LDA_ZP
	mem.Data[0xFFFD] = 0x42
	mem.Data[0x0042] = 0x37

	// when:
	statusBefore := cpu.Status
	cyclesUsed := cpu.Execute(3, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != 3 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAZeroPageXCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDA_ZPX
	mem.Data[0xFFFD] = 0x42
	mem.Data[0x0042] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(4, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != 4 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAZeroPageXCanLoadAValueIntoTheARegisterWhenItWraps(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0xFF
	mem.Data[0xFFFC] = INS_LDA_ZPX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x007F] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(4, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != 4 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAAbsoluteCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDA_ABS
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4480] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAAbsoluteXCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 1
	mem.Data[0xFFFC] = INS_LDA_ABSX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4481] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAAbsoluteXCanLoadAValueIntoTheARegisterWhenItCrossesAPageBoundary(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0xFF
	mem.Data[0xFFFC] = INS_LDA_ABSX
	mem.Data[0xFFFD] = 0x02
	mem.Data[0xFFFE] = 0x44 // 0x4402
	mem.Data[0x4501] = 0x37 // 0x4402 + 0xFF cross page boundary!
	const NUM_CYCLES = 5
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 5", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAAbsoluteYCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 1
	mem.Data[0xFFFC] = INS_LDA_ABSY
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4481] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAAbsoluteYCanLoadAValueIntoTheARegisterWhenItCrossesAPageBoundary(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0xFF
	mem.Data[0xFFFC] = INS_LDA_ABSY
	mem.Data[0xFFFD] = 0x02
	mem.Data[0xFFFE] = 0x44 // 0x4402
	mem.Data[0x4501] = 0x37 // 0x4402 + 0xFF cross page boundary!
	const NUM_CYCLES = 5
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 5", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAIndirectXCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x04
	mem.Data[0xFFFC] = INS_LDA_INDX
	mem.Data[0xFFFD] = 0x02
	mem.Data[0x0006] = 0x00 // 0x02 + 0x04
	mem.Data[0x0007] = 0x80
	mem.Data[0x8000] = 0x37
	const NUM_CYCLES = 6
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAIndirectYCanLoadAValueIntoTheARegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x04
	mem.Data[0xFFFC] = INS_LDA_INDY
	mem.Data[0xFFFD] = 0x02
	mem.Data[0x0002] = 0x00
	mem.Data[0x0003] = 0x80
	mem.Data[0x8004] = 0x37 // 0x8000 + 0x04
	const NUM_CYCLES = 5
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDAIndirectYCanLoadAValueIntoTheARegisterWhenItCrossesAPage(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0xFF
	mem.Data[0xFFFC] = INS_LDA_INDY
	mem.Data[0xFFFD] = 0x02
	mem.Data[0x0002] = 0x02
	mem.Data[0x0003] = 0x80
	mem.Data[0x8101] = 0x37 // 0x8002 + 0xFF
	const NUM_CYCLES = 6
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.A != 0x37 {
		t.Errorf("A = %#02x, want 0x37", cpu.A)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDX
func TestLDXImmediateCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDX_IM
	mem.Data[0xFFFD] = 0x84
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(2, mem)

	// then:
	if cpu.X != 0x84 {
		t.Errorf("X = %#02x, want 0x84", cpu.X)
	}

	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN == 0 {
		t.Error("N flag should be set")
	}
	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXImmediateCanAffectTheZeroFlag(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x44
	mem.Data[0xFFFC] = INS_LDX_IM
	mem.Data[0xFFFD] = 0x0
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(2, mem)

	// then:
	if cpu.X != 0x0 {
		t.Errorf("X = %#02x, want 0x84", cpu.X)
	}

	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be clear")
	}
	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXZeroPageCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDX_ZP
	mem.Data[0xFFFD] = 0x42
	mem.Data[0x0042] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(3, mem)

	// then:
	if cpu.X != 0x37 {
		t.Errorf("X = %#02x, want 0x37", cpu.X)
	}

	if cyclesUsed != 3 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXZeroPageYCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDX_ZPY
	mem.Data[0xFFFD] = 0x42
	mem.Data[0x0042] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(4, mem)

	// then:
	if cpu.X != 0x37 {
		t.Errorf("X = %#02x, want 0x37", cpu.X)
	}

	if cyclesUsed != 4 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXZeroPageYCanLoadAValueIntoTheXRegisterWhenItWraps(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0xFF
	mem.Data[0xFFFC] = INS_LDX_ZPY
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x007F] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(4, mem)

	// then:
	if cpu.X != 0x37 {
		t.Errorf("X = %#02x, want 0x37", cpu.X)
	}

	if cyclesUsed != 4 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXAbsoluteCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDX_ABS
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4480] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.X != 0x37 {
		t.Errorf("X = %#02x, want 0x37", cpu.X)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXAbsoluteYCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 1
	mem.Data[0xFFFC] = INS_LDX_ABSY
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4481] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.X != 0x37 {
		t.Errorf("X = %#02x, want 0x37", cpu.X)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDXAbsoluteYCanLoadAValueIntoTheXRegisterWhenItCrossesAPageBoundary(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0xFF
	mem.Data[0xFFFC] = INS_LDX_ABSY
	mem.Data[0xFFFD] = 0x02
	mem.Data[0xFFFE] = 0x44 // 0x4402
	mem.Data[0x4501] = 0x37 // 0x4402 + 0xFF cross page boundary!
	const NUM_CYCLES = 5
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.X != 0x37 {
		t.Errorf("X = %#02x, want 0x37", cpu.X)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 5", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

// LDY
func TestLDYImmediateCanLoadAValueIntoTheYRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDY_IM
	mem.Data[0xFFFD] = 0x84
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(2, mem)

	// then:
	if cpu.Y != 0x84 {
		t.Errorf("Y = %#02x, want 0x84", cpu.Y)
	}

	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN == 0 {
		t.Error("N flag should be set")
	}
	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYImmediateCanAffectTheZeroFlag(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x44
	mem.Data[0xFFFC] = INS_LDY_IM
	mem.Data[0xFFFD] = 0x0
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(2, mem)

	// then:
	if cpu.Y != 0x0 {
		t.Errorf("Y = %#02x, want 0x84", cpu.Y)
	}

	if cyclesUsed != 2 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ == 0 {
		t.Error("Z flag should be set")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be clear")
	}
	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYZeroPageCanLoadAValueIntoTheYRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDY_ZP
	mem.Data[0xFFFD] = 0x42
	mem.Data[0x0042] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(3, mem)

	// then:
	if cpu.Y != 0x37 {
		t.Errorf("Y = %#02x, want 0x37", cpu.Y)
	}

	if cyclesUsed != 3 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYZeroPageXCanLoadAValueIntoTheXRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDY_ZPX
	mem.Data[0xFFFD] = 0x42
	mem.Data[0x0042] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(4, mem)

	// then:
	if cpu.Y != 0x37 {
		t.Errorf("Y = %#02x, want 0x37", cpu.Y)
	}

	if cyclesUsed != 4 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYZeroPageXCanLoadAValueIntoTheYRegisterWhenItWraps(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0xFF
	mem.Data[0xFFFC] = INS_LDY_ZPX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0x007F] = 0x37
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(4, mem)

	// then:
	if cpu.Y != 0x37 {
		t.Errorf("Y = %#02x, want 0x37", cpu.Y)
	}

	if cyclesUsed != 4 {
		t.Errorf("cyclesUsed = %d, want 2", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYAbsoluteCanLoadAValueIntoTheYRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	mem.Data[0xFFFC] = INS_LDY_ABS
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4480] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.Y != 0x37 {
		t.Errorf("Y = %#02x, want 0x37", cpu.Y)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYAbsoluteXCanLoadAValueIntoTheYRegister(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 1
	mem.Data[0xFFFC] = INS_LDY_ABSX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFFE] = 0x44 // 0x4480
	mem.Data[0x4481] = 0x37
	const NUM_CYCLES = 4
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.Y != 0x37 {
		t.Errorf("Y = %#02x, want 0x37", cpu.Y)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 4", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestLDYAbsoluteXCanLoadAValueIntoTheYRegisterWhenItCrossesAPageBoundary(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0xFF
	mem.Data[0xFFFC] = INS_LDY_ABSX
	mem.Data[0xFFFD] = 0x02
	mem.Data[0xFFFE] = 0x44 // 0x4402
	mem.Data[0x4501] = 0x37 // 0x4402 + 0xFF cross page boundary!
	const NUM_CYCLES = 5
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if cpu.Y != 0x37 {
		t.Errorf("Y = %#02x, want 0x37", cpu.Y)
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want 5", cyclesUsed)
	}

	if cpu.Status&FlagZ != 0 {
		t.Error("Z flag should be clear")
	}

	if cpu.Status&FlagN != 0 {
		t.Error("N flag should be set")
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestSTAZeroPageCanStoreTheARegisterIntoMemory(t *testing.T) {
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

func TestSTAAZeroPageXCanStoreTheARegisterIntoMemory(t *testing.T) {
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

func TestSTAAbsoluteCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ABS
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFEE] = 0x00
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

func TestSTAAbsoluteXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ABSX
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFEE] = 0x00
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

func TestSTAAbsoluteYCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_ABSY
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFEE] = 0x00
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

func TestSTAAIndirectXCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_INDX
	mem.Data[0xFFFD] = 0x20

	mem.Data[0x002F] = 0x30
	mem.Data[0x0030] = 0x80
	mem.Data[0x8030] = 0x00
	var NUM_CYCLES int32 = 6
	statusBefore := cpu.Status

	// when:
	cyclesUsed := cpu.Execute(NUM_CYCLES, mem)

	// then:
	if mem.Data[0x8030] != 0x2F {
		t.Errorf("mem.Data[0x8030] = %#02x, want 0x2F", mem.Data[0x8030])
	}

	if cyclesUsed != NUM_CYCLES {
		t.Errorf("cyclesUsed = %d, want %d", cyclesUsed, NUM_CYCLES)
	}

	expectUnaffectedFlags(t, statusBefore, cpu.Status)
}

func TestSTAAIndirectYCanStoreTheARegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x0F
	cpu.A = 0x2F
	mem.Data[0xFFFC] = INS_STA_INDY
	mem.Data[0xFFFD] = 0x20
	mem.Data[0x0020] = 0x30
	mem.Data[0x0021] = 0x00
	mem.Data[0x003F] = 0x00
	var NUM_CYCLES int32 = 5
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
func TestSTXZeroPageCanStoreTheXRegisterIntoMemory(t *testing.T) {
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

func TestSTXZeroPageYCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x2F
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

func TestSTXAbsoluteCanStoreTheXRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.X = 0x2F
	mem.Data[0xFFFC] = INS_STX_ABS
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFEE] = 0x00
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
func TestSTYZeroPageCanStoreTheYRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x2F
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

func TestSTYZeroPageXCanStoreTheYRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x2F
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

func TestSTYAbsoluteCanStoreTheYRegisterIntoMemory(t *testing.T) {
	// given:
	cpu, mem := newTestCPU()
	cpu.Y = 0x2F
	mem.Data[0xFFFC] = INS_STY_ABS
	mem.Data[0xFFFD] = 0x80
	mem.Data[0xFFEE] = 0x00
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
