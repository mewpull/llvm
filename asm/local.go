// Problems to solve.
//
// phi instructions can reference local variables defined in basic blocks not
// yet visited when translating basic blocks in linear order.
//
// Terminator instructions can reference basic blocks not yet visited when
// translating basic blocks in linear order.
//
// The function parameters, basic blocks and local variables (produced by the
// result of instructions) of a function may be unnamed. They are assigned the
// first unused local ID (e.g. %42) when traversing the body of the function in
// linear order; where function parameters are assigned first, then for each
// basic block, assign an ID to the basic block and then to the result of its
// instructions. Note, instructions that produce void results are ignored.
// Non-value instructions (e.g. store) are always ignored. Notably, the call
// instruction may be ignored if the callee has a void return.

// TODO: make concurrent :)

package asm

import (
	"github.com/llir/ll/ast"
	"github.com/llir/llvm/internal/enc"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/pkg/errors"
)

// funcGen is a generator for a given IR function.
type funcGen struct {
	// Module generator.
	gen *generator
	// LLVM IR function being generated.
	f *ir.Func
	// locals maps from local identifier (without '%' prefix) to corresponding IR
	// value.
	locals map[ir.LocalIdent]value.Value
	// isVoid reports whether the given instruction or terminator has a void
	// result type. Only call, invoke and callbr instructions are considered.
	// This is important as it determines the assignment of IDs to unnamed local
	// variables.
	isVoid map[local]bool
}

// newFuncGen returns a new generator for the given IR function.
func newFuncGen(gen *generator, f *ir.Func) *funcGen {
	return &funcGen{
		gen:    gen,
		f:      f,
		locals: make(map[ir.LocalIdent]value.Value),
		isVoid: make(map[local]bool),
	}
}

// resolveLocals resolves the local variables (function parameters, basic
// blocks, results of instructions and terminators) of the given function body.
func (fgen *funcGen) resolveLocals(old ast.FuncBody) error {
	// Index local identifiers and create scaffolding IR local variables (without
	// bodies but with types).
	oldBlocks := old.Blocks()
	if err := fgen.createLocals(oldBlocks); err != nil {
		return errors.WithStack(err)
	}
	// Translate AST instructions to IR.
	if err := fgen.translateInsts(oldBlocks); err != nil {
		return errors.WithStack(err)
	}
	// Translate AST terminators to IR.
	return fgen.translateTerms(oldBlocks)
}

// === [ Create and index IR ] =================================================

// local is a local variable.
type local interface {
	value.Named
	// ID returns the ID of the local identifier.
	ID() int64
	// SetID sets the ID of the local identifier.
	SetID(id int64)
	// IsUnnamed reports whether the local identifier is unnamed.
	IsUnnamed() bool
}

// createLocals indexes local identifiers and creates scaffolding IR local
// variables (without bodies but with types) of the given function.
//
// post-condition: fgen.locals maps from local identifier (without '%' prefix)
// to corresponding skeleton IR value.
func (fgen *funcGen) createLocals(oldBlocks []ast.BasicBlock) error {
	// Create local variable skeletons (without bodies but with types).
	if err := fgen.newLocals(oldBlocks); err != nil {
		return errors.WithStack(err)
	}
	// Assign local IDs.
	//
	// Note: the type of call instructions and invoke terminators must be
	// determined before assigning local IDs, as they may be values or non-values
	// based on return type. This is done by fgen.newLocals.
	if err := fgen.assignIDs(fgen.f); err != nil {
		return errors.WithStack(err)
	}
	// Index local identifiers.
	return fgen.indexLocals()
}

// newLocals creates scaffolding IR local variables (without bodies but with
// types) of the given function.
func (fgen *funcGen) newLocals(oldBlocks []ast.BasicBlock) error {
	// Note: function parameters are already translated in gen.irFuncHeader.
	f := fgen.f
	f.Blocks = make([]*ir.Block, len(oldBlocks))
	for i, oldBlock := range oldBlocks {
		block := &ir.Block{}
		if n, ok := oldBlock.Name(); ok {
			block.LocalIdent = labelIdent(n)
		}
		if oldInsts := oldBlock.Insts(); len(oldInsts) > 0 {
			block.Insts = make([]ir.Instruction, len(oldInsts))
			for j, oldInst := range oldInsts {
				inst, err := fgen.newInst(oldInst)
				if err != nil {
					return errors.WithStack(err)
				}
				block.Insts[j] = inst
			}
		}
		term, err := fgen.newTerm(oldBlock.Term())
		if err != nil {
			return errors.WithStack(err)
		}
		block.Term = term
		block.Parent = f
		f.Blocks[i] = block
	}
	return nil
}

