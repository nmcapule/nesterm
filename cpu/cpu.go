// Package cpu implements the 6502 microprocessor.
package cpu

// 6502 Addressing Modes.
var (
  PLACE = iota // Placeholder for Unknown Op Codes
  ABSOL        // Absolute
  ABSOX        // Absolute,X
  ABSOY        // Absolute,Y
  ACCUM        // Accumulator
  IDIDX        // Indirect Indexed
  IDREC        // Indirect
  IDXDI        // Indexed Indirect
  IMMED        // Immediate
  IMPLY        // Implied
  RELAT        // Relative
  ZEROP        // Zero Page
  ZEROX        // Zero Page,X
  ZEROY        // Zero Page,Y
)

var modetable = []string{
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

// Cpu is an implementation of the 6502 microprocessor.
type Cpu struct {
	memory     [0xFFFF]uint8 // Addressable memory in 6502
	pc         uint16        // Program Counter
	sp         uint8         // Stack Pointer
	a, x, y, p uint8         // Accumulator; X, Y Register; Processor Status
}
