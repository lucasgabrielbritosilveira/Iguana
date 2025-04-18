package cpu

import (
	"iguana/bus"
)

type CPU struct {
	Instructions  []Instruction
	Accumulator   uint8
	X             uint8
	Y             uint8
	PC            uint16
	SP            uint8
	Status        uint8
	fetched       uint8
	addr_abs      uint16
	addr_relative uint16
	opcode        uint8
	cycles        uint8
	bus           *bus.Bus
}

func (cpu *CPU) fetch() uint8 {
	return 0
}

func (cpu *CPU) clock() {

}
func (cpu *CPU) read(address uint16) uint8 {
	return cpu.bus.Read(address)
}

func (cpu *CPU) write(address uint16, data uint8) {
	cpu.bus.Write(address, data)

}
