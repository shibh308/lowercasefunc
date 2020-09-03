package lowercasefunc

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || !fd.Name.IsExported() {
				continue
			}
			name := fd.Name.String()
			lower := strings.ToLower(string(name[0])) + name[1:]
			target := getFuncObj(pass.Pkg, lower)
			if target == nil {
				continue
			}
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
		}
	}

	return nil, nil
}

type LowerCaseFunc struct{
	UpperDecl    *ast.FuncDecl
	LowerObj     types.Object
	CalledPos     []token.Pos
}

func runDetect(pass *analysis.Pass) (interface{}, error) {
	var result []LowerCaseFunc
	for _, f := range pass.Files {
		for _, decl := range f.Decls {
			fd, ok := decl.(*ast.FuncDecl)
			if !ok || !fd.Name.IsExported() {
				continue
			}
			name := fd.Name.String()
			lower := strings.ToLower(string(name[0])) + name[1:]
			target := getFuncObj(pass.Pkg, lower)
			if target == nil {
				continue
			}
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
				LowerCaseFunc{
					UpperDecl: fd,
					LowerObj: target,
					CalledPos: calledPos,
				})
		}
	}

	return result, nil
}

