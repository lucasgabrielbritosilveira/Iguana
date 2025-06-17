package cpu

import (
	"iguana/bus"
)

type CPU struct {
	Instructions []Instruction
	Accumulator  uint8
	X            uint8
	Y            uint8
	PC           uint16
	SP           uint8
	Status       map[string]uint8
	fetched      uint8
	addrAbs      uint16
	addrRelative uint16
	opcode       uint8
	cycles       uint8
	bus          *bus.Bus
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
		cpu.fetched = cpu.read(cpu.addrAbs)
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
	cpu.Accumulator = 0
	cpu.X = 0
	cpu.Y = 0
	cpu.SP = 0xFD
	cpu.SetFlags(0)

	cpu.PC = uint16(cpu.read(0xFFFD))<<8 | uint16(cpu.read(0xFFFC))

	cpu.addrAbs = 0
	cpu.addrRelative = 0
	cpu.fetched = 0

	cpu.cycles = 8

}
func (cpu *CPU) irq() {
	if cpu.Status["I"] == 0 {
		cpu.write(0x0100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
		cpu.SP--
		cpu.write(0x0100+uint16(cpu.SP), uint8(cpu.PC&0x00FF))
		cpu.SP--

		cpu.Status["I"] = 1
		cpu.Status["U"] = 1
		cpu.Status["B"] = 0

		cpu.write(0x0100+uint16(cpu.SP), cpu.Flags())
		cpu.SP--

		cpu.addrAbs = 0xFFFE
		cpu.PC = uint16(cpu.read(0xFFFF))<<8 | uint16(cpu.read(0xFFFE))

		cpu.cycles = 7

	}
}
func (cpu *CPU) nmi() {
	cpu.write(0x0100+uint16(cpu.SP), uint8((cpu.PC>>8)&0x00FF))
	cpu.SP--
	cpu.write(0x0100+uint16(cpu.SP), uint8(cpu.PC&0x00FF))
	cpu.SP--

	cpu.Status["I"] = 1
	cpu.Status["U"] = 1
	cpu.Status["B"] = 0

	cpu.write(0x0100+uint16(cpu.SP), cpu.Flags())
	cpu.SP--

	cpu.addrAbs = 0xFFFA
	cpu.PC = uint16(cpu.read(0xFFFB))<<8 | uint16(cpu.read(0xFFFA))

	cpu.cycles = 8
}

func (cpu *CPU) rti() uint8 {
	cpu.SP++
	cpu.SetFlags(cpu.read(0x0100 + uint16(cpu.SP)))
	cpu.SP++

	cpu.PC = uint16(cpu.read(0x0100 + uint16(cpu.SP)))
	cpu.SP++
	cpu.PC |= uint16(cpu.read(0x0100+uint16(cpu.SP))) << 8
	return 0
}
func (cpu *CPU) read(address uint16) uint8 {
	return cpu.bus.Read(address)
}

func (cpu *CPU) write(address uint16, data uint8) {
	cpu.bus.Write(address, data)

}
