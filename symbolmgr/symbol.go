package symbolmgr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/linuxerwang/goats-html/pkgmgr"
)

const (
	TypeArg = iota
	TypeImport
	TypeVar
	TypeFor
)

// Symbol represents a symbol defined and used in a tag scope.
type Symbol struct {
	Name    string
	Type    int
	IsPb    bool
	PkgImpt *pkgmgr.PkgImport
}

func (s *Symbol) String() string {
	return fmt.Sprintf("name: %s, type: %d, is_pb: %t, PkgImpt: %+v", s.Name, s.Type, s.IsPb, s.PkgImpt)
}

// SymbolMgr manages all symbols of a template.
type SymbolMgr struct {
	symStack []map[string]*Symbol
}

// Push pushes a symbol map into its inner stack.
func (sm *SymbolMgr) Push(symbols map[string]*Symbol) {
	sm.symStack = append(sm.symStack, symbols)
}

// Pop pops and drop the top symbol map in the stack.
func (sm *SymbolMgr) Pop() {
	sm.symStack = sm.symStack[:len(sm.symStack)-1]
}

// Peek returns the top symbol map in the stack and the map is updatable.
func (sm *SymbolMgr) Peek() map[string]*Symbol {
	return sm.symStack[len(sm.symStack)-1]
}

// Size returns the size of the symbol stack.
func (sm *SymbolMgr) Size() int {
	return len(sm.symStack)
}

// GetSymbol returns a symbol for the given name. It searches the name in
// each map of stack from top to bottom, so that if a name was defined twice,
// the higher one in stack shadows the lower one.
func (sm *SymbolMgr) GetSymbol(name string) (*Symbol, error) {
	for i := len(sm.symStack) - 1; i >= 0; i-- {
		symbols := sm.symStack[i]
		if s, ok := symbols[name]; ok {
			return s, nil
		}
	}
	return nil, errors.New("unknown symbol: " + name)
}

// ExpandType expands the given PB type with full qualified name, or the
// input unchanged if it's not a PB type.
func (sm *SymbolMgr) ExpandType(varType string) string {
	n := strings.Index(varType, ".")
	if n < 1 {
		return varType
	}

	if sym, ok := sm.symStack[0][varType[:n]]; ok {
		if sym.Type != TypeImport {
			return varType
		}

		pbPkg := sym.PkgImpt.PbPkg()
		if pbPkg == "" {
			return varType
		}

		return pbPkg + varType[n:]
	}

	return varType
}

func (sm *SymbolMgr) AllSymbols() map[string]*Symbol {
	all := make(map[string]*Symbol)
	for _, m := range sm.symStack {
		for k, v := range m {
			all[k] = v
		}
	}
	return all
}

func New() *SymbolMgr {
	mgr := &SymbolMgr{}
	return mgr
}
