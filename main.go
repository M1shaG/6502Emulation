package main

import (
	CPU6502 "nesem/6502"
	Bus "nesem/bus"
)

func main() {
	bus := Bus.NewBus()
	cpu := &CPU6502.CPU{}
	cpu.ConnectBus(bus)

	bus.RAM[0xFFFC] = 0x00
	bus.RAM[0xFFFD] = 0x00

	bus.RAM[0x0000] = 0xA9
	bus.RAM[0x0001] = 0x42

	cpu.Reset()

	for i := 0; i < 10; i++ {
		cpu.Clock()
	}
}
