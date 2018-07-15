// Package cpu implements the 6502 microprocessor.
package cpu

import (
	"fmt"
)

type addrMode int8

// 6502 Addressing Modes.
const (
	PLACE addrMode = iota // Placeholder for Unknown Op Codes
	ABSOL                 // Absolute
	ABSOX                 // Absolute,X
	ABSOY                 // Absolute,Y
	ACCUM                 // Accumulator
	IDIDX                 // Indirect Indexed
	IDREC                 // Indirect
	IDXDI                 // Indexed Indirect
	IMMED                 // Immediate
	IMPLY                 // Implied
	RELAT                 // Relative
	ZEROP                 // Zero Page
	ZEROX                 // Zero Page,X
	ZEROY                 // Zero Page,Y
)

var modesize = []uint16{
	1, // PLACE
	3, // ABSOL
	3, // ABSOX
	3, // ABSOY
	1, // ACCUM
	2, // IDIDX
	3, // IDREC
	2, // IDXDI
	3, // IMMED
	1, // IMPLY
	2, // RELAT
	2, // ZEROP
	2, // ZEROX
	2, // ZEROY
	// More
}

var modetable = []addrMode{
	//00   0x01   0x02   0x03   0x04   0x05   0x06   0x07
	//08   0x09   0x0a   0x0b   0x0c   0x0d   0x0e   0x0f
	IMPLY, IDXDI, PLACE, PLACE, PLACE, ZEROP, ZEROP, PLACE, // 0x00
	IMPLY, IMMED, ACCUM, PLACE, PLACE, ABSOL, ABSOL, PLACE, // 0x08
	RELAT, IDIDX, PLACE, PLACE, PLACE, ZEROX, ZEROX, PLACE, // 0x10
	IMPLY, ABSOY, PLACE, PLACE, PLACE, ABSOX, ABSOX, PLACE, // 0x18
	ABSOL, IDXDI, PLACE, PLACE, ZEROP, ZEROP, ZEROP, PLACE, // 0x20
	IMPLY, IMMED, ACCUM, PLACE, ABSOL, ABSOL, ABSOL, PLACE, // 0x28
	RELAT, IDIDX, PLACE, PLACE, PLACE, ZEROX, ZEROX, PLACE, // 0x30
	IMPLY, ABSOY, PLACE, PLACE, PLACE, ABSOX, ABSOX, PLACE, // 0x38
	IMPLY, IDXDI, PLACE, PLACE, PLACE, ZEROP, ZEROP, PLACE, // 0x40
	IMPLY, IMMED, ACCUM, PLACE, ABSOL, ABSOL, ABSOL, PLACE, // 0x48
	RELAT, IDIDX, PLACE, PLACE, PLACE, ZEROX, ZEROX, PLACE, // 0x50
	IMPLY, ABSOY, PLACE, PLACE, PLACE, ABSOX, ABSOX, PLACE, // 0x58
	IMPLY, IDXDI, PLACE, PLACE, PLACE, ZEROP, ZEROP, PLACE, // 0x60
	IMPLY, IMMED, ACCUM, PLACE, IDREC, ABSOL, ABSOX, PLACE, // 0x68
	RELAT, IDIDX, PLACE, PLACE, PLACE, ZEROX, ZEROX, PLACE, // 0x70
	IMPLY, ABSOY, PLACE, PLACE, PLACE, ABSOX, ABSOL, PLACE, // 0x78
	//00   0x01   0x02   0x03   0x04   0x05   0x06   0x07
	//08   0x09   0x0a   0x0b   0x0c   0x0d   0x0e   0x0f
	PLACE, IDXDI, PLACE, PLACE, ZEROP, ZEROP, ZEROP, PLACE, // 0x80
	IMPLY, PLACE, IMPLY, PLACE, ABSOL, ABSOL, ABSOL, PLACE, // 0x88
	RELAT, IDIDX, PLACE, PLACE, ZEROX, ZEROX, ZEROY, PLACE, // 0x90
	IMPLY, ABSOY, IMPLY, PLACE, PLACE, ABSOX, PLACE, PLACE, // 0x98
	IMMED, IDXDI, IMMED, PLACE, ZEROP, ZEROP, ZEROP, PLACE, // 0xa0
	IMPLY, IMMED, IMPLY, PLACE, ABSOL, ABSOL, ABSOL, PLACE, // 0xa8
	RELAT, IDIDX, PLACE, PLACE, ZEROX, ZEROX, ZEROY, PLACE, // 0xb0
	IMPLY, ABSOY, IMPLY, PLACE, ABSOX, ABSOX, ABSOY, PLACE, // 0xb8
	IMMED, IDXDI, PLACE, PLACE, ZEROP, ZEROP, ZEROP, PLACE, // 0xc0
	IMPLY, IMMED, IMPLY, PLACE, ABSOL, ABSOL, ABSOL, PLACE, // 0xc8
	RELAT, IDIDX, PLACE, PLACE, PLACE, ZEROX, ZEROX, PLACE, // 0xd0
	IMPLY, ABSOY, PLACE, PLACE, PLACE, ABSOX, ABSOX, PLACE, // 0xd8
	IMMED, IDXDI, PLACE, PLACE, ZEROP, ZEROP, ZEROP, PLACE, // 0xe0
	IMPLY, IMMED, IMPLY, PLACE, ABSOL, ABSOL, ABSOL, PLACE, // 0xe8
	RELAT, IDIDX, PLACE, PLACE, PLACE, ZEROX, ZEROX, PLACE, // 0xf0
	IMPLY, ABSOY, PLACE, PLACE, PLACE, ABSOX, ABSOX, PLACE, // 0xf8
	//00   0x01   0x02   0x03   0x04   0x05   0x06   0x07
	//08   0x09   0x0a   0x0b   0x0c   0x0d   0x0e   0x0f
}

