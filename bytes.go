package main

// Manipulate with bits
//
// cpu.Status |= FlagX -> Set X status to 1
// cpu.Status &^= FlagX -> Set X status to 0
// cpu.Status ^= FlagX -> test whether X is set
// cpu.Status ^= FlagX -> toggle X, e.g X was 0 became 1, was 1 became 0
const (
	FlagC byte = 1 << iota
	FlagZ
	FlagI
	FlagD
	FlagB
	_
	FlagV
	FlagN
)

// Sets the given status flag bit to 1.
func (c *CPU) SetFlag(flag byte) {
	c.Status |= flag
}

// Sets the given status flag bit to 0.
func (c *CPU) ClearFlag(flag byte) {
	c.Status &^= flag
}

// Checks is the current flag active
func (c *CPU) GetFlag(flag byte) bool {
	return c.Status&flag != 0
}
