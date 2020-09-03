package main

import (
	"github.com/shibh308/lowercasefunc"
	"golang.org/x/tools/go/analysis/unitchecker"
)

func main() { unitchecker.Main(lowercasefunc.Analyzer) }

