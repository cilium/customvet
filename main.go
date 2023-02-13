package main

import (
	"github.com/cilium/customvet/analysis/ioreadall"
	"github.com/cilium/customvet/analysis/timeafter"
	"github.com/cilium/customvet/analysis/typeassertion"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(timeafter.Analyzer, ioreadall.Analyzer, typeassertion.Analyzer)
}
