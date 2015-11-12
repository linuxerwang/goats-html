package pkgmgr

type dummyAliasGetter struct {
	alias string
}

func (dag dummyAliasGetter) Alias() string {
	return dag.alias
}

type dummyAliasReferer struct {
	aliases  map[string]bool
	paths    map[string]bool
	requires map[string]bool
}

func (dar *dummyAliasReferer) Aliases() map[string]bool {
	return dar.aliases
}

func (dar *dummyAliasReferer) RefByAlias(alias string, forInterface bool) AliasGetter {
	dar.aliases[alias] = forInterface
	return dummyAliasGetter{alias: alias}
}

func (dar *dummyAliasReferer) RefByPath(pkgPath string, forInterface bool) AliasGetter {
	dar.paths[pkgPath] = forInterface
	return dummyAliasGetter{alias: pkgPath}
}

func (dar *dummyAliasReferer) RefClosureRequire(str string) {
	dar.requires[str] = true
}

func NewDummyAliasReferer() *dummyAliasReferer {
	return &dummyAliasReferer{
		aliases:  map[string]bool{},
		paths:    map[string]bool{},
		requires: map[string]bool{},
	}
}
