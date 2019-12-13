package ir

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// --- [ Vector instructions ] -------------------------------------------------

// ~~~ [ extractelement ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstExtractElement is an LLVM IR extractelement instruction.
type InstExtractElement struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Vector.
	X value.Value
	// Element index.
	Index value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewExtractElement returns a new extractelement instruction based on the given
// vector and element index.
func NewExtractElement(x, index value.Value) *InstExtractElement {
	inst := &InstExtractElement{X: x, Index: index}
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstExtractElement) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstExtractElement) Type() types.Type {
	t, ok := inst.X.Type().(*types.VectorType)
	if !ok {
		panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", inst.X.Type()))
	}
	return t.ElemType
}

// LLString returns the LLVM syntax representation of the instruction.
func (inst *InstExtractElement) LLString() string {
	// 'extractelement' X=TypeValue ',' Index=TypeValue Metadata=(','
	// MetadataAttachment)+?
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "extractelement %s, %s", inst.X, inst.Index)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ insertelement ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstInsertElement is an LLVM IR insertelement instruction.
type InstInsertElement struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Vector.
	X value.Value
	// Element to insert.
	Elem value.Value
	// Element index.
	Index value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewInsertElement returns a new insertelement instruction based on the given
// vector, element and element index.
func NewInsertElement(x, elem, index value.Value) *InstInsertElement {
	inst := &InstInsertElement{X: x, Elem: elem, Index: index}
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstInsertElement) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstInsertElement) Type() types.Type {
	t, ok := inst.X.Type().(*types.VectorType)
	if !ok {
		panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", inst.X.Type()))
	}
	return t
}

// LLString returns the LLVM syntax representation of the instruction.
func (inst *InstInsertElement) LLString() string {
	// 'insertelement' X=TypeValue ',' Elem=TypeValue ',' Index=TypeValue
	// Metadata=(',' MetadataAttachment)+?
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "insertelement %s, %s, %s", inst.X, inst.Elem, inst.Index)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}

// ~~~ [ shufflevector ] ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// InstShuffleVector is an LLVM IR shufflevector instruction.
type InstShuffleVector struct {
	// Name of local variable associated with the result.
	LocalIdent
	// Vectors.
	X, Y value.Value
	// Shuffle mask.
	Mask value.Value

	// extra.

	// (optional) Metadata.
	Metadata
}

// NewShuffleVector returns a new shufflevector instruction based on the given
// vectors and shuffle mask.
func NewShuffleVector(x, y, mask value.Value) *InstShuffleVector {
	inst := &InstShuffleVector{X: x, Y: y, Mask: mask}
	return inst
}

// String returns the LLVM syntax representation of the instruction as a
// type-value pair.
func (inst *InstShuffleVector) String() string {
	return fmt.Sprintf("%s %s", inst.Type(), inst.Ident())
}

// Type returns the type of the instruction.
func (inst *InstShuffleVector) Type() types.Type {
	// Cache type if not present.
	xType, ok := inst.X.Type().(*types.VectorType)
	if !ok {
		panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", inst.X.Type()))
	}
	maskType, ok := inst.Mask.Type().(*types.VectorType)
	if !ok {
		panic(fmt.Errorf("invalid vector type; expected *types.VectorType, got %T", inst.Mask.Type()))
	}
	return types.NewVector(maskType.Len, xType.ElemType)
}

// LLString returns the LLVM syntax representation of the instruction.
func (inst *InstShuffleVector) LLString() string {
	// 'shufflevector' X=TypeValue ',' Y=TypeValue ',' Mask=TypeValue
	// Metadata=(',' MetadataAttachment)+?
	buf := &strings.Builder{}
	fmt.Fprintf(buf, "%s = ", inst.Ident())
	fmt.Fprintf(buf, "shufflevector %s, %s, %s", inst.X, inst.Y, inst.Mask)
	for _, md := range inst.Metadata {
		fmt.Fprintf(buf, ", %s", md)
	}
	return buf.String()
}