// TODO(nmcapule): These aren't really accurate.
var cycles = []uint8{
	//01 02 03 04 05 06 07
	//09 0a 0b 0c 0d 0e 0f
	2, 2, 2, 2, 2, 2, 2, 2, // 0x00
	2, 2, 2, 2, 2, 2, 2, 2, // 0x08
	2, 2, 2, 2, 2, 2, 2, 2, // 0x10
	2, 2, 2, 2, 2, 2, 2, 2, // 0x18
	2, 2, 2, 2, 2, 2, 2, 2, // 0x20
	2, 2, 2, 2, 2, 2, 2, 2, // 0x28
	2, 2, 2, 2, 2, 2, 2, 2, // 0x30
	2, 2, 2, 2, 2, 2, 2, 2, // 0x38
	2, 2, 2, 2, 2, 2, 2, 2, // 0x40
	2, 2, 2, 2, 2, 2, 2, 2, // 0x48
	2, 2, 2, 2, 2, 2, 2, 2, // 0x50
	2, 2, 2, 2, 2, 2, 2, 2, // 0x58
	2, 2, 2, 2, 2, 2, 2, 2, // 0x60
	2, 2, 2, 2, 2, 2, 2, 2, // 0x68
	2, 2, 2, 2, 2, 2, 2, 2, // 0x70
	2, 2, 2, 2, 2, 2, 2, 2, // 0x78
	//01 02 03 04 05 06 07
	//09 0a 0b 0c 0d 0e 0f
	2, 2, 2, 2, 2, 2, 2, 2, // 0x80
	2, 2, 2, 2, 2, 2, 2, 2, // 0x88
	2, 2, 2, 2, 2, 2, 2, 2, // 0x90
	2, 2, 2, 2, 2, 2, 2, 2, // 0x98
	2, 2, 2, 2, 2, 2, 2, 2, // 0xa0
	2, 2, 2, 2, 2, 2, 2, 2, // 0xa8
	2, 2, 2, 2, 2, 2, 2, 2, // 0xb0
	2, 2, 2, 2, 2, 2, 2, 2, // 0xb8
	2, 2, 2, 2, 2, 2, 2, 2, // 0xc0
	2, 2, 2, 2, 2, 2, 2, 2, // 0xc8
	2, 2, 2, 2, 2, 2, 2, 2, // 0xd0
	2, 2, 2, 2, 2, 2, 2, 2, // 0xd8
	2, 2, 2, 2, 2, 2, 2, 2, // 0xe0
	2, 2, 2, 2, 2, 2, 2, 2, // 0xe8
	2, 2, 2, 2, 2, 2, 2, 2, // 0xf0
	2, 2, 2, 2, 2, 2, 2, 2, // 0xf8
	//01 02 03 04 05 06 07
	//09 0a 0b 0c 0d 0e 0f
}

