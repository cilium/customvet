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

import (
	"fmt"
	"time"
)

func TimeAfterNotInForLoop() {
	<-time.After(time.Millisecond)
}

func TimeAfterNotInForLoop2() {
	select {
	case <-time.After(time.Millisecond):
	}
}

func TimeAfterInForLoop() {
	for i, l := 0, 3; i < l; i++ {
		fmt.Printf("time_after_perumutation_%d", i)
		<-time.After(time.Millisecond) // want `use of time.After in a for loop is prohibited, use inctimer instead`
	}
}

func TimeAfterInForRangeLoop() {
	for _, n := range []int{0, 1, 2} {
		fmt.Printf("time_after_perumutation_%d", n)
		<-time.After(time.Millisecond) // want `use of time.After in a for loop is prohibited, use inctimer instead`
	}
}
