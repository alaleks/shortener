package checkosexit_test

import (
	"testing"

	"github.com/alaleks/shortener/pkg/checkosexit"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestCheckOsExit(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), checkosexit.Analyzer, "./...")
}