// http://www.emulator101.com/reference/6502-reference.html
var optable = []string{
	//00   0x01   0x02   0x03   0x04   0x05   0x06   0x07
	//08   0x09   0x0a   0x0b   0x0c   0x0d   0x0e   0x0f
	"BRK", "ORA", "___", "___", "___", "ORA", "ASL", "___", // 0x00
	"PHP", "ORA", "ASL", "___", "___", "ORA", "ASL", "___", // 0x08
	"BPL", "ORA", "___", "___", "___", "ORA", "ASL", "___", // 0x10
	"CLC", "ORA", "___", "___", "___", "ORA", "ASL", "___", // 0x18
	"JSR", "AND", "___", "___", "BIT", "AND", "ROL", "___", // 0x20
	"PLP", "AND", "ROL", "___", "BIT", "AND", "ROL", "___", // 0x28
	"BMI", "AND", "___", "___", "___", "AND", "ROL", "___", // 0x30
	"SEC", "AND", "___", "___", "___", "AND", "ROL", "___", // 0x38
	"RTI", "EOR", "___", "___", "___", "EOR", "LSR", "___", // 0x40
	"PHA", "EOR", "LSR", "___", "JMP", "EOR", "LSR", "___", // 0x48
	"BVC", "EOR", "___", "___", "___", "EOR", "LSR", "___", // 0x50
	"CLI", "EOR", "___", "___", "___", "EOR", "LSR", "___", // 0x58
	"RTS", "ADC", "___", "___", "___", "ADC", "ROR", "___", // 0x60
	"PLA", "ADC", "ROR", "___", "JMP", "ADC", "ROR", "___", // 0x68
	"BVS", "ADC", "___", "___", "___", "ADC", "ROR", "___", // 0x70
	"SEI", "ADC", "___", "___", "___", "ADC", "ROR", "___", // 0x78
	//00   0x01   0x02   0x03   0x04   0x05   0x06   0x07
	//08   0x09   0x0a   0x0b   0x0c   0x0d   0x0e   0x0f
	"___", "STA", "___", "___", "STY", "STA", "STX", "___", // 0x80
	"DEY", "___", "TXA", "___", "STY", "STA", "STX", "___", // 0x88
	"BCC", "STA", "___", "___", "STY", "STA", "STX", "___", // 0x90
	"TYA", "STA", "TXS", "___", "___", "STA", "___", "___", // 0x98
	"LDY", "LDA", "LDX", "___", "LDY", "LDA", "LDX", "___", // 0xa0
	"TAY", "LDA", "TAX", "___", "LDY", "LDA", "LDX", "___", // 0xa8
	"BCS", "LDA", "___", "___", "LDY", "LDA", "LDX", "___", // 0xb0
	"CLV", "LDA", "TSX", "___", "LDY", "LDA", "LDX", "___", // 0xb8
	"CPY", "CMP", "___", "___", "CPY", "CMP", "DEC", "___", // 0xc0
	"INY", "CMP", "DEX", "___", "CPY", "CMP", "DEC", "___", // 0xc8
	"BNE", "CMP", "___", "___", "___", "CMP", "DEC", "___", // 0xd0
	"CLD", "CMP", "___", "___", "___", "CMP", "DEC", "___", // 0xd8
	"CPX", "SBC", "___", "___", "CPX", "SBC", "INC", "___", // 0xe0
	"INX", "SBC", "NOP", "___", "CPX", "SBC", "INC", "___", // 0xe8
	"BEQ", "SBC", "___", "___", "___", "SBC", "INC", "___", // 0xf0
	"SED", "SBC", "___", "___", "___", "SBC", "INC", "___", // 0xf8
	//00   0x01   0x02   0x03   0x04   0x05   0x06   0x07
	//08   0x09   0x0a   0x0b   0x0c   0x0d   0x0e   0x0f
}

