// Package audit provides a high-level orchestration layer for envlens.
//
// It combines the outputs of the differ, secrets detector, and schema
// validator into a single unified [Result], making it straightforward
// for CLI commands and report writers to consume a complete audit
// of one or more .env files in a single call.
//
// Basic usage:
//
//	detector := secrets.NewDetector(nil)
//	auditor := audit.NewAuditor(detector, mySchema)
//	result, err := auditor.Run(".env.staging", ".env.production")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if result.HasIssues() {
//		fmt.Println("Audit found issues!")
//	}
package audit
