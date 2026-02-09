package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/lema-ai/rlslinter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("rlslinter", New)
}

// New creates a new instance of the gormlinter plugin
func New(settings any) (register.LinterPlugin, error) {
	return &RLSLinter{}, nil
}

// RLSLinter implements the LinterPlugin interface
type RLSLinter struct{}

// BuildAnalyzers returns the analyzers provided by this linter
func (l *RLSLinter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

// GetLoadMode returns the load mode required by this linter
func (l *RLSLinter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
