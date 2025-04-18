package cpu

type Instruction struct {
	Value   int
	Address int
	Cycle   int
}

func (cpu *CPU) imp() uint8 {
	cpu.fetched = cpu.Accumulator
	return 0
}
func (cpu *CPU) zp0() uint8 {
	tmp := cpu.read(cpu.PC)
	cpu.addr_abs = uint16(tmp)
	cpu.PC++
	cpu.addr_abs &= 0x00FF
	return 0
}
func (cpu *CPU) zpy() uint8 {
	tmp := cpu.read(cpu.PC + uint16(cpu.Y))
	cpu.addr_abs = uint16(tmp)
	cpu.PC++
	cpu.addr_abs &= 0x00FF
	return 0
}
func (cpu *CPU) abs() uint8 {
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addr_abs = (high_addr << 8) | low_addr
	return 0
}
func (cpu *CPU) aby() uint8 {
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.addr_abs += uint16(cpu.Y)
	if cpu.addr_abs&0xFF00 != (high_addr << 8) {
		return 1
	}
	return 0
}
func (cpu *CPU) izx() uint8 {
	var t uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var low_addr uint16 = uint16(cpu.read(uint16(t+uint16(cpu.X)) & 0x00FF))
	var high_addr uint16 = uint16(cpu.read(uint16(t+uint16(cpu.X)+1) & 0x00FF))
	cpu.addr_abs = (high_addr << 8) | low_addr
	return 0
}
func (cpu *CPU) imm() uint8 {
	cpu.addr_abs = cpu.PC
	cpu.PC++
	return 0
}
func (cpu *CPU) zpx() uint8 {
	cpu.addr_abs = uint16(cpu.read(cpu.PC + uint16(cpu.Y)))
	cpu.PC++
	cpu.addr_abs &= 0x00FF
	return 0
}
func (cpu *CPU) rel() uint8 {
	cpu.addr_relative = uint16(cpu.read(cpu.PC))
	cpu.PC++
	if cpu.addr_relative&0x80 != 0 {
		cpu.addr_relative |= 0xFF00
	}
	return 0
}
func (cpu *CPU) abx() uint8 {
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.addr_abs += uint16(cpu.X)

	if cpu.addr_abs&0xFF00 != (high_addr << 8) {
		return 1
	}
	return 0
}
func (cpu *CPU) ind() uint8 {
	var high_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var low_addr uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var ptr uint16 = (high_addr << 8) | low_addr
	cpu.addr_abs = uint16(cpu.read(ptr+1)<<8 | cpu.read(ptr))
	return 0
}
func (cpu *CPU) izy() uint8 {
	var t uint16 = uint16(cpu.read(cpu.PC))
	cpu.PC++
	var low_addr uint16 = uint16(cpu.read(t & 0x00FF))
	var high_addr uint16 = uint16(cpu.read(t + 1&0x00FF))
	cpu.addr_abs = (high_addr << 8) | low_addr
	cpu.addr_abs += uint16(cpu.Y)
	if cpu.addr_abs&0x00FF != (high_addr << 8) {
		return 1
	}
	return 0
}

// Opcodes

func adc() {

}

func and() {

}

func asl() {

}

func bcc() {

}

func bcs() {

}

func beq() {

}

func bit() {

}

func bmi() {

}

func bne() {

}

func bpl() {

}

func brk() {

}

func bvc() {

}

func bvs() {

}

func clc() {

}

func cld() {

}

func cli() {

}

func clv() {

}

func cmp() {

}

func cpx() {

}

func cpy() {

}

func dec() {

}

func dex() {

}

func dey() {

}

func eor() {

}

func inc() {

}

func inx() {

}

func iny() {

}

func jmp() {

}

func jsr() {

}

func lda() {

}

func ldx() {

}

func ldy() {

}

func lsr() {

}

func nop() {

}

func ora() {

}

func pha() {

}

func php() {

}

func pla() {

}

func plp() {

}

func rol() {

}

func ror() {

}

func rti() {

}

func rts() {

}

func sbc() {

}

func sec() {

}

func sed() {

}

func sei() {

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
