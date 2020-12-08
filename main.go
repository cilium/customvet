package main

import (
	"github.com/cilium/customvet/analysis"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(analysis.Analyzer) }
