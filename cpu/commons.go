package cpu

func (cpu *CPU) Add(value uint16) {
	tmp := uint16(cpu.Accumulator) + value + uint16(cpu.Status["C"])
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
	cpu.addrAbs = cpu.PC + cpu.addrRelative

	if (cpu.addrAbs & 0xFF00) != (cpu.PC & 0xFF00) {
		cpu.cycles++
	}
	cpu.PC = cpu.addrAbs
}

func (cpu *CPU) Compare(parameter uint8) {
	cpu.fetch()
	tmp := uint16(parameter) - uint16(cpu.fetched)
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

func (cpu *CPU) Load(parameter *uint8) {
	cpu.fetch()
	*parameter = cpu.fetched
	if cpu.Status["Z"] == 0 {
		cpu.Accumulator = 0
	} else {
		cpu.Accumulator = 1
	}
	cpu.Status["N"] = cpu.Accumulator & 0x80
}

func (cpu *CPU) Transfer(source *uint8, target *uint8) {
	*target = *source
	if *target == 0 {
		cpu.Status["Z"] = 0
	} else {
		cpu.Status["Z"] = 1
	}
	cpu.Status["N"] = *target & 0x80
}

func (cpu *CPU) Flags() byte {
	var flags byte
	flags |= cpu.Status["C"] << 0
	flags |= cpu.Status["Z"] << 1
	flags |= cpu.Status["I"] << 2
	flags |= cpu.Status["D"] << 3
	flags |= cpu.Status["B"] << 4
	flags |= cpu.Status["U"] << 5
	flags |= cpu.Status["V"] << 6
	flags |= cpu.Status["N"] << 7
	return flags
}

func (cpu *CPU) SetFlags(flags byte) {
	cpu.Status["C"] = (flags >> 0) & 1
	cpu.Status["Z"] = (flags >> 1) & 1
	cpu.Status["I"] = (flags >> 2) & 1
	cpu.Status["D"] = (flags >> 3) & 1
	cpu.Status["B"] = (flags >> 4) & 1
	cpu.Status["U"] = (flags >> 5) & 1
	cpu.Status["V"] = (flags >> 6) & 1
	cpu.Status["N"] = (flags >> 7) & 1
}
