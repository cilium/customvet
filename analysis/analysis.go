// Copyright 2020 Authors of Cilium
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package analysis

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `This is a custom vet check for the cilium codebase.`

var Analyzer = &analysis.Analyzer{
	Name:     "cilium",
	Doc:      Doc,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

type checker interface {
	check(n ast.Node) (bool, string, token.Pos)
}

type visitor func(ast.Node) bool

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if v(node) {
		return v
	}
	return nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspct, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("analyzer is not type *inspector.Inspector")
	}

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.File)(nil),
		(*ast.ImportSpec)(nil),
	}
	checks := []checker{
		&timeAfterCheck{},
	}
	inspct.Preorder(nodeFilter, func(n ast.Node) {
		for _, c := range checks {
			ok, report, pos := c.check(n)
			if !ok {
				pass.Reportf(pos, report)
			}
		}
	})
	return nil, nil
}

func isPkgDot(expr ast.Expr, pkg, name string) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	res := ok && isIdent(sel.X, pkg) && isIdent(sel.Sel, name)
	return res
}

func isIdent(expr ast.Expr, ident string) bool {
	id, ok := expr.(*ast.Ident)
	return ok && id.Name == ident
}
