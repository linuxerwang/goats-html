package pkgmgr

import (
	"fmt"
	"io"
	"log"
	"path"
	"strings"

	"github.com/linuxerwang/goats-html/util"
)

type GenType int

const (
	GenInterfaceImports GenType = iota
	GenImplImports
)

type AliasGetter interface {
	Alias() string
}

type AliasReferer interface {
	RefByAlias(string, bool) AliasGetter
	RefByPath(string, bool) AliasGetter
	RefClosureRequire(string)
}

type PkgImport struct {
	name  string
	alias string
	path  string
	pbPkg string
}

func (pi *PkgImport) generateImports(output io.Writer) {
	if pi.alias != "" {
		io.WriteString(output, fmt.Sprintf("%s \"%s\"\n", pi.alias, pi.path))
	} else {
		io.WriteString(output, fmt.Sprintf("\"%s\"\n", pi.path))
	}
}

func (pi *PkgImport) Name() string {
	return pi.name
}

func (pi *PkgImport) SetName(name string) {
	pi.name = name
}

func (pi *PkgImport) Alias() string {
	return pi.alias
}

func (pi *PkgImport) SetAlias(alias string) {
	pi.alias = alias
}

func (pi *PkgImport) Path() string {
	return pi.path
}

func (pi *PkgImport) SetPath(path string) {
	pi.path = path
}

func (pi *PkgImport) PbPkg() string {
	return pi.pbPkg
}

func (pi *PkgImport) SetPbPkg(pbPkg string) {
	pi.pbPkg = pbPkg
}

// PkgRefs maintains package imports for a specific template.
type PkgRefs struct {
	pkgMgr          *PkgManager
	pkgs            map[string]bool // map from package path to whether it's for interface.
	closureRequires map[string]bool // set of closure require targets.
}

func (pr *PkgRefs) ParseTmplCall(callStmt string) (pkgPath, relPkgPath, callName string) {
	return pr.pkgMgr.ParseTmplCall(callStmt)
}

func (pr *PkgRefs) RefByPath(pkgPath string, forInterface bool) AliasGetter {
	pi := pr.pkgMgr.PkgByPath(pkgPath)
	if pi != nil {
		pr.pkgs[pi.path] = forInterface
	}
	return pi
}

func (pr *PkgRefs) RefByAlias(alias string, forInterface bool) AliasGetter {
	pi := pr.pkgMgr.PkgByAlias(alias)
	if pi != nil {
		pr.pkgs[pi.path] = forInterface
	}
	return pi
}

func (pr *PkgRefs) RefClosureRequire(require string) {
	pr.closureRequires[require] = true
}

func (pr *PkgRefs) GenerateImports(output io.Writer, genType GenType) {
	for path, forIface := range pr.pkgs {
		if genType == GenInterfaceImports && !forIface {
			continue
		}
		if genType == GenImplImports && forIface {
			continue
		}
		pr.pkgMgr.PkgByPath(path).generateImports(output)
	}
}

func (pr *PkgRefs) GenerateRequires(output io.Writer) {
	for r, _ := range pr.closureRequires {
		io.WriteString(output, fmt.Sprintf("goog.require('%s');\n", r))
	}
}

func (pr *PkgRefs) DumpRequires(requires map[string]bool) {
	for r, _ := range pr.closureRequires {
		requires[r] = true
	}
}

// PackageManager maintains all package imports in a template file.
type PkgManager struct {
	aliasId         int
	tmplPkg         string
	outputPkgPrefix string
	pkgsForPath     map[string]*PkgImport // maps from package path to PkgImport
	pkgsForAlias    map[string]*PkgImport // maps from package alias to PkgImport
}

func (pm *PkgManager) AddImport(alias, pkgPath string, pbPkg string) *PkgImport {
	if pkgPath == "" {
		return nil
	}

	alias = util.TrimWhiteSpaces(alias)

	if pi, found := pm.pkgsForAlias[alias]; found {
		return pi
	}

	pkgPath = util.TrimWhiteSpaces(pkgPath)
	pkgName := path.Base(pkgPath)

	if alias == "" {
		// Generate alias.
		alias = fmt.Sprintf("__%s_%d", pkgName, pm.aliasId)
		pm.aliasId++
	}

	pkgPath = path.Clean(pkgPath)

	if pi, found := pm.pkgsForPath[pkgPath]; found {
		return pi
	}

	pi := &PkgImport{
		name:  pkgName,
		alias: alias,
		path:  pkgPath,
		pbPkg: pbPkg,
	}

	pm.pkgsForPath[pi.path] = pi
	pm.pkgsForAlias[pi.alias] = pi

	return pi
}

func (pm *PkgManager) ParseTmplCall(callStmt string) (pkgPath, relPkgPath, callName string) {
	callStmt = util.TrimWhiteSpaces(callStmt)

	if !strings.Contains(callStmt, "#") {
		log.Fatal(`Call to template must contain a "#".`)
	}

	parts := strings.SplitN(callStmt, "#", 2)

	if strings.HasPrefix(callStmt, "#") {
		// In-file call.
		return "", "", parts[1]
	}

	pkgPath = strings.Replace(parts[0], ".html", "_html", -1)

	if strings.HasPrefix(pkgPath, "/") {
		// Absolute reference
		return path.Join(pm.outputPkgPrefix, pkgPath[1:]), pkgPath[1+len(pm.tmplPkg):], parts[1]
	} else {
		return path.Join(pm.outputPkgPrefix, pm.tmplPkg, pkgPath), pkgPath, parts[1]
	}
}

func (pm *PkgManager) PkgByPath(pkgPath string) *PkgImport {
	pkgPath = path.Clean(util.TrimWhiteSpaces(pkgPath))
	if pi, found := pm.pkgsForPath[pkgPath]; found {
		return pi
	}
	return nil
}

func (pm *PkgManager) PkgByAlias(alias string) *PkgImport {
	if pi, found := pm.pkgsForAlias[alias]; found {
		return pi
	}
	return nil
}

func (pm *PkgManager) GetPkgsForAlias() map[string]*PkgImport {
	return pm.pkgsForAlias
}

func (pm *PkgManager) CreatePkgRefs() *PkgRefs {
	return &PkgRefs{
		pkgMgr:          pm,
		pkgs:            map[string]bool{},
		closureRequires: map[string]bool{},
	}
}

func New(tmplPkg, outputPkgPrefix string) *PkgManager {
	return &PkgManager{
		aliasId:         0,
		tmplPkg:         tmplPkg,
		outputPkgPrefix: outputPkgPrefix,
		pkgsForPath:     map[string]*PkgImport{},
		pkgsForAlias:    map[string]*PkgImport{},
	}
}
