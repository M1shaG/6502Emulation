package CPU6502

var opcodeTable = [256]Instruction{
	// 0x00-0x0F
	{"BRK", (*CPU).BRK, (*CPU).IMM, false, 7}, {"ORA", (*CPU).ORA, (*CPU).IZX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 3}, {"ORA", (*CPU).ORA, (*CPU).ZP0, false, 3}, {"ASL", (*CPU).ASL, (*CPU).ZP0, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"PHP", (*CPU).PHP, (*CPU).IMP, true, 3}, {"ORA", (*CPU).ORA, (*CPU).IMM, false, 2}, {"ASL", (*CPU).ASL, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"ORA", (*CPU).ORA, (*CPU).ABS, false, 4}, {"ASL", (*CPU).ASL, (*CPU).ABS, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},

	// 0x10-0x1F
	{"BPL", (*CPU).BPL, (*CPU).REL, false, 2}, {"ORA", (*CPU).ORA, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"ORA", (*CPU).ORA, (*CPU).ZPX, false, 4}, {"ASL", (*CPU).ASL, (*CPU).ZPX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"CLC", (*CPU).CLC, (*CPU).IMP, true, 2}, {"ORA", (*CPU).ORA, (*CPU).ABY, false, 4}, {"???", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"ORA", (*CPU).ORA, (*CPU).ABX, false, 4}, {"ASL", (*CPU).ASL, (*CPU).ABX, false, 7}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},

	// 0x20-0x2F
	{"JSR", (*CPU).JSR, (*CPU).ABS, false, 6}, {"AND", (*CPU).AND, (*CPU).IZX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"BIT", (*CPU).BIT, (*CPU).ZP0, false, 3}, {"AND", (*CPU).AND, (*CPU).ZP0, false, 3}, {"ROL", (*CPU).ROL, (*CPU).ZP0, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"PLP", (*CPU).PLP, (*CPU).IMP, true, 4}, {"AND", (*CPU).AND, (*CPU).IMM, false, 2}, {"ROL", (*CPU).ROL, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"BIT", (*CPU).BIT, (*CPU).ABS, false, 4}, {"AND", (*CPU).AND, (*CPU).ABS, false, 4}, {"ROL", (*CPU).ROL, (*CPU).ABS, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},

	// 0x30-0x3F
	{"BMI", (*CPU).BMI, (*CPU).REL, false, 2}, {"AND", (*CPU).AND, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"AND", (*CPU).AND, (*CPU).ZPX, false, 4}, {"ROL", (*CPU).ROL, (*CPU).ZPX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"SEC", (*CPU).SEC, (*CPU).IMP, true, 2}, {"AND", (*CPU).AND, (*CPU).ABY, false, 4}, {"???", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"AND", (*CPU).AND, (*CPU).ABX, false, 4}, {"ROL", (*CPU).ROL, (*CPU).ABX, false, 7}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},

	// 0x40-0x4F
	{"RTI", (*CPU).RTI, (*CPU).IMP, true, 6}, {"EOR", (*CPU).EOR, (*CPU).IZX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 3}, {"EOR", (*CPU).EOR, (*CPU).ZP0, false, 3}, {"LSR", (*CPU).LSR, (*CPU).ZP0, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"PHA", (*CPU).PHA, (*CPU).IMP, true, 3}, {"EOR", (*CPU).EOR, (*CPU).IMM, false, 2}, {"LSR", (*CPU).LSR, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"JMP", (*CPU).JMP, (*CPU).ABS, false, 3}, {"EOR", (*CPU).EOR, (*CPU).ABS, false, 4}, {"LSR", (*CPU).LSR, (*CPU).ABS, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},

	// 0x50-0x5F
	{"BVC", (*CPU).BVC, (*CPU).REL, false, 2}, {"EOR", (*CPU).EOR, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"EOR", (*CPU).EOR, (*CPU).ZPX, false, 4}, {"LSR", (*CPU).LSR, (*CPU).ZPX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"CLI", (*CPU).CLI, (*CPU).IMP, true, 2}, {"EOR", (*CPU).EOR, (*CPU).ABY, false, 4}, {"???", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"EOR", (*CPU).EOR, (*CPU).ABX, false, 4}, {"LSR", (*CPU).LSR, (*CPU).ABX, false, 7}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},

	// 0x60-0x6F
	{"RTS", (*CPU).RTS, (*CPU).IMP, true, 6}, {"ADC", (*CPU).ADC, (*CPU).IZX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 3}, {"ADC", (*CPU).ADC, (*CPU).ZP0, false, 3}, {"ROR", (*CPU).ROR, (*CPU).ZP0, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"PLA", (*CPU).PLA, (*CPU).IMP, true, 4}, {"ADC", (*CPU).ADC, (*CPU).IMM, false, 2}, {"ROR", (*CPU).ROR, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"JMP", (*CPU).JMP, (*CPU).IND, false, 5}, {"ADC", (*CPU).ADC, (*CPU).ABS, false, 4}, {"ROR", (*CPU).ROR, (*CPU).ABS, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},

	// 0x70-0x7F
	{"BVS", (*CPU).BVS, (*CPU).REL, false, 2}, {"ADC", (*CPU).ADC, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"ADC", (*CPU).ADC, (*CPU).ZPX, false, 4}, {"ROR", (*CPU).ROR, (*CPU).ZPX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"SEI", (*CPU).SEI, (*CPU).IMP, true, 2}, {"ADC", (*CPU).ADC, (*CPU).ABY, false, 4}, {"???", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"ADC", (*CPU).ADC, (*CPU).ABX, false, 4}, {"ROR", (*CPU).ROR, (*CPU).ABX, false, 7}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},

	// 0x80-0x8F
	{"???", (*CPU).NOP, (*CPU).IMM, false, 2}, {"STA", (*CPU).STA, (*CPU).IZX, false, 6}, {"???", (*CPU).NOP, (*CPU).IMM, false, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"STY", (*CPU).STY, (*CPU).ZP0, false, 3}, {"STA", (*CPU).STA, (*CPU).ZP0, false, 3}, {"STX", (*CPU).STX, (*CPU).ZP0, false, 3}, {"???", (*CPU).XXX, (*CPU).IMP, true, 3},
	{"DEY", (*CPU).DEY, (*CPU).IMP, true, 2}, {"???", (*CPU).NOP, (*CPU).IMM, false, 2}, {"TXA", (*CPU).TXA, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"STY", (*CPU).STY, (*CPU).ABS, false, 4}, {"STA", (*CPU).STA, (*CPU).ABS, false, 4}, {"STX", (*CPU).STX, (*CPU).ABS, false, 4}, {"???", (*CPU).XXX, (*CPU).IMP, true, 4},

	// 0x90-0x9F
	{"BCC", (*CPU).BCC, (*CPU).REL, false, 2}, {"STA", (*CPU).STA, (*CPU).IZY, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"STY", (*CPU).STY, (*CPU).ZPX, false, 4}, {"STA", (*CPU).STA, (*CPU).ZPX, false, 4}, {"STX", (*CPU).STX, (*CPU).ZPY, false, 4}, {"???", (*CPU).XXX, (*CPU).IMP, true, 4},
	{"TYA", (*CPU).TYA, (*CPU).IMP, true, 2}, {"STA", (*CPU).STA, (*CPU).ABY, false, 5}, {"TXS", (*CPU).TXS, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 5}, {"STA", (*CPU).STA, (*CPU).ABX, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},

	// 0xA0-0xAF
	{"LDY", (*CPU).LDY, (*CPU).IMM, false, 2}, {"LDA", (*CPU).LDA, (*CPU).IZX, false, 6}, {"LDX", (*CPU).LDX, (*CPU).IMM, false, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"LDY", (*CPU).LDY, (*CPU).ZP0, false, 3}, {"LDA", (*CPU).LDA, (*CPU).ZP0, false, 3}, {"LDX", (*CPU).LDX, (*CPU).ZP0, false, 3}, {"???", (*CPU).XXX, (*CPU).IMP, true, 3},
	{"TAY", (*CPU).TAY, (*CPU).IMP, true, 2}, {"LDA", (*CPU).LDA, (*CPU).IMM, false, 2}, {"TAX", (*CPU).TAX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"LDY", (*CPU).LDY, (*CPU).ABS, false, 4}, {"LDA", (*CPU).LDA, (*CPU).ABS, false, 4}, {"LDX", (*CPU).LDX, (*CPU).ABS, false, 4}, {"???", (*CPU).XXX, (*CPU).IMP, true, 4},

	// 0xB0-0xBF
	{"BCS", (*CPU).BCS, (*CPU).REL, false, 2}, {"LDA", (*CPU).LDA, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"LDY", (*CPU).LDY, (*CPU).ZPX, false, 4}, {"LDA", (*CPU).LDA, (*CPU).ZPX, false, 4}, {"LDX", (*CPU).LDX, (*CPU).ZPY, false, 4}, {"???", (*CPU).XXX, (*CPU).IMP, true, 4},
	{"CLV", (*CPU).CLV, (*CPU).IMP, true, 2}, {"LDA", (*CPU).LDA, (*CPU).ABY, false, 4}, {"TSX", (*CPU).TSX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 4},
	{"LDY", (*CPU).LDY, (*CPU).ABX, false, 4}, {"LDA", (*CPU).LDA, (*CPU).ABX, false, 4}, {"LDX", (*CPU).LDX, (*CPU).ABY, false, 4}, {"???", (*CPU).XXX, (*CPU).IMP, true, 4},

	// 0xC0-0xCF
	{"CPY", (*CPU).CPY, (*CPU).IMM, false, 2}, {"CMP", (*CPU).CMP, (*CPU).IZX, false, 6}, {"???", (*CPU).NOP, (*CPU).IMM, false, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"CPY", (*CPU).CPY, (*CPU).ZP0, false, 3}, {"CMP", (*CPU).CMP, (*CPU).ZP0, false, 3}, {"DEC", (*CPU).DEC, (*CPU).ZP0, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"INY", (*CPU).INY, (*CPU).IMP, true, 2}, {"CMP", (*CPU).CMP, (*CPU).IMM, false, 2}, {"DEX", (*CPU).DEX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2},
	{"CPY", (*CPU).CPY, (*CPU).ABS, false, 4}, {"CMP", (*CPU).CMP, (*CPU).ABS, false, 4}, {"DEC", (*CPU).DEC, (*CPU).ABS, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},

	// 0xD0-0xDF
	{"BNE", (*CPU).BNE, (*CPU).REL, false, 2}, {"CMP", (*CPU).CMP, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"CMP", (*CPU).CMP, (*CPU).ZPX, false, 4}, {"DEC", (*CPU).DEC, (*CPU).ZPX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"CLD", (*CPU).CLD, (*CPU).IMP, true, 2}, {"CMP", (*CPU).CMP, (*CPU).ABY, false, 4}, {"NOP", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"CMP", (*CPU).CMP, (*CPU).ABX, false, 4}, {"DEC", (*CPU).DEC, (*CPU).ABX, false, 7}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},

	// 0xE0-0xEF
	{"CPX", (*CPU).CPX, (*CPU).IMM, false, 2}, {"SBC", (*CPU).SBC, (*CPU).IZX, false, 6}, {"???", (*CPU).NOP, (*CPU).IMM, false, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"CPX", (*CPU).CPX, (*CPU).ZP0, false, 3}, {"SBC", (*CPU).SBC, (*CPU).ZP0, false, 3}, {"INC", (*CPU).INC, (*CPU).ZP0, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 5},
	{"INX", (*CPU).INX, (*CPU).IMP, true, 2}, {"SBC", (*CPU).SBC, (*CPU).IMM, false, 2}, {"NOP", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).SBC, (*CPU).IMM, false, 2},
	{"CPX", (*CPU).CPX, (*CPU).ABS, false, 4}, {"SBC", (*CPU).SBC, (*CPU).ABS, false, 4}, {"INC", (*CPU).INC, (*CPU).ABS, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},

	// 0xF0-0xFF
	{"BEQ", (*CPU).BEQ, (*CPU).REL, false, 2}, {"SBC", (*CPU).SBC, (*CPU).IZY, false, 5}, {"???", (*CPU).XXX, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 8},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"SBC", (*CPU).SBC, (*CPU).ZPX, false, 4}, {"INC", (*CPU).INC, (*CPU).ZPX, false, 6}, {"???", (*CPU).XXX, (*CPU).IMP, true, 6},
	{"SED", (*CPU).SED, (*CPU).IMP, true, 2}, {"SBC", (*CPU).SBC, (*CPU).ABY, false, 4}, {"NOP", (*CPU).NOP, (*CPU).IMP, true, 2}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
	{"???", (*CPU).NOP, (*CPU).IMP, true, 4}, {"SBC", (*CPU).SBC, (*CPU).ABX, false, 4}, {"INC", (*CPU).INC, (*CPU).ABX, false, 7}, {"???", (*CPU).XXX, (*CPU).IMP, true, 7},
}
