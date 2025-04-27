package cpu

func (cpu *CPU) Add(value uint16) {
	var temp uint16 = uint16(cpu.Accumulator) + value + uint16(cpu.Status["C"])
	if temp > 255 {
		cpu.Status["C"] = 1
	}
	if temp&0x00FF == 0 {
		cpu.Status["Z"] = 1
	}
	if (uint16(cpu.Accumulator^cpu.fetched))&(uint16(cpu.Accumulator)^temp)&0x0080 != 0 {
		cpu.Status["V"] = 1
	}
	if temp&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	cpu.Accumulator = uint8(temp & 0x00FF)
}

func (cpu *CPU) Branch() {
	cpu.cycles++
	cpu.addr_abs = cpu.PC + cpu.addr_relative

	if (cpu.addr_abs & 0xFF00) != (cpu.PC & 0xFF00) {
		cpu.cycles++
	}
	cpu.PC = cpu.addr_abs
}
