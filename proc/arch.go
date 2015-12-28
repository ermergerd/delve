package proc

import (
	"encoding/binary"
	"runtime"
)

type Arch interface {
	SetGStructOffset(ver GoVersion, iscgo bool)
	DecodePtr(b []byte) uintptr
	EncodePtr(p uintptr, b []byte)
	PtrSize() int
	BreakpointInstruction() []byte
	BreakpointSize() int
	GStructOffset() uint64
	HardwareBreakpointUsage() []bool
	SetHardwareBreakpointUsage(int, bool)
	
	// Also implement the byte order interface to encode and decode multi-byte values
	binary.ByteOrder
}

type AMD64 struct {
	ptrSize                 int
	breakInstruction        []byte
	breakInstructionLen     int
	gStructOffset           uint64
	hardwareBreakpointUsage []bool
	
	binary.ByteOrder
}

func AMD64Arch() *AMD64 {
	var breakInstr = []byte{0xCC}

	return &AMD64{
		ptrSize:                 8,
		breakInstruction:        breakInstr,
		breakInstructionLen:     len(breakInstr),
		hardwareBreakpointUsage: make([]bool, 4),
		ByteOrder:               binary.LittleEndian,
	}
}

func (a *AMD64) SetGStructOffset(ver GoVersion, isextld bool) {
	switch runtime.GOOS {
	case "darwin":
		a.gStructOffset = 0x8a0
	case "linux":
		a.gStructOffset = 0xfffffffffffffff0
		if isextld || ver.AfterOrEqual(GoVersion{1, 5, -1, 2, 0}) || ver.IsDevel() {
			a.gStructOffset += 8
		}
	}
}

func (a *AMD64) DecodePtr(b []byte) uintptr {
	return uintptr(binary.LittleEndian.Uint64(b))
}

func (a *AMD64) EncodePtr(p uintptr, b []byte) {
	binary.LittleEndian.PutUint64(b, uint64(p))
}

func (a *AMD64) PtrSize() int {
	return a.ptrSize
}

func (a *AMD64) BreakpointInstruction() []byte {
	return a.breakInstruction
}

func (a *AMD64) BreakpointSize() int {
	return a.breakInstructionLen
}

func (a *AMD64) GStructOffset() uint64 {
	return a.gStructOffset
}

func (a *AMD64) HardwareBreakpointUsage() []bool {
	return a.hardwareBreakpointUsage
}

func (a *AMD64) SetHardwareBreakpointUsage(reg int, set bool) {
	a.hardwareBreakpointUsage[reg] = set
}

type ARM struct {
	ptrSize                 int
	breakInstruction        []byte
	breakInstructionLen     int
	gStructOffset           uint64
	hardwareBreakpointUsage []bool
	
	binary.ByteOrder
}

func ARMArch() *ARM {
	var breakInstr = []byte{0xFE, 0xDE, 0xFF, 0xE7}
	//var breakInstr = []byte{0xE1, 0x20, 0x00, 0x70}

	return &ARM{
		ptrSize:                 4,
		breakInstruction:        breakInstr,
		breakInstructionLen:     len(breakInstr),
		hardwareBreakpointUsage: make([]bool, 4),
		ByteOrder:               binary.LittleEndian,
	}
}

func (a *ARM) SetGStructOffset(ver GoVersion, isextld bool) {
	switch runtime.GOOS {
	case "darwin":
		a.gStructOffset = 0x8a0
	case "linux":
		a.gStructOffset = 0xfffffffffffffff0
		if isextld || ver.AfterOrEqual(GoVersion{1, 5, -1, 2, 0}) || ver.IsDevel() {
			a.gStructOffset += 8
		}
	}
}

func (a *ARM) DecodePtr(b []byte) uintptr {
	return uintptr(binary.LittleEndian.Uint32(b))
}

func (a *ARM) EncodePtr(p uintptr, b []byte) {
	binary.LittleEndian.PutUint32(b, uint32(p))
}

func (a *ARM) PtrSize() int {
	return a.ptrSize
}

func (a *ARM) BreakpointInstruction() []byte {
	return a.breakInstruction
}

func (a *ARM) BreakpointSize() int {
	return a.breakInstructionLen
}

func (a *ARM) GStructOffset() uint64 {
	return a.gStructOffset
}

func (a *ARM) HardwareBreakpointUsage() []bool {
	return a.hardwareBreakpointUsage
}

func (a *ARM) SetHardwareBreakpointUsage(reg int, set bool) {
	a.hardwareBreakpointUsage[reg] = set
}