// indexLocals indexes local identifiers of the given function.
func (fgen *funcGen) indexLocals() error {
	// Index function parameters.
	f := fgen.f
	for _, param := range f.Params {
		if err := fgen.addLocal(param.LocalIdent, param); err != nil {
			return errors.WithStack(err)
		}
	}
	// Index basic blocks.
	for _, block := range f.Blocks {
		if err := fgen.addLocal(block.LocalIdent, block); err != nil {
			return errors.WithStack(err)
		}
		// Index instructions.
		for _, inst := range block.Insts {
			v, ok := inst.(local)
			if !ok || fgen.isVoid[v] {
				// Skip non-value instructions.
				continue
			}
			ident := localIdentOfValue(v)
			if err := fgen.addLocal(ident, v); err != nil {
				return errors.WithStack(err)
			}
		}
		// Index terminator.
		v, ok := block.Term.(local)
		if !ok || fgen.isVoid[v] {
			// Skip non-value terminators.
			continue
		}
		ident := localIdentOfValue(v)
		if err := fgen.addLocal(ident, v); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// ### [ Helper functions ] ####################################################

// addLocal adds the local variable with the given local identifier to the map
// of local variables of the function.
func (fgen *funcGen) addLocal(ident ir.LocalIdent, v value.Value) error {
	if prev, ok := fgen.locals[ident]; ok {
		return errors.Errorf("local identifier %q already present; prev `%s`, new `%s`", ident.Ident(), prev, v)
	}
	fgen.locals[ident] = v
	return nil
}

// localIdentOfValue returns the local identifier of the given local variable.
func localIdentOfValue(v local) ir.LocalIdent {
	if v.IsUnnamed() {
		return ir.LocalIdent{LocalID: v.ID()}
	}
	return ir.LocalIdent{LocalName: v.Name()}
}

// NOTE: assignIDs is copied from ir/func.go with slight modifications to handle
// type lookup of callees using fgen.isVoid instead of callee.Type(), as the
// callees are not yet resolved when we need to assign IDs.

// assignIDs assigns IDs to unnamed local variables.
func (fgen *funcGen) assignIDs(f *ir.Func) error {
	if len(f.Blocks) == 0 {
		return nil
	}
	id := int64(0)
	setName := func(n local) error {
		if n.IsUnnamed() {
			if n.ID() != 0 && id != n.ID() {
				want := id
				got := n.ID()
				return errors.Errorf("invalid local ID in function %q, expected %s, got %s", f.Ident(), enc.LocalID(want), enc.LocalID(got))
			}
			n.SetID(id)
			id++
		}
		return nil
	}
	for _, param := range f.Params {
		// Assign local IDs to unnamed parameters of function definitions.
		if err := setName(param); err != nil {
			return errors.WithStack(err)
		}
	}
	for _, block := range f.Blocks {
		// Assign local IDs to unnamed basic blocks.
		if err := setName(block); err != nil {
			return errors.WithStack(err)
		}
		for _, inst := range block.Insts {
			n, ok := inst.(local)
			if !ok {
				continue
			}
			// Skip void instructions.
			// TODO: Check if any other value instructions than call may have void
			// type.
			if fgen.isVoidValue(n) {
				continue
			}
			// Assign local IDs to unnamed local variables.
			if err := setName(n); err != nil {
				return errors.WithStack(err)
			}
		}
		n, ok := block.Term.(local)
		if !ok {
			continue
		}
		if fgen.isVoidValue(n) {
			continue
		}
		if err := setName(n); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// isVoidValue reports whether the given named value is a non-value (i.e. a call
// instruction, or an invoke or callbr terminator with void-return type).
func (fgen *funcGen) isVoidValue(n local) bool {
	switch n.(type) {
	case *ir.InstCall, *ir.TermInvoke, *ir.TermCallBr:
		return fgen.isVoid[n]
	}
	return false
}
