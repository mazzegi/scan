package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/mazzegi/scan"
)

type Point struct {
	X, Y, Z float64
}

func (p Point) String() string {
	return fmt.Sprintf("[%.3f, %.3f, %.3f]", p.X, p.Y, p.Z)
}

var input = `
p1: 1.2, -5.7, 8.5
p2: 6.6, 51.8, -12.3
p2: -34.45534, 4e-3, 21
`

func main() {
	tplLine, err := scan.ParseTemplate("line", "{{name: string}}:{{point: point}}")
	if err != nil {
		panic(err)
	}
	tplPoint, err := scan.ParseTemplate("point", "{{x:float}},{{y:float}},{{z:float}}")
	if err != nil {
		panic(err)
	}
	funcs := scan.BuiltinFuncs()
	funcs.Add("point", func(s string) (any, error) {
		var p Point
		res, err := tplPoint.Eval(s, funcs)
		if err != nil {
			return nil, err
		}
		err = res.Scan(&p.X, &p.Y, &p.Z)
		if err != nil {
			return nil, err
		}
		return p, nil
	})

	scanner := bufio.NewScanner(bytes.NewBufferString(input))
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		if ln == "" {
			continue
		}
		var name string
		var pt Point
		res, err := tplLine.Eval(ln, funcs)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		err = res.Scan(&name, &pt)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%q => %s\n", name, pt)
	}
}
