package docgen

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// internal
// ========
type file struct {
	fset     *token.FileSet
	astFile  *ast.File
	src      []byte
	filename string

	main bool

	decls map[string]string
	boxes []Box
}

func (f *file) walk(fn func(ast.Node) bool) {
	ast.Walk(walker(fn), f.astFile)
}

func (f *file) find() []Box {
	f.findBoxCalls()
	return f.boxes
}

func (f *file) findBoxCalls() {
	// TODO: figure out how to extract type function call from node graph
	f.walk(func(node ast.Node) bool {
		ce, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		isLoad := isPkgDot(ce.Fun, "sea", "Load")
		isLoadWithDefault := isPkgDot(ce.Fun, "sea", "LoadWithDefault")
		if !(isLoadWithDefault || isLoad) || len(ce.Args) < 1 {
			return true
		}

		box := Box{}
		if isLoad {
			box.args = append(box.args, `"sea.Load"`)
		}
		if isLoadWithDefault {
			box.args = append(box.args, `"sea.LoadWithDefault"`)
		}

		for _, arg := range ce.Args {
			switch x := arg.(type) {
			case *ast.BasicLit:
				box.args = append(box.args, x.Value)
			}
		}
		f.boxes = append(f.boxes, box)
		return true
	})
}

// helpers
// =======
func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	return ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}

// wrap a function to fulfill ast.Visitor interface
type walker func(ast.Node) bool

func (w walker) Visit(node ast.Node) ast.Visitor {
	if w(node) {
		return w
	}
	return nil
}

// exports
// =======
type Box struct {
	args []string
}

func FindCalls(filename string) ([]Box, error) {

	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	f := &file{fset: fset, astFile: astFile, src: nil, filename: filename}
	f.decls = make(map[string]string)

	return f.find(), nil

}
