package goats

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"strings"
)

type GenType int

const (
	GenInterfaceImports GenType = iota
	GenImplImports
)

type pkgImport struct {
	name  string
	alias string
	path  string
}

func (pi *pkgImport) generateImports(buffer *bytes.Buffer) {
	if pi.alias != "" {
		buffer.WriteString(fmt.Sprintf("%s \"%s\"\n", pi.alias, pi.path))
	} else {
		buffer.WriteString(fmt.Sprintf("\"%s\"\n", pi.path))
	}
}

func (pi *pkgImport) Alias() string {
	return pi.alias
}

// PkgRefs maintains package imports for a specific template.
type PkgRefs struct {
	pkgMgr *PkgManager
	pkgs   map[string]bool // map from package path to whether it's for interface.
}

func (pr *PkgRefs) ParseTmplCall(callStmt string) (pkgPath, callName string) {
	return pr.pkgMgr.ParseTmplCall(callStmt)
}

func (pr *PkgRefs) RefByPath(pkgPath string, forInterface bool) *pkgImport {
	pi := pr.pkgMgr.PkgByPath(pkgPath)
	if pi != nil {
		pr.pkgs[pi.path] = forInterface
	}
	return pi
}

func (pr *PkgRefs) RefByAlias(alias string) *pkgImport {
	pi := pr.pkgMgr.PkgByAlias(alias)
	if pi != nil {
		pr.pkgs[pi.path] = true
	}
	return pi
}

func (pr *PkgRefs) GenerateImports(buffer *bytes.Buffer, genType GenType) {
	for path, forIface := range pr.pkgs {
		if genType == GenInterfaceImports && !forIface {
			continue
		}
		if genType == GenImplImports && forIface {
			continue
		}
		pr.pkgMgr.PkgByPath(path).generateImports(buffer)
	}
}

// PackageManager maintains all package imports in a template file.
type PkgManager struct {
	aliasId      int
	tmplPkg      string
	pkgsForPath  map[string]*pkgImport // maps from package path to pkgImport
	pkgsForAlias map[string]*pkgImport // maps from package alias to pkgImport
}

func (pm *PkgManager) AddImport(alias, pkgPath string) *pkgImport {
	if pkgPath == "" {
		return nil
	}

	alias = TrimWhiteSpaces(alias)

	if pi, found := pm.pkgsForAlias[alias]; found {
		return pi
	}

	pkgPath = TrimWhiteSpaces(pkgPath)
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

	pi := &pkgImport{
		name:  pkgName,
		alias: alias,
		path:  pkgPath,
	}

	pm.pkgsForPath[pi.path] = pi
	pm.pkgsForAlias[pi.alias] = pi

	return pi
}

func (pm *PkgManager) ParseTmplCall(callStmt string) (pkgPath, callName string) {
	callStmt = TrimWhiteSpaces(callStmt)

	if !strings.Contains(callStmt, "#") {
		log.Fatal(`Call to template must contain a "#".`)
	}

	parts := strings.SplitN(callStmt, "#", 2)

	if strings.HasPrefix(callStmt, "#") {
		// In-file call.
		return "", parts[1]
	}

	pkgPath = strings.Replace(parts[0], ".html", "_html", -1)

	if strings.HasPrefix(pkgPath, "/") {
		// Absolute reference
		return pkgPath[1:], parts[1]
	} else {
		return path.Join(pm.tmplPkg, pkgPath), parts[1]
	}
}

func (pm *PkgManager) PkgByPath(pkgPath string) *pkgImport {
	pkgPath = path.Clean(TrimWhiteSpaces(pkgPath))
	if pi, found := pm.pkgsForPath[pkgPath]; found {
		return pi
	}
	return nil
}

func (pm *PkgManager) PkgByAlias(alias string) *pkgImport {
	if pi, found := pm.pkgsForAlias[alias]; found {
		return pi
	}
	return nil
}

func (pm *PkgManager) CreatePkgRefs() *PkgRefs {
	return &PkgRefs{
		pkgMgr: pm,
		pkgs:   map[string]bool{},
	}
}

func NewPkgManager(tmplPkg string) *PkgManager {
	return &PkgManager{
		aliasId:      0,
		tmplPkg:      tmplPkg,
		pkgsForPath:  map[string]*pkgImport{},
		pkgsForAlias: map[string]*pkgImport{},
	}
}
