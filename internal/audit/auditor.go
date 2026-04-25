// Package audit provides environment file auditing capabilities,
// combining diff analysis, secret detection, and schema validation
// into a unified audit report summary.
package audit

import (
	"github.com/user/envlens/internal/differ"
	"github.com/user/envlens/internal/parser"
	"github.com/user/envlens/internal/schema"
	"github.com/user/envlens/internal/secrets"
)

// Result holds the combined output of an audit run.
type Result struct {
	Diffs    []differ.DiffEntry
	Findings []secrets.Finding
	Violations []schema.Violation
	BaseFile  string
	TargetFile string
}

// Auditor orchestrates diff, secret detection, and schema validation.
type Auditor struct {
	detector *secrets.Detector
	schema   *schema.Schema
}

// NewAuditor creates an Auditor. schema may be nil to skip validation.
func NewAuditor(detector *secrets.Detector, s *schema.Schema) *Auditor {
	return &Auditor{detector: detector, schema: s}
}

// Run performs a full audit comparing base against target env files.
func (a *Auditor) Run(baseFile, targetFile string) (*Result, error) {
	baseEntries, err := parser.Parse(baseFile)
	if err != nil {
		return nil, err
	}

	targetEntries, err := parser.Parse(targetFile)
	if err != nil {
		return nil, err
	}

	diffs := differ.Diff(baseEntries, targetEntries)

	findings := a.detector.Scan(targetEntries)

	var violations []schema.Violation
	if a.schema != nil {
		violations = a.schema.Validate(targetEntries)
	}

	return &Result{
		Diffs:      diffs,
		Findings:   findings,
		Violations: violations,
		BaseFile:   baseFile,
		TargetFile: targetFile,
	}, nil
}

// HasIssues returns true if the audit result contains any findings or violations.
func (r *Result) HasIssues() bool {
	return len(r.Findings) > 0 || len(r.Violations) > 0
}
