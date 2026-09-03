package main

type word uint16 // handy name for 16bit

// 64 KB
const MAX_MEM = 1024 * 64

type Mem struct {
	Data [MAX_MEM]byte
}

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

type Processor interface {
	Reset(memory *Mem)
	FetchByte(Cycles *uint32, memory *Mem) byte
	ReadByteFromByte(Cycles *uint32, Address byte, memory *Mem) byte
	ReadByteFromWord(Cycles *uint32, Address word, memory *Mem) byte
	ReadWord(Cycles *uint32, Address word, memory *Mem) byte
	FetchWord(Cycles *uint32, memory *Mem) word
	Execute(memory *Mem)
}

