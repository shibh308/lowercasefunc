package lowercasefunc

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"path/filepath"
	"strings"
)

const doc = "lowercasefunc is ..."

var Analyzer = &analysis.Analyzer{
	Name: "lowercasefunc",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

var DetectAnalyzer = &analysis.Analyzer{
	Name: "lowercasefuncdetector",
	Doc:  doc,
	Run:  runDetect,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func getFuncObj(pkg *types.Package, name string) types.Object {
	if obj := pkg.Scope().Lookup(name); obj == nil {
		return nil
	} else if _, ok := obj.(*types.Func); ok {
		return obj
	}
	return nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	var result []FuncPair

	exportedFunc := make(map[string]*ast.FuncDecl)
	unexportedFunc := make(map[string]*ast.FuncDecl)

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			name := fd.Name.String()
			if fd.Name.IsExported() {
				exportedFunc[name] = fd
			} else {
				unexportedFunc[name] = fd
			}
		}
	}
	for name, fd := range exportedFunc {
		lower := strings.ToLower(string(name[0])) + name[1:]
		ufd, ok := unexportedFunc[lower]
		if !ok {
			continue
		}
		target := pass.TypesInfo.Defs[ufd.Name]
		body := fd.Body
		var calledPos []token.Pos
		ast.Inspect(body, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if ok && pass.TypesInfo.Uses[id] == target {
				calledPos = append(calledPos, id.Pos())
			}
			return true
		})
		p := pass.Fset.Position(target.Pos())
		s := fmt.Sprintf("{FuncName:%s, TargetPos:\"%s:%d:%d\", CalledPos:[", fd.Name.String(), filepath.Base(p.Filename), p.Column, p.Line)
		var cpStr []string
		for _, cp := range calledPos {
			p := pass.Fset.Position(cp)
			cpStr = append(cpStr, fmt.Sprintf("\"%s:%d:%d\"", filepath.Base(p.Filename), p.Column, p.Line))
		}
		pass.Reportf(fd.Pos(), s + strings.Join(cpStr, ", ") + "]}\n")
		result = append(result,
			FuncPair{
				UpperDecl: fd,
				LowerDecl: ufd,
				CalledPos: calledPos,
			})
	}
	return nil, nil
}

type FuncPair struct{
	UpperDecl    *ast.FuncDecl
	LowerDecl    *ast.FuncDecl
	CalledPos     []token.Pos
}

func runDetect(pass *analysis.Pass) (interface{}, error) {
	var result []FuncPair

	exportedFunc := make(map[string]*ast.FuncDecl)
	unexportedFunc := make(map[string]*ast.FuncDecl)

	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			name := fd.Name.String()
			if fd.Name.IsExported() {
				exportedFunc[name] = fd
			} else {
				unexportedFunc[name] = fd
			}
		}
	}
	for name, fd := range exportedFunc {
		lower := strings.ToLower(string(name[0])) + name[1:]
		ufd, ok := unexportedFunc[lower]
		if !ok {
			continue
		}
		target := pass.TypesInfo.Defs[ufd.Name]
		body := fd.Body
		var calledPos []token.Pos
		ast.Inspect(body, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if ok && pass.TypesInfo.Uses[id] == target {
				calledPos = append(calledPos, id.Pos())
			}
			return true
		})
		result = append(result,
			FuncPair{
				UpperDecl: fd,
				LowerDecl: ufd,
				CalledPos: calledPos,
			})
	}

	return result, nil
}

