package plugin

import (
	"github.com/golangci/plugin-module-register/register"
	"github.com/lema-ai/rlslinter/internal/analyzer"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("gormlinter", New)
}

// New creates a new instance of the gormlinter plugin
func New(settings any) (register.LinterPlugin, error) {
	return &GormLinter{}, nil
}

// GormLinter implements the LinterPlugin interface
type GormLinter struct{}

// BuildAnalyzers returns the analyzers provided by this linter
func (l *GormLinter) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{analyzer.Analyzer}, nil
}

// GetLoadMode returns the load mode required by this linter
func (l *GormLinter) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
