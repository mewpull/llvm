package ir

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// --- [ Unary instructions ] --------------------------------------------------

// ~~~ [ fneg ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstFNeg is an LLVM IR fneg instruction.
type InstFNeg struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Operand.
	X value.Value // floating-point scalar or floating-point vector

	// extra.

	// (optional) Fast math flags.
	FastMathFlags []enum.FastMathFlag
	// (optional) Metadata.
	Metadata
}

// NewFNeg returns a new fneg instruction based on the given operand.
func NewFNeg(x value.Value) *InstFNeg {
	inst := &InstFNeg{X: x}
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstFNeg) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstFNeg) Type() types.Type {
	return inst.X.Type()
}

// LLString returns the LLVM syntax representation of the instruction.
func (inst *InstFNeg) LLString() string {
	// 'fneg' FastMathFlags=FastMathFlag* X=TypeValue Metadata=(','
	// MetadataAttachment)+?
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	buf.WriteString("fneg")
	for _, flag := range inst.FastMathFlags {
		fmt.Fprintf(buf, " %s", flag)
	}
	fmt.Fprintf(buf, " %s", inst.X)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}
