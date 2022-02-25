package main

import (
	"bytes"
	"fmt"

	"github.com/mazzegi/scan"
)

type Point struct {
	X, Y, Z float64
}

func (p Point) String() string {
	return fmt.Sprintf("[%.3f, %.3f, %.3f]", p.X, p.Y, p.Z)
}

type NamedPoint struct {
	Name  string
	Point Point
}

var input = `
p1: 1.2, -5.7, 8.5
p2: 6.6, 51.8, -12.3
p2: -34.45534, 4e-3, 21
`

func main() {

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

	ps, err := scan.Lines[NamedPoint]("{{name: string}}:{{point: point}}", funcs, bytes.NewBufferString(input))
	if err != nil {
		panic(err)
	}
	for _, p := range ps {
		fmt.Printf("%q => %s\n", p.Name, p.Point)
	}
}
