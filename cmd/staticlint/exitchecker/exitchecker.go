// Package implementing the exit checker
package exitchecker

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const Doc = `check for os.Exit in main() functions.

The exit checker looks for :

	os.Exit()

calls in main() functions.`

var ExitCheckAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  Doc,
	Run:  checkExit,
}

// checkExit walks the os.Exit method calls in main() functions.
func checkExit(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if x, ok := node.(*ast.FuncDecl); ok {
				if x.Name.String() == "main" {
					for _, stmt := range x.Body.List {
						if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
							if call, ok := exprStmt.X.(*ast.CallExpr); ok {
								if fun, ok := call.Fun.(*ast.SelectorExpr); ok {
									pkg, ok := fun.X.(*ast.Ident)
									if ok {
										if (pkg.Name + "." + fun.Sel.Name) == "os.Exit" {
											pass.ReportRangef(exprStmt, "has os.Exit function")
										}
									}
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return nil, nil
}
