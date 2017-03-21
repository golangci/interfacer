// Copyright (c) 2015, Daniel Martí <mvdan@mvdan.cc>
// See LICENSE for licensing information

package interfacer

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/loader"
)

//go:generate sh -c "go list std | go run generate/std/main.go -o std.go"
//go:generate gofmt -w -s std.go

type cache struct {
	loader.Config

	cur pkgCache

	grabbed map[string]pkgCache
}

type pkgCache struct {
	exp, unexp typeSet
}

type typeSet struct {
	ifaces map[string]string
	funcs  map[string]string
}

func (c *cache) isFuncType(t string) bool {
	if stdFuncs[t] {
		return true
	}
	if s := c.cur.exp.funcs[t]; s != "" {
		return true
	}
	return c.cur.unexp.funcs[t] != ""
}

func (c *cache) ifaceOf(t string) string {
	if s := stdIfaces[t]; s != "" {
		return s
	}
	if s := c.cur.exp.ifaces[t]; s != "" {
		return s
	}
	return c.cur.unexp.ifaces[t]
}

func (c *cache) grabNames(pkg *types.Package) {
	c.fillCache(pkg)
	c.cur = c.grabbed[pkg.Path()]
}

func (c *cache) fillCache(pkg *types.Package) {
	path := pkg.Path()
	if _, e := c.grabbed[path]; e {
		return
	}
	for _, imp := range pkg.Imports() {
		c.fillCache(imp)
	}
	cur := pkgCache{
		exp: typeSet{
			ifaces: make(map[string]string),
			funcs:  make(map[string]string),
		},
		unexp: typeSet{
			ifaces: make(map[string]string),
			funcs:  make(map[string]string),
		},
	}
	addTypes := func(impPath string, ifs, funs map[string]string, top bool) {
		fullName := func(name string) string {
			if !top {
				return impPath + "." + name
			}
			return name
		}
		for iftype, name := range ifs {
			if _, e := stdIfaces[iftype]; e {
				continue
			}
			if ast.IsExported(name) {
				cur.exp.ifaces[iftype] = fullName(name)
			}
		}
		for ftype, name := range funs {
			if stdFuncs[ftype] {
				continue
			}
			if ast.IsExported(name) {
				cur.exp.funcs[ftype] = fullName(name)
			} else {
				cur.unexp.funcs[ftype] = fullName(name)
			}
		}
	}
	for _, imp := range pkg.Imports() {
		pc := c.grabbed[imp.Path()]
		addTypes(imp.Path(), pc.exp.ifaces, pc.exp.funcs, false)
	}
	ifs, funs := FromScope(pkg.Scope())
	addTypes(path, ifs, funs, true)
	c.grabbed[path] = cur
}
