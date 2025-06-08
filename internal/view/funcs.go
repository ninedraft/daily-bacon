package view

import (
	"fmt"
	"strings"
	"text/template"
)

var funcs = template.FuncMap{
	"dict":        dict,
	"formatFloat": formatFloat,
	"joinFloats":  joinFloats,
	"replace":     strings.ReplaceAll,
}

// dict builds a map from a sequence of key/value pairs.
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict requires an even number of args; got %d", len(values))
	}
	m := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key := fmt.Sprint(values[i])
		m[key] = values[i+1]
	}
	return m, nil
}

// formatFloat prints a float with two decimals, trimming trailing zeros.
func formatFloat(f float64) string {
	// e.g. 1.00 → "1"; 1.20 → "1.2"; 1.23 → "1.23"
	s := fmt.Sprintf("%.2f", f)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// joinFloats joins a slice of floats into lines of formatted numbers.
func joinFloats(vals []float64) string {
	var b strings.Builder
	for i, v := range vals {
		if i > 0 {
			b.WriteString("\n")
		}
		b.WriteString(formatFloat(v))
	}
	return b.String()
}
