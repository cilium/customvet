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

package testdata

import "fmt"

type A struct{}

func (a A) foo() string {
	return "foo"
}

func TypeAssertionAssignmentNoCheck() {
	a := A{}
	var ai interface{} = a
	b := ai.(A) // want `use of an unchecked type assertion is prohibited`
	fmt.Println(b)
}

func TypeAssertionCallNoCheck() {
	a := A{}
	var ai interface{} = a
	b := ai.(A).foo() // want `use of an unchecked type assertion is prohibited`
	fmt.Println(b)
}

func TypeAssertionAssignmentWithCheck() {
	a := A{}
	var ai interface{} = a
	b, ok := ai.(A)
	fmt.Println(b, ok)
}

func TypeAssertionAssignmentWithExplicitNoCheck() {
	a := A{}
	var ai interface{} = a
	b, _ := ai.(A)
	fmt.Println(b)
}
