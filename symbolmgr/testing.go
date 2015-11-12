package symbolmgr

import (
	"github.com/linuxerwang/goats-html/pkgmgr"
)

type DummySymbolMgrBuilder struct {
	smgr *SymbolMgr
}

func (b *DummySymbolMgrBuilder) Build() *SymbolMgr {
	return b.smgr
}

func (b *DummySymbolMgrBuilder) AddArg(name string, isPb bool) {
	symbol := &Symbol{
		Name: name,
		Type: TypeArg,
		IsPb: true,
	}
	b.smgr.Peek()[name] = symbol
}

func (b *DummySymbolMgrBuilder) AddImport(name string, isPb bool, path, pbPkg string) {
	symbol := &Symbol{
		Name:    name,
		Type:    TypeImport,
		IsPb:    true,
		PkgImpt: &pkgmgr.PkgImport{},
	}
	symbol.PkgImpt.SetName(name)
	symbol.PkgImpt.SetAlias(name)
	symbol.PkgImpt.SetPath(path)
	symbol.PkgImpt.SetPbPkg(pbPkg)
	b.smgr.Peek()[name] = symbol
}

func NewSymbolMgrBuilder() *DummySymbolMgrBuilder {
	smgr := New()
	sm := make(map[string]*Symbol)
	smgr.Push(sm)

	return &DummySymbolMgrBuilder{
		smgr: smgr,
	}
}
