package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `rlslinter detects unsupported GORM methods that fail with the org-wrapping transaction system.

The lemmata GORM wrapper automatically wraps queries in transactions and commits them immediately.
This breaks methods that return data to be read after the transaction closes, causing
Row-Level Security (RLS) to fail and potentially returning empty or incorrect results.
`

var Analyzer = &analysis.Analyzer{
	Name:     "rlslinter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

var prohibitedMethods = map[string]methodInfo{
	"Row": {
		message:     "GORM .Row() is not supported - data accessed after transaction closes",
		replacement: "Use .First() or .Find() instead",
		explanation: `
Why this fails:
  • The GORM wrapper commits transactions immediately after query execution
  • .Row() returns data to be read AFTER the transaction closes
  • app.current_org_id is no longer set when data is accessed
  • Row-Level Security (RLS) won't filter by organization
  • Queries on org tables may return EMPTY results or data from wrong org

To suppress (ONLY if querying non-org tables):
  //nolint:rlslinter // Not querying org tables
  row := db.Row()`,
	},
	"Rows": {
		message:     "GORM .Rows() is not supported - data accessed after transaction closes",
		replacement: "Use .Find() to retrieve multiple records",
		explanation: `
Why this fails:
  • The GORM wrapper commits transactions immediately after query execution
  • .Rows() returns data to be read AFTER the transaction closes
  • app.current_org_id is no longer set when data is accessed
  • Row-Level Security (RLS) won't filter by organization
  • Queries on org tables may return EMPTY results or data from wrong org

To suppress (ONLY if querying non-org tables):
  //nolint:rlslinter // Not querying org tables
  rows, _ := db.Rows()`,
	},
	"Scan": {
		message:     "GORM .Scan() is not supported - data accessed after transaction closes",
		replacement: "Use .Pluck() for single column or .Find() for multiple columns",
		explanation: `
Why this fails:
  • The GORM wrapper commits transactions immediately after query execution
  • .Scan() returns data to be read AFTER the transaction closes
  • app.current_org_id is no longer set when data is accessed
  • Row-Level Security (RLS) won't filter by organization
  • Queries on org tables may return EMPTY results or data from wrong org

To suppress (ONLY if querying non-org tables):
  //nolint:rlslinter // Not querying org tables
  db.Scan(&result)`,
	},
}

type methodInfo struct {
	message     string
	replacement string
	explanation string
}

func run(pass *analysis.Pass) (any, error) {
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}

	insp.Preorder(nodeFilter, func(n ast.Node) {
		callExpr := n.(*ast.CallExpr)

		// Check if this is a selector expression (e.g., db.Row())
		selector, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			return
		}

		methodName := selector.Sel.Name

		// Check if method is in our prohibited list
		methodDetails, isProhibited := prohibitedMethods[methodName]
		if !isProhibited {
			return
		}

		// Check if receiver is a GORM type
		if !isGormType(pass.TypesInfo, selector.X) {
			return
		}

		// Check for suppression comment
		if hasSuppression(pass, callExpr) {
			return
		}

		// Report the issue
		reportIssue(pass, callExpr, methodName, methodDetails)
	})

	return nil, nil
}

// isGormType checks if the expression's type is from GORM
func isGormType(info *types.Info, expr ast.Expr) bool {
	typ := info.TypeOf(expr)
	if typ == nil {
		return false
	}

	// Handle pointer types
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	// Get the named type
	named, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	obj := named.Obj()
	if obj == nil || obj.Pkg() == nil {
		return false
	}

	pkgPath := obj.Pkg().Path()

	// Check if it's from GORM or generated GORM queries
	return pkgPath == "gorm.io/gorm" ||
		pkgPath == "gorm.io/gorm/clause" ||
		strings.HasPrefix(pkgPath, "github.com/lema.ai/lemmata/db/gorm") ||
		strings.HasPrefix(pkgPath, "github.com/lema-ai/lemmata/db/gorm")
}

// hasSuppression checks if there's a nolint comment suppressing this linter
func hasSuppression(pass *analysis.Pass, node ast.Node) bool {
	file := pass.Fset.File(node.Pos())
	if file == nil {
		return false
	}

	// Get line number of the call
	line := file.Line(node.Pos())

	// Look for comment on same line or line above
	for _, commentGroup := range pass.Files[0].Comments {
		for _, comment := range commentGroup.List {
			commentLine := file.Line(comment.Pos())

			// Check if comment is on the same line or line above
			if commentLine == line || commentLine == line-1 {
				text := comment.Text
				// Check for //nolint:rlslinter or //nolint
				if strings.Contains(text, "nolint:rlslinter"); strings.Contains(text, "nolint") && !strings.Contains(text, "nolint:") {
					return true
				}
			}
		}
	}

	return false
}

// reportIssue reports a diagnostic for a prohibited GORM method call
func reportIssue(pass *analysis.Pass, callExpr *ast.CallExpr, methodName string, info methodInfo) {
	selector := callExpr.Fun.(*ast.SelectorExpr)

	var pos token.Pos
	var end token.Pos

	// Position at the method name itself
	pos = selector.Sel.Pos()
	end = selector.Sel.End()

	message := info.message + "\n" +
		info.replacement + "\n" +
		info.explanation

	pass.Report(analysis.Diagnostic{
		Pos:     pos,
		End:     end,
		Message: message,
	})
}
