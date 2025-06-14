package cpu

type AddressingMode struct {
	Run func() uint8
	ID  string
}

type Operator struct {
	Run func() uint8
	ID  string
}
type Instruction struct {
	Value          int
	Address        int
	Cycle          uint8
	AddressingMode AddressingMode
	Operator       Operator
}

// Addressing Modes

func (cpu *CPU) imp() uint8 {
	cpu.fetched = cpu.Accumulator
	return 0
}
func (cpu *CPU) zp0() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.PC))
	cpu.addrAbs &= 0x00FF
	cpu.PC++
	return 0
}

func (cpu *CPU) zpy() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.PC + uint16(cpu.Y)))
	cpu.addrAbs &= 0x00FF
	cpu.PC++
	return 0
}

func (cpu *CPU) abs() uint8 {
	lowAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	highAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addrAbs = (highAddr << 8) | lowAddr
	return 0
}

func (cpu *CPU) aby() uint8 {
	lowAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	highAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addrAbs = (highAddr << 8) | lowAddr
	cpu.addrAbs += uint16(cpu.Y)
	if cpu.addrAbs&0xFF00 != (highAddr << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) izx() uint8 {
	temp := uint16(cpu.read(cpu.PC))
	lowAddr := uint16(cpu.read(uint16(temp+uint16(cpu.X)) & 0x00FF))
	highAddr := uint16(cpu.read(uint16(temp+uint16(cpu.X)+1) & 0x00FF))
	cpu.addrAbs = (highAddr << 8) | lowAddr
	cpu.PC++
	return 0
}

func (cpu *CPU) imm() uint8 {
	cpu.addrAbs = cpu.PC
	cpu.PC++
	return 0
}

func (cpu *CPU) zpx() uint8 {
	cpu.addrAbs = uint16(cpu.read(cpu.PC + uint16(cpu.X)))
	cpu.addrAbs &= 0x00FF
	cpu.PC++
	return 0
}

func (cpu *CPU) rel() uint8 {
	cpu.addrRelative = uint16(cpu.read(cpu.PC))
	if cpu.addrRelative&0x80 != 0 {
		cpu.addrRelative |= 0xFF00
	}
	cpu.PC++
	return 0
}

func (cpu *CPU) abx() uint8 {
	lowAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	highAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addrAbs = (highAddr << 8) | lowAddr
	cpu.addrAbs += uint16(cpu.X)
	if cpu.addrAbs&0xFF00 != (highAddr << 8) {
		return 1
	}
	return 0

}
func (cpu *CPU) ind() uint8 {
	lowAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	highAddr := uint16(cpu.read(cpu.PC))
	cpu.PC++
	ptr := (highAddr << 8) | lowAddr
	if lowAddr == 0x00FF {
		cpu.addrAbs = uint16(cpu.read(ptr&0xFF00))<<8 | uint16(cpu.read(ptr))
	} else {
		cpu.addrAbs = uint16(cpu.read(ptr+1))<<8 | uint16(cpu.read(ptr))
	}
	return 0
}

func (cpu *CPU) izy() uint8 {
	tmp := uint16(cpu.read(cpu.PC))
	lowAddr := uint16(cpu.read(tmp & 0x00FF))
	highAddr := uint16(cpu.read(tmp + 1&0x00FF))

	cpu.addrAbs = (highAddr << 8) | lowAddr
	cpu.addrAbs += uint16(cpu.Y)

	cpu.PC++
	if cpu.addrAbs&0x00FF != (highAddr << 8) {
		return 1
	}
	return 0
}

// Arithmetics

func (cpu *CPU) adc() uint8 {
	cpu.fetch()
	cpu.Add(uint16(cpu.fetched))
	return 1
}

func (cpu *CPU) sbc() uint8 {
	value := uint16(cpu.fetched) ^ 0x00FF
	cpu.Add(value)
	return 1
}

func (cpu *CPU) dec() uint8 {
	cpu.fetch()
	tmp := cpu.fetched - 1
	cpu.write(cpu.addrAbs, tmp&0x00FF)
	if tmp == 0 {
		cpu.Status["Z"] = 1
	} else if tmp&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 0
}

func (cpu *CPU) dex() uint8 {
	cpu.X--
	if cpu.X == 0 {
		cpu.Status["Z"] = 1
	} else if cpu.X&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 0
}

func (cpu *CPU) dey() uint8 {
	cpu.Y--
	if cpu.Y == 0 {
		cpu.Status["Z"] = 1
	} else if cpu.Y&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 0
}

func (cpu *CPU) inc() uint8 {
	cpu.fetch()
	tmp := cpu.fetched + 1
	cpu.write(cpu.addrAbs, tmp&0x00FF)
	if tmp == 0 {
		cpu.Status["Z"] = 1
	} else if tmp&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 0
}

func (cpu *CPU) inx() uint8 {
	cpu.X++
	if cpu.X == 0 {
		cpu.Status["Z"] = 1
	} else if cpu.X&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 0
}

func (cpu *CPU) iny() uint8 {
	cpu.Y++
	if cpu.Y == 0 {
		cpu.Status["Z"] = 1
	} else if cpu.Y&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 0
}

// Compare

func (cpu *CPU) cmp() uint8 {
	cpu.Compare(cpu.Accumulator)
	return 1
}

func (cpu *CPU) cpx() uint8 {
	cpu.Compare(cpu.X)
	return 0
}

func (cpu *CPU) cpy() uint8 {
	cpu.Compare(cpu.Y)
	return 0
}

// Bitwise

func asl() {

}

func (cpu *CPU) and() uint8 {
	cpu.fetch()
	cpu.Accumulator &= cpu.fetched
	if cpu.Accumulator == 0 {
		cpu.Status["Z"] = 1
	} else if cpu.Accumulator&0x80 != 0 {
		cpu.Status["N"] = 1
	}
	return 1
}

func (cpu *CPU) eor() uint8 {
	cpu.fetch()
	cpu.Accumulator = ^cpu.fetched
	if cpu.Accumulator == 0 {
		cpu.Status["Z"] = 0
	} else {
		cpu.Status["Z"] = 1
	}
	cpu.Status["N"] = cpu.Accumulator & 0x80
	return 1
}

func (cpu *CPU) ora() uint8 {
	cpu.fetch()
	cpu.Accumulator = cpu.Accumulator | cpu.fetched
	if cpu.Accumulator == 0 {
		cpu.Status["Z"] = 0
	} else {
		cpu.Status["Z"] = 1
	}
	cpu.Status["N"] = cpu.Accumulator & 0x80
	return 1
}

func (cpu *CPU) bit() uint8 {
	cpu.fetch()
	temp := cpu.Accumulator & cpu.fetched
	if temp&0x00FF == 0 {
		cpu.Status["Z"] = 0
	} else {
		cpu.Status["Z"] = 1
	}
	cpu.Status["N"] = cpu.fetched & (1 << 7)
	cpu.Status["V"] = cpu.fetched & (1 << 6)
	return 0
}

func (cpu *CPU) rol() uint8 {
	cpu.fetch()
	temp := uint16((cpu.fetched << 1) | cpu.Status["C"])
	cpu.Status["C"] = uint8(temp & 0xFF00)
	if (temp & 0x00FF) == 0x0000 {
		cpu.Status["Z"] = 0
	} else {
		cpu.Status["Z"] = 1
	}
	cpu.Status["N"] = uint8(temp & 0x0080)

	if cpu.Instructions[cpu.opcode].AddressingMode.ID == "IMP" {
		cpu.Accumulator = uint8(temp & 0x00FF)
	} else {
		cpu.write(cpu.addrAbs, uint8(temp&0x00FF))
	}
	return 0
}

func ror() {

}

func lsr() {

}

// Branches and Jumps

func (cpu *CPU) bcc() uint8 {
	if cpu.Status["C"] == 0 {
		cpu.Branch()
	}
	return 0
}

func (cpu *CPU) bcs() uint8 {
	if cpu.Status["C"] == 1 {
		cpu.Branch()
	}
	return 0
}

func (cpu *CPU) beq() uint8 {
	if cpu.Status["Z"] == 1 {
		cpu.Branch()
	}
	return 0
}

func (cpu *CPU) bmi() uint8 {
	if cpu.Status["N"] == 1 {
		cpu.Branch()
	}
	return 0

}

func (cpu *CPU) bne() uint8 {
	if cpu.Status["Z"] == 0 {
		cpu.Branch()
	}
	return 0
}

func (cpu *CPU) bpl() uint8 {
	if cpu.Status["N"] == 0 {
		cpu.Branch()
	}
	return 0
}

func (cpu *CPU) bvc() uint8 {
	if cpu.Status["V"] == 0 {
		cpu.Branch()
	}
	return 0
}

func (cpu *CPU) bvs() uint8 {
	if cpu.Status["V"] == 1 {
		cpu.Branch()
	}
	return 0
}

func brk() {

}

func jsr() {

}

func jmp() {

}

func rti() {

}

func rts() {

}

// Status

func (cpu *CPU) clc() uint8 {
	cpu.Status["C"] = 0
	return 0
}

func (cpu *CPU) cld() uint8 {
	cpu.Status["D"] = 0
	return 0

}

func (cpu *CPU) cli() uint8 {
	cpu.Status["I"] = 0
	return 0
}

func (cpu *CPU) clv() uint8 {
	cpu.Status["V"] = 0
	return 0
}

func sec() {

}

func sed() {

}

func sei() {

}

//Stack

func (cpu *CPU) pha() uint8 {
	cpu.write(uint16(cpu.SP)+0x0100, cpu.Accumulator)
	cpu.SP--
	return 0
}

func (cpu *CPU) php() uint8 {
	cpu.write(uint16(cpu.SP)+0x0100, cpu.Flags()|0x10)
	cpu.Status["B"] = 0
	cpu.Status["U"] = 0
	cpu.SP--
	return 0
}

func (cpu *CPU) pla() uint8 {
	cpu.SP++
	cpu.Accumulator = cpu.read(0x0100 + uint16(cpu.SP))
	if cpu.Accumulator == 0 {
		cpu.Status["Z"] = 0
	} else {
		cpu.Status["Z"] = 1
	}
	cpu.Status["N"] = cpu.Accumulator & 0x80
	return 0
}

func (cpu *CPU) plp() uint8 {
	cpu.SP++
	cpu.SetFlags(cpu.read(0x0100 + uint16(cpu.SP)))
	cpu.Status["U"] = 1
	return 0
}

// Data

func (cpu *CPU) lda() uint8 {
	cpu.Load(&cpu.Accumulator)
	return 1
}

func (cpu *CPU) ldx() uint8 {
	cpu.Load(&cpu.X)
	return 1
}

func (cpu *CPU) ldy() uint8 {
	cpu.Load(&cpu.Y)
	return 1
}

func (cpu *CPU) sta() uint8 {
	cpu.write(cpu.addrAbs, cpu.Accumulator)
	return 0
}

func (cpu *CPU) stx() uint8 {
	cpu.write(cpu.addrAbs, cpu.X)
	return 0
}

func (cpu *CPU) sty() uint8 {
	cpu.write(cpu.addrAbs, cpu.Y)
	return 0
}

func (cpu *CPU) tax() uint8 {
	cpu.Transfer(&cpu.Accumulator, &cpu.X)
	return 0
}

func (cpu *CPU) tay() uint8 {
	cpu.Transfer(&cpu.Accumulator, &cpu.Y)
	return 0
}

func (cpu *CPU) tsx() uint8 {
	cpu.Transfer(&cpu.SP, &cpu.X)
	return 0
}

func (cpu *CPU) txa() uint8 {
	cpu.Transfer(&cpu.X, &cpu.Accumulator)
	return 0
}

func (cpu *CPU) txs() uint8 {
	cpu.SP = cpu.X
	return 0
}

func (cpu *CPU) tya() uint8 {
	cpu.Transfer(&cpu.Y, &cpu.Accumulator)
	return 0
}

// Not an operator
func (cpu *CPU) nop() uint8 {
	return 0
}
