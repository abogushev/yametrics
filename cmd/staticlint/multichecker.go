// Кастомный линтер
//
// Включает в себя:
//
// - SA* staticcheck: категория проверок SA под кодовым названием staticcheck включает все проверки, связанные с правильностью кода. https://staticcheck.io/docs/checks/#SA
//
// - simple: проверки категории S под кодовым названием «простые» содержат все проверки, связанные с упрощением кода. https://staticcheck.io/docs/checks/#S
//
// - stylecheck: категория проверок ST под кодовым названием stylecheck содержит все проверки, связанные со стилистическими вопросами. https://staticcheck.io/docs/checks/#ST
//
// - printf: анализатор, который проверяет согласованность строк и аргументов формата Printf.
//
// - shadow: анализатор, который проверяет теневые переменные.
//
// - structtag: анализатор, который проверяет правильность формирования тегов полей структуры.
//
// - OSExitCheckAnalyzer: анализатор проверки вызова os.Exist в функции main пакета main
package main

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

func main() {
	var checks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if len(v.Analyzer.Name) >= 2 && v.Analyzer.Name[0:2] == "SA" {
			checks = append(checks, v.Analyzer)
		}
	}
	for _, v := range simple.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	for _, v := range stylecheck.Analyzers {
		checks = append(checks, v.Analyzer)
	}

	checks = append(checks, printf.Analyzer)
	checks = append(checks, shadow.Analyzer)
	checks = append(checks, structtag.Analyzer)
	checks = append(checks, OSExitCheckAnalyzer)

	multichecker.Main(checks...)
}

var OSExitCheckAnalyzer = &analysis.Analyzer{
	Name: "osExitCheck",
	Doc:  "check for os exit",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			if file.Name.Name != "main" {
				return true
			}
			fun, ok := node.(*ast.FuncDecl)
			if !ok {
				return true
			}
			if fun.Name.Name != "main" {
				return true
			}
			for _, stmt := range fun.Body.List {
				exprStmt, ok := stmt.(*ast.ExprStmt)
				if !ok {
					continue
				}
				f, ok := exprStmt.X.(*ast.CallExpr)
				if !ok {
					continue
				}
				ff, ok := f.Fun.(*ast.SelectorExpr)
				if !(ok && ff.Sel.Name == "Exit") {
					continue
				}
				ident, ok := ff.X.(*ast.Ident)
				if !(ok && ident.Name == "os") {
					continue
				}
				pass.Reportf(file.Pos(), "calling os.Exist is unwanted")
			}
			return true
		})
	}
	return nil, nil
}
