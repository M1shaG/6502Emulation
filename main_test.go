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

	// when:
	mem.Data[0xFFFC] = INS_LDA_IM
	mem.Data[0xFFFD] = 0x84
	const NUM_CYCLES = 1

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

	// when:
	statusBefore := cpu.Status
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

	// when:
	statusBefore := cpu.Status
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

	// when:
	statusBefore := cpu.Status
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

	// when:
	statusBefore := cpu.Status
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