// Memory is an implementation of the 6502 addressable memory.
//
// Memory Map:
//  0x0000..0x00FF    - Zero Page
//    0x0000..0x1FFF  - Mirror (See Set 1)
//  0x2000..0x401F    - IO Register
//    0x2000..0x3FFF  - Mirror (See Set 2)
//  0x4020..0x5FFF    - Expansion ROM
//  0x6000..0x7FFF    - SRAM
//  0x8000..0xFFFF    - PRG ROM
//    0x8000..0xBFFF  - 16K Memory Bank 1
//    0xC000..0xFFFF  - 16K Memory Bank 2; Mirrors Bank 1 depending on cartridge (See Set 3)
//
// Mirrors
//  Set 1) 0x0800 increments
//    0x0000..0x07FF
//    0x0800..0x0FFF
//    0x1000..0x17FF
//    0x1800..0x1FFF
//  Set 2) 0x0008 increments
//    0x2000..0x2007
//    0x2008..0x200F
//    0x2010..0x2017
//    0x2018..0x201F
//    ...
//    0x3FF7..0x3FFF
//  Set 3) 0x4000 increments (optional)
//    0x8000..0xBFFF
//    0xC000..0xFFFF
type Memory struct {
	m         [0xFFFF]uint8
	mirrorPRG bool // Set if want to mirror PRG 1 and PRG 2
}

// Computes the effective address of the given absolute address.
func (m *Memory) faddr(addr uint16) uint16 {
	if addr >= 0x0000 && addr <= 0x1FFF {
		return (addr & 0x0800)
	}
	if addr >= 0x2000 && addr <= 0x3FFF {
		return (addr & 0x0008) + 0x2000
	}
	if addr >= 0x8000 && addr <= 0xFFFF && m.mirrorPRG {
		return (addr & 0x4000) + 0x8000
	}
	return addr
}

// Dump returns a whole copy of the virtual memory.
func (m *Memory) Dump() [0xFFFF]uint8 {
	var vm [0xFFFF]uint8
	for i := 0; i <= 0xFFFF; i++ {
		vm[i] = m.Get(uint16(i))
	}
	return vm
}

// Set sets the byte on the given memory address.
// Returns true if page boundary crossed.
func (m *Memory) Set(addr uint16, v uint8) bool {
	m.m[m.faddr(addr)] = v
	return false
}

// Get gets the byte on the given memory address.
func (m *Memory) Get(addr uint16) uint8 {
	return m.m[m.faddr(addr)]
}

// Get16 gets a 16bit from the given memory address.
func (m *Memory) Get16(addr uint16) uint16 {
	hi := uint16(m.Get(addr + 1))
	lo := uint16(m.Get(addr))
	return hi<<8 | lo
}

// Processor status flags.
// http://nesdev.com/6502.txt
const (
	FlagCarry       uint8 = 0x01
	FlagZero              = 0x02
	FlagNoInterrupt       = 0x04
	FlagDecimal           = 0x08
	FlagBreak             = 0x10
	FlagUnused            = 0x20
	FlagOverflow          = 0x40
	FlagSign              = 0x80 // aka. Negative Flag
)

// CPU is an implementation of the 6502 microprocessor.
type CPU struct {
	memory     Memory // Addressable memory in 6502
	pc         uint16 // Program Counter
	sp         uint8  // Stack Pointer
	a, x, y, p uint8  // Accumulator; X, Y Register; Processor Status
}

// String implements the Stringer interface.
func (c CPU) String() string {
	s := `6502 (2A03)
  Program Counter: %x
  Stack Pointer:   %x
  Registers:
    A = %x
    X = %x
    Y = %x
  Status Flags:
    Carry             = %t
    Zero              = %t
    Interrupt Disable = %t
    Decimal Mode      = %t
    Break             = %t
    (unused)          = %t
    Overflow          = %t
    Negative          = %t`

	return fmt.Sprintf(s, c.pc, c.sp, c.a, c.x, c.y,
		c.isflag(FlagCarry), c.isflag(FlagZero),
		c.isflag(FlagNoInterrupt), c.isflag(FlagDecimal),
		c.isflag(FlagBreak), c.isflag(FlagUnused),
		c.isflag(FlagOverflow), c.isflag(FlagSign))
}

