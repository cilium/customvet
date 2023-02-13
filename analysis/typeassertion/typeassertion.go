// Copyright 2023 Authors of Cilium
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

package typeassertion

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// Analyzer is the global for the multichecker
var Analyzer = &analysis.Analyzer{
	Name:     "typeassertion",
	Doc:      "This checks for unchecked type assertions.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

type visitor func(ast.Node) bool

func (v visitor) Visit(node ast.Node) ast.Visitor {
	if v(node) {
		return v
	}
	return nil
}

type typeAssertionLoc struct {
	LParen token.Pos
	RParen token.Pos
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspct, ok := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	if !ok {
		return nil, fmt.Errorf("analyzer is not type *inspector.Inspector")
	}

	var (
		nodeFilter = []ast.Node{
			(*ast.TypeAssertExpr)(nil),
			(*ast.AssignStmt)(nil),
		}
		validTypeAssertions = make(map[typeAssertionLoc]struct{})
	)
	inspct.Preorder(nodeFilter, func(n ast.Node) {
		switch stmt := n.(type) {
		case *ast.AssignStmt:
			if len(stmt.Lhs) == 2 && len(stmt.Rhs) == 1 {
				ta, ok := stmt.Rhs[0].(*ast.TypeAssertExpr)
				if ok {
					validTypeAssertions[typeAssertionLoc{
						LParen: ta.Lparen,
						RParen: ta.Rparen,
					}] = struct{}{}
				}
			}
		case *ast.TypeAssertExpr:
			_, ok := validTypeAssertions[typeAssertionLoc{
				LParen: stmt.Lparen,
				RParen: stmt.Rparen,
			}]
			if !ok {
				pass.Reportf(n.Pos(), "use of an unchecked type assertion is prohibited")
			}
		}
	})
	return nil, nil
}
