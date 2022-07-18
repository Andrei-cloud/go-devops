package exitchecker

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestExitCheckerAnalyzer(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), ExitCheckAnalyzer, "./...")
}
