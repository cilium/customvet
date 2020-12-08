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
)

const (
	timeAfterPkg  = "time"
	timeAfterFunc = "After"
)

type timeAfterCheck struct {
	aliases []string
}

func (t *timeAfterCheck) check(n ast.Node) (ok bool, report string, pos token.Pos) {
	ok = true
	switch stmt := n.(type) {
	case *ast.File:
		t.aliases = []string{timeAfterPkg}
	case *ast.ImportSpec:
		// Collect aliases.
		pkg := stmt.Path.Value
		if pkg == fmt.Sprintf("%q", timeAfterPkg) {
			if stmt.Name != nil {
				t.aliases = append(t.aliases, stmt.Name.Name)
			}
		}
	case *ast.ForStmt:
		ast.Walk(visitor(func(node ast.Node) bool {
			switch expr := node.(type) {
			case *ast.CallExpr:
				for _, pkg := range t.aliases {
					if isPkgDot(expr.Fun, pkg, timeAfterFunc) {
						ok = false
						report = fmt.Sprintf("use of %s.After in a for loop is prohibited, use inctimer instead", pkg)
						pos = node.Pos()
					}
				}
			}
			return ok
		}), stmt.Body)
	}
	return
}
