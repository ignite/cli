//go:build ignore

//go:generate go run gen.go

// Code generator for math/overflow.
package main

import (
	"os"
	"strconv"
	"text/template"

	_ "embed"
)

func main() {
	file, err := os.Create("overflow_generated.gno")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	types := make([]Type, 0, 10)
	for _, bs := range [...]int{0, 8, 16, 32, 64} {
		var s string
		if bs != 0 {
			s = strconv.Itoa(bs)
		} else {
			bs = 64
		}
		types = append(
			types,
			Type{
				Name:   "int" + s,
				Short:  s,
				Signed: true,
				Min:    -(1 << (bs - 1)),
			},
			Type{
				Name:  "uint" + s,
				Short: "u" + s,
			},
		)
	}
	if err := parsedTmpl.Execute(file, types); err != nil {
		panic(err)
	}
}

type Type struct {
	Name   string
	Short  string
	Signed bool
	Min    int64
}

var parsedTmpl = template.Must(template.New("").Parse(tmpl))

//go:embed template.go.tpl
var tmpl string
