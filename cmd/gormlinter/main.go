package main

import (
	"github.com/lema.ai/lemmata/tools/gormlinter/internal/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(analyzer.Analyzer)
}
