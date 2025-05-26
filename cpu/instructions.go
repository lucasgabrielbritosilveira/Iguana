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
	cpu.addr_abs = uint16(cpu.read(cpu.PC))
	cpu.addr_abs &= 0x00FF
	cpu.PC++
	return 0
}

func (cpu *CPU) zpy() uint8 {
	cpu.addr_abs = uint16(cpu.read(cpu.PC + uint16(cpu.Y)))
	cpu.addr_abs &= 0x00FF
	cpu.PC++
	return 0
}

func (cpu *CPU) abs() uint8 {
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addr_abs = (high_addr << 8) | low_addr
	return 0
}

func (cpu *CPU) aby() uint8 {
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.addr_abs += uint16(cpu.Y)
	if cpu.addr_abs&0xFF00 != (high_addr << 8) {
		return 1
	}
	return 0
}

func (cpu *CPU) izx() uint8 {
	var temp uint16 = uint16(cpu.read(cpu.PC))
	var low_addr uint16 = uint16(cpu.read(uint16(temp+uint16(cpu.X)) & 0x00FF))
	var high_addr uint16 = uint16(cpu.read(uint16(temp+uint16(cpu.X)+1) & 0x00FF))
	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.PC++
	return 0
}

func (cpu *CPU) imm() uint8 {
	cpu.addr_abs = cpu.PC
	cpu.PC++
	return 0
}

func (cpu *CPU) zpx() uint8 {
	cpu.addr_abs = uint16(cpu.read(cpu.PC + uint16(cpu.X)))
	cpu.addr_abs &= 0x00FF
	cpu.PC++
	return 0
}

func (cpu *CPU) rel() uint8 {
	cpu.addr_relative = uint16(cpu.read(cpu.PC))
	if cpu.addr_relative&0x80 != 0 {
		cpu.addr_relative |= 0xFF00
	}
	cpu.PC++
	return 0
}

func (cpu *CPU) abx() uint8 {
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.addr_abs += uint16(cpu.X)
	if cpu.addr_abs&0xFF00 != (high_addr << 8) {
		return 1
	}
	return 0

}
func (cpu *CPU) ind() uint8 {
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var ptr uint16 = (high_addr << 8) | low_addr
	if low_addr == 0x00FF {
		cpu.addr_abs = uint16(cpu.read(ptr&0xFF00))<<8 | uint16(cpu.read(ptr))
	} else {
		cpu.addr_abs = uint16(cpu.read(ptr+1))<<8 | uint16(cpu.read(ptr))
	}
	return 0
}

func (cpu *CPU) izy() uint8 {
	var tmp uint16 = uint16(cpu.read(cpu.PC))
	var low_addr uint16 = uint16(cpu.read(tmp & 0x00FF))
	var high_addr uint16 = uint16(cpu.read(tmp + 1&0x00FF))

	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.addr_abs += uint16(cpu.Y)

	cpu.PC++
	if cpu.addr_abs&0x00FF != (high_addr << 8) {
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
	var value uint16 = uint16(cpu.fetched) ^ 0x00FF
	cpu.Add(value)
	return 1
}

func (cpu *CPU) dec() uint8 {
	cpu.fetch()
	var tmp uint8 = cpu.fetched - 1
	cpu.write(cpu.addr_abs, tmp&0x00FF)
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
	var tmp uint8 = cpu.fetched + 1
	cpu.write(cpu.addr_abs, tmp&0x00FF)
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

func eor() {

}

func ora() {

}

func bit() {

}

func rol() {

}

func ror() {

}

func lsr() {

}

// Branchs and Jumps

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

func pha() {

}

func php() {

}

func pla() {

}

func plp() {

}

// Data

func lda() {

}

func ldx() {

}

func ldy() {

}

func sta() {

}

func stx() {

}

func sty() {

}

func tax() {

}

func tay() {

}

func tsx() {

}

func txa() {

}

func txs() {

}

func tya() {

}

// Not an operator
func nop() {

}
