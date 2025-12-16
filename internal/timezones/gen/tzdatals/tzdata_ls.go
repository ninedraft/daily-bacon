// This tool finds golang generated tzdata in zip format and extracts list of locations to be embedded into this project.
package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

const tzdataFileMode = 0o600

func main() {
	status, err := run()
	if err != nil {
		log.Printf("[ERROR]: %v", err)
	}
	os.Exit(status)
}

func run() (int, error) {
	flag.Parse()

	outputFilename := flag.Arg(0)
	if outputFilename == "" {
		return 1, errors.New("missing output filename argument")
	}

	zzipdataFilepath := filepath.Join(os.Getenv("GOROOT"), "src/time/tzdata/zzipdata.go")
	log.Printf("parsing tzdata file: %v", zzipdataFilepath)

	zipdata, err := extractZipdata(zzipdataFilepath)
	if err != nil {
		return 1, fmt.Errorf("extract zipdata: %w", err)
	}

	tzdata, err := zip.NewReader(
		bytes.NewReader([]byte(zipdata)),
		int64(len(zipdata)),
	)
	if err != nil {
		return 1, fmt.Errorf("open zip reader: %w", err)
	}

	var timezones []string
	err = fs.WalkDir(tzdata, ".", func(fpath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || path.Ext(d.Name()) != "" {
			timezones = append(timezones, fpath)
		}
		return nil
	})
	if err != nil {
		log.Printf("[ERROR]: walking tzdata %v", err)
	}

	slices.Sort(timezones)
	data := strings.Join(timezones, "\n")

	if err := os.WriteFile(outputFilename, []byte(data), tzdataFileMode); err != nil {
		return 1, fmt.Errorf("write file %s: %w", outputFilename, err)
	}

	return 0, nil
}

// extractZipdata parses a Go file, finds const zipdata, and returns its evaluated string value.
func extractZipdata(filePath string) (string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return "", fmt.Errorf("parse file %s: %w", filePath, err)
	}

	// Find the expression assigned to const zipdata
	targetExpr, err := findZipdataExpr(file)
	if err != nil {
		return "", err
	}

	// Type-check and evaluate the constant expression.
	conf := types.Config{
		Importer: importer.Default(),
		Error: func(err error) {
			log.Printf("[ERROR]: %v", err)
		},
	}

	info := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	_, _ = conf.Check("tzdata", fset, []*ast.File{file}, info)

	tv, ok := info.Types[targetExpr]
	if !ok {
		return "", errors.New("cannot evaluate zipdata: missing type/value info")
	}
	if tv.Value == nil {
		return "", errors.New("zipdata is not a constant expression or has no constant value")
	}
	if !isString(tv.Type) {
		return "", fmt.Errorf("zipdata has non-string type: %v", tv.Type)
	}

	s := constant.StringVal(tv.Value)
	return s, nil
}

func isString(ttype types.Type) bool {
	if ttype == nil {
		return false
	}

	basic, ok := ttype.(*types.Basic)
	if !ok || basic == nil {
		return false
	}

	return basic.Kind() == types.UntypedString || basic.Kind() == types.String
}

func findZipdataExpr(f *ast.File) (ast.Expr, error) {
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok || gen.Tok != token.CONST {
			continue
		}
		for _, spec := range gen.Specs {
			vs, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// Build mapping from name index to the expression used.
			// Cases:
			//   const a, b = "x", "y"        -> len(Names)==len(Values)
			//   const a, b = "xy"            -> one value for all names (b gets same)
			for i, name := range vs.Names {
				if name == nil || name.Name != "zipdata" {
					continue
				}
				switch {
				case len(vs.Values) == len(vs.Names):
					return vs.Values[i], nil
				case len(vs.Values) == 1:
					return vs.Values[0], nil
				default:
					// Handles iota-like or uninitialized consts (not applicable here)
					return nil, errors.New("zipdata found, but value expression is not explicitly provided")
				}
			}
		}
	}
	return nil, errors.New(`const "zipdata" not found in the file`)
}
