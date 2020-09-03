package lowercasefunc_test

import (
	"testing"

	"github.com/shibh308/lowercasefunc"
	"golang.org/x/tools/go/analysis/analysistest"
)

// TestAnalyzer is a test for Analyzer.
func TestAnalyzer(t *testing.T) {

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, lowercasefunc.Analyzer, "a")
}

