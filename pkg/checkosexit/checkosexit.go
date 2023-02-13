// Package checkosexit implements checking the call to
// os.Exit in the main function of the main package.
package checkosexit

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

// Analyzer is the check for calls to os.Exit.
var Analyzer = &analysis.Analyzer{
	Name: "checkosexit",
	Doc:  "checking the use of a direct call to os.Exit in the main function of the main package",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.File:
				// if the package is not main, then we do not check further
				if x.Name.Name != "main" {
					return false
				}

				_, ok := x.Scope.Objects["tests"]
				if ok {
					return false
				}
			case *ast.FuncDecl:
				// if the function is not main, then we do not check further
				if x.Name.Name != "main" {
					return false
				}
			case *ast.SelectorExpr:
				// check direct call to os.Exit
				if fmt.Sprintf("%v", x.X) == "os" && x.Sel.Name == "Exit" {
					pass.Report(analysis.Diagnostic{
						Pos:     x.Pos(),
						Message: "direct call to os.Exit found",
					})
				}
			}

			return true
		})
	}

	return nil, nil
}
