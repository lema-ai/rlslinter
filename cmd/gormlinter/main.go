package main

import (
	"github.com/lema-ai/rlslinter/internal/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
