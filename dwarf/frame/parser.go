// Package frame contains data structures and
// related functions for parsing and searching
// through Dwarf .debug_frame data.
package frame

import (
	"bytes"
	"encoding/binary"

	"github.com/derekparker/delve/dwarf/util"
)

type parsefunc func(*parseContext) parsefunc

type parseContext struct {
	buf     *bytes.Buffer
	entries FrameDescriptionEntries
	common  *CommonInformationEntry
	frame   *FrameDescriptionEntry
	length  uint32
	order   binary.ByteOrder
	ptrSize uint
}

// Parse takes in data (a byte slice) and returns a slice of
// commonInformationEntry structures. Each commonInformationEntry
// has a slice of frameDescriptionEntry structures.
func Parse(data []byte, order binary.ByteOrder, ptrSize uint) FrameDescriptionEntries {
	var (
		buf  = bytes.NewBuffer(data)
		pctx = &parseContext{buf: buf, entries: NewFrameIndex(), order: order, ptrSize: ptrSize}
	)

	for fn := parselength; buf.Len() != 0; {
		fn = fn(pctx)
	}

	return pctx.entries
}

func cieEntry(data []byte) bool {
	return bytes.Equal(data, []byte{0xff, 0xff, 0xff, 0xff})
}

func parselength(ctx *parseContext) parsefunc {
	var data = ctx.buf.Next(8)

	ctx.length = binary.LittleEndian.Uint32(data[:4]) - 4 // take off the length of the CIE id / CIE pointer.

	if cieEntry(data[4:]) {
		ctx.common = &CommonInformationEntry{Length: ctx.length}
		return parseCIE
	}

	ctx.frame = &FrameDescriptionEntry{Length: ctx.length, CIE: ctx.common}
	return parseFDE
}

func parseFDE(ctx *parseContext) parsefunc {
	r := ctx.buf.Next(int(ctx.length))
	
	if ctx.ptrSize == 4 {
		ctx.frame.begin = uintptr(ctx.order.Uint32(r[:4]))
		ctx.frame.end = uintptr(ctx.order.Uint32(r[4:8]))
		
		// The rest of this entry consists of the instructions
		// so we can just grab all of the data from the buffer
		// cursor to length.
		ctx.frame.Instructions = r[8:]
	} else if ctx.ptrSize == 8 {
		ctx.frame.begin = uintptr(ctx.order.Uint64(r[:8]))
		ctx.frame.end = uintptr(ctx.order.Uint64(r[8:16]))
		
		// The rest of this entry consists of the instructions
		// so we can just grab all of the data from the buffer
		// cursor to length.
		ctx.frame.Instructions = r[16:]
	}

	// Insert into the tree after setting address range begin
	// otherwise compares won't work.
	ctx.entries = append(ctx.entries, ctx.frame)
	ctx.length = 0

	return parselength
}

func parseCIE(ctx *parseContext) parsefunc {
	data := ctx.buf.Next(int(ctx.length))
	buf := bytes.NewBuffer(data)
	// parse version
	ctx.common.Version = data[0]

	// parse augmentation
	ctx.common.Augmentation, _ = util.ParseString(buf)

	// parse code alignment factor
	ctx.common.CodeAlignmentFactor, _ = util.DecodeULEB128(buf)

	// parse data alignment factor
	ctx.common.DataAlignmentFactor, _ = util.DecodeSLEB128(buf)

	// parse return address register
	returnAddressRegister, _ := util.DecodeULEB128(buf)
	ctx.common.ReturnAddressRegister = uintptr(returnAddressRegister)

	// parse initial instructions
	// The rest of this entry consists of the instructions
	// so we can just grab all of the data from the buffer
	// cursor to length.
	ctx.common.InitialInstructions = buf.Bytes() //ctx.buf.Next(int(ctx.length))
	ctx.length = 0

	return parselength
}
