package proc

import "fmt"
import "bytes"
import sys "golang.org/x/sys/unix"

type Regs struct {
	regs *sys.PtraceRegs
}

func (r *Regs) String() string {
	var buf bytes.Buffer
	var regs = []struct {
		k string
		v uint32
	}{
		{"R0", r.regs.Uregs[0]},
		{"R1", r.regs.Uregs[1]},
		{"R2", r.regs.Uregs[2]},
		{"R3", r.regs.Uregs[3]},
		{"R4", r.regs.Uregs[4]},
		{"R5", r.regs.Uregs[5]},
		{"R6", r.regs.Uregs[6]},
		{"R7", r.regs.Uregs[7]},
		{"R8", r.regs.Uregs[8]},
		{"R9", r.regs.Uregs[9]},
		{"R10 (g)", r.regs.Uregs[10]},
		{"R11", r.regs.Uregs[11]},
		{"R12", r.regs.Uregs[12]},
		{"R13 (SP)", r.regs.Uregs[13]},
		{"R14 (LR)", r.regs.Uregs[14]},
		{"R15 (PC)", r.regs.PC()},
	}
	for _, reg := range regs {
		fmt.Fprintf(&buf, "%8s = %0#16x\n", reg.k, reg.v)
	}
	return buf.String()
}

func (r *Regs) PC() uintptr {
	return uintptr(r.regs.PC())
}

func (r *Regs) SP() uintptr {
	return uintptr(r.regs.Uregs[13])
}

func (r *Regs) CX() uintptr {
	return uintptr(r.regs.Uregs[1])
}

func (r *Regs) TLS() uintptr {
	return uintptr(r.regs.Uregs[10])
}

func (r *Regs) SetPC(thread *Thread, pc uintptr) (err error) {
	r.regs.SetPC(uint64(pc))
	thread.dbp.execPtraceFunc(func() { err = sys.PtraceSetRegs(thread.Id, r.regs) })
	return
}

func registers(thread *Thread) (Registers, error) {
	var (
		regs sys.PtraceRegs
		err  error
	)
	thread.dbp.execPtraceFunc(func() { err = sys.PtraceGetRegs(thread.Id, &regs) })
	if err != nil {
		return nil, err
	}
	return &Regs{&regs}, nil
}
