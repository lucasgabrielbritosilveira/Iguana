package cpu

func (cpu *CPU) Add(value uint16) {
	var tmp uint16 = uint16(cpu.Accumulator) + value + uint16(cpu.Status["C"])
	if tmp > 255 {
		cpu.Status["C"] = 1
	}
	if tmp&0x00FF == 0 {
		cpu.Status["Z"] = 1
	}
	if (uint16(cpu.Accumulator^cpu.fetched))&(uint16(cpu.Accumulator)^tmp)&0x0080 != 0 {
		cpu.Status["V"] = 1
	}
	if tmp&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	cpu.Accumulator = uint8(tmp & 0x00FF)
}

func (cpu *CPU) Branch() {
	cpu.cycles++
	cpu.addr_abs = cpu.PC + cpu.addr_relative

	if (cpu.addr_abs & 0xFF00) != (cpu.PC & 0xFF00) {
		cpu.cycles++
	}
	cpu.PC = cpu.addr_abs
}

func (cpu *CPU) Compare(parameter uint8) {
	cpu.fetch()
	var tmp uint16 = uint16(parameter) - uint16(cpu.fetched)
	if tmp > 255 {
		cpu.Status["C"] = 1
	}
	if tmp&0x00FF == 0 {
		cpu.Status["Z"] = 1
	}
	if tmp&0x80 != 0 {
		cpu.Status["N"] = 1
	}
}
