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
	Status        map[string]uint8
	fetched       uint8
	addr_abs      uint16
	addr_relative uint16
	opcode        uint8
	cycles        uint8
	bus           *bus.Bus
}

func NewCPU() CPU {
	return CPU{
		Status: map[string]uint8{
			"C": 0,
			"Z": 0,
			"I": 0,
			"D": 0,
			"U": 0,
			"V": 0,
			"B": 0,
			"N": 0,
		},
	}
}
func (cpu *CPU) fetch() uint8 {
	if cpu.Instructions[cpu.opcode].AddressingMode.ID != "IMP" {
		cpu.fetched = cpu.read(cpu.addr_abs)
	}
	return cpu.fetched
}

func (cpu *CPU) clock() {
	if cpu.cycles == 0 {
		cpu.opcode = cpu.read(cpu.PC)
		cpu.PC++
		cpu.cycles = cpu.Instructions[cpu.opcode].Cycle
		cpu.cycles += cpu.Instructions[cpu.opcode].AddressingMode.Run() & cpu.Instructions[cpu.opcode].Operator.Run()
	}
	cpu.cycles--

}

func (cpu *CPU) reset() {

}
func (cpu *CPU) irq() {

}
func (cpu *CPU) nmi() {

}
func (cpu *CPU) read(address uint16) uint8 {
	return cpu.bus.Read(address)
}

func (cpu *CPU) write(address uint16, data uint8) {
	cpu.bus.Write(address, data)

}
