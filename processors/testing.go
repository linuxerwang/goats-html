package processors

import (
	"io"
)

type dummyAliasGetter struct {
	alias string
}

func (dag dummyAliasGetter) Alias() string {
	return dag.alias
}

type dummyAliasReferer struct {
	aliases map[string]bool
	paths   map[string]bool
}

func (dar *dummyAliasReferer) RefByAlias(alias string, forInterface bool) AliasGetter {
	dar.aliases[alias] = forInterface
	return dummyAliasGetter{alias: alias}
}

func (dar *dummyAliasReferer) RefByPath(pkgPath string, forInterface bool) AliasGetter {
	dar.paths[pkgPath] = forInterface
	return dummyAliasGetter{alias: pkgPath}
}

func NewDummyAliasReferer() *dummyAliasReferer {
	return &dummyAliasReferer{
		aliases: map[string]bool{},
		paths:   map[string]bool{},
	}
}

type dummyProcessor struct {
	BaseProcessor
	Called bool
}

func (dp *dummyProcessor) Process(writer io.Writer, context *TagContext) {
	io.WriteString(writer, "DUMMY")
	dp.Called = true
}

func NewDummyProcessor() *dummyProcessor {
	return &dummyProcessor{}
}