// Step executes a single cycle of the CPU.
func (c *CPU) Step() {
	op := c.memory.Get(c.pc)
	addr := c.pc + 1
	mode := modetable[op]

	// Move counter at the end of instruction and args
	c.pc += modesize[mode]

	addr, pbc := c.calcaddr(addr, mode)
	if pbc {
		fmt.Println("Page boundary crossed")
	}

	// TODO(nmcapule): Write me
	fmt.Println(addr)
}

// Returns true if processor status flag is set.
func (c *CPU) isflag(flag uint8) bool {
	return c.p&flag != 0
}

// Set processor status flag.
func (c *CPU) flag(flag uint8) {
	c.p |= flag
}

// Clear processor status flag.
func (c *CPU) unflag(flag uint8) {
	c.p &= ^flag
}

// Assign value to processor status flag.
func (c *CPU) setflag(flag uint8, v bool) {
	if v {
		c.flag(flag)
	} else {
		c.unflag(flag)
	}
}

func (c *CPU) calcflags(result int16, mask uint8) {
	// TODO(nmcapule): Im not sure about this one though.
	if mask&FlagCarry != 0 {
		c.setflag(FlagCarry, result&0xF00 != 0)
	}
	if mask&FlagZero != 0 {
		c.setflag(FlagZero, result == 0)
	}
	// TODO(nmcapule): Im not sure about this one though.
	if mask&FlagOverflow != 0 {
		c.setflag(FlagZero, result > 0x0FF)
	}
	// Sign bit is just the most significant bit in a byte.
	if mask&FlagSign != 0 {
		c.setflag(FlagSign, result&0x80 != 0)
	}
}

// Calculates the address to be used by the op. It is assumed that the program
// counter has already moved before this method is called.
//
// Returns the address and whether addressing mode caused a page boundary crossing.
func (c *CPU) calcaddr(addr uint16, mode addrMode) (uint16, bool) {
	switch mode {
	case ABSOL: // Absolute
		return c.memory.Get16(addr), false
	case ABSOX: // Absolute + X
		base := c.memory.Get16(addr)
		addr = base + uint16(c.x)
		return addr, crossPage(addr, base)
	case ABSOY: // Absolute + Y
		base := c.memory.Get16(addr)
		addr = base + uint16(c.y)
		return addr, crossPage(addr, base)
	case IMMED: // Immediate
		return addr, false
	case IDREC:
		return c.memory.Get16(c.memory.Get16(addr)), false
	case IDXDI: // Indexed Indirect
		// TODO(ncapule): Check out fogleman/nes bug
		return c.memory.Get16(uint16(c.memory.Get(addr) + c.x)), false
	case IDIDX: // Indirect Indexed
		// TODO(ncapule): Check out fogleman/nes bug
		base := c.memory.Get16(uint16(c.memory.Get(addr)))
		addr = base + uint16(c.y)
		return addr, crossPage(addr, base)
	case RELAT: // Relative
		offset := uint16(c.memory.Get(addr))
		// TODO(ncapule): Why?
		if offset >= 0x80 {
			return c.pc + offset - 0x100, false
		}
		return c.pc + offset, false
	case ZEROP: // Zero Page
		return uint16(c.memory.Get(addr)), false
	case ZEROX: // Zero Page + X
		return uint16(c.memory.Get(addr) + c.x), false
	case ZEROY: // Zero Page + Y
		return uint16(c.memory.Get(addr) + c.y), false
	default:
		// ACCUM, IMPLY
		return 0, false
	}
}

func (c *CPU) adc(x uint8) int16 {
	// TODO(nmcapule)

	c.calcflags(int16(c.a), FlagSign|FlagZero|FlagCarry|FlagOverflow)

	return int16(0)
}

func (c *CPU) and(x uint8) int16 {
	c.a &= x

	c.calcflags(int16(c.a), FlagSign|FlagZero)

	return int16(c.a)
}

func (c *CPU) clc() {
	c.unflag(FlagCarry)
}

func (c *CPU) cld() {
	c.unflag(FlagDecimal)
}

func (c *CPU) cli() {
	c.unflag(FlagNoInterrupt)
}

func (c *CPU) clv() {
	c.unflag(FlagOverflow)
}

// Checks if addr a and b differ pages (page is per 256 bytes).
func crossPage(a, b uint16) bool {
	return a&0xFF00 != b&0xFF00
}
