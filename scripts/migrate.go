package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"ariga.io/atlas-provider-gorm/gormschema"
	"github.com/dukk308/golang-clean-arch-starter/database"
)

func main() {
	sb := &strings.Builder{}
	loadModels(sb)

	io.WriteString(os.Stdout, sb.String()) //nolint
}

func loadModels(sb *strings.Builder) {
	models := database.Models

	stmts, err := gormschema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err) // nolint: revive
		os.Exit(1)
	}

	clean := strings.TrimSpace(stmts)
	clean = strings.TrimSuffix(clean, ";")

	sb.WriteString(clean) // nolint: revive
	sb.WriteString(";\n") // nolint: revive
}
