package Bus

type Word uint16

type Bus struct {
	RAM [64 * 1024]byte
}

func NewBus() *Bus {
	b := &Bus{}
	return b
}

func (b *Bus) Write(addr Word, data byte) {
	b.RAM[addr] = data
}

func (b *Bus) Read(addr Word) byte {
	return b.RAM[addr]
}
