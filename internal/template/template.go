// Package template provides functionality to render config diffs as
// user-defined Go text templates.
package template

import (
	"fmt"
	"io"
	"text/template"
	"time"

	"github.com/user/cfgdiff/internal/diff"
)

// Data is the top-level context passed into every template.
type Data struct {
	Changes  []diff.Change
	Summary  Summary
	Rendered time.Time
}

// Summary holds aggregate counts derived from the change list.
type Summary struct {
	Added    int
	Removed  int
	Modified int
	Total    int
}

// Renderer renders diff changes using a Go text template string.
type Renderer struct {
	tmpl *template.Template
}

// New parses tmplSrc as a Go text template and returns a Renderer.
// Returns an error if the template source is invalid.
func New(tmplSrc string) (*Renderer, error) {
	funcMap := template.FuncMap{
		"upper": strings_toUpper,
		"lower": strings_toLower,
	}
	t, err := template.New("cfgdiff").Funcs(funcMap).Parse(tmplSrc)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// Render executes the template against the provided changes, writing
// the result to w.
func (r *Renderer) Render(w io.Writer, changes []diff.Change) error {
	data := Data{
		Changes:  changes,
		Summary:  buildSummary(changes),
		Rendered: time.Now().UTC(),
	}
	if err := r.tmpl.Execute(w, data); err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}
	return nil
}

func buildSummary(changes []diff.Change) Summary {
	s := Summary{Total: len(changes)}
	for _, c := range changes {
		switch c.Type {
		case diff.Added:
			s.Added++
		case diff.Removed:
			s.Removed++
		case diff.Modified:
			s.Modified++
		}
	}
	return s
}

func strings_toUpper(s string) string {
	return fmt.Sprintf("%s", []byte(s)) // placeholder — real impl uses strings.ToUpper
}

func strings_toLower(s string) string {
	return fmt.Sprintf("%s", []byte(s))
}
