package main

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/mazzegi/scan"
)

type Command struct {
	Action                 string
	X0, X1, Y0, Y1, Z0, Z1 int
}

func (c Command) String() string {
	return fmt.Sprintf("%q: IX=[%d,%d], IY=[%d,%d], IZ=[%d,%d]", c.Action, c.X0, c.X1, c.Y0, c.Y1, c.Z0, c.Z1)
}

const input = `
on x=-46..2,y=-26..20,z=-39..5
on x=0..44,y=-44..0,z=-19..32
on x=-44..10,y=-20..28,z=4..48
on x=-12..39,y=-35..9,z=-12..36
on x=-20..31,y=-15..36,z=-44..2
on x=-18..36,y=-43..8,z=-12..41
on x=-24..27,y=-5..39,z=-30..24
on x=-40..4,y=-20..28,z=-33..14
on x=-16..30,y=-16..31,z=-12..32
on x=-3..48,y=-27..18,z=-5..39
off x=-37..-22,y=24..41,z=20..38
on x=-33..16,y=2..49,z=-46..4
off x=-9..2,y=29..41,z=-45..-32
on x=-7..44,y=-27..22,z=-40..6
off x=-43..-30,y=17..27,z=-43..-34
on x=-13..35,y=-19..25,z=-45..2
off x=31..45,y=36..49,z=12..28
on x=-16..33,y=-20..27,z=-34..16
off x=-43..-30,y=-1..17,z=2..13
on x=-7..47,y=-8..39,z=-2..44
on x=-51481..-16686,y=-55882..-41735,z=31858..57273
on x=62605..83371,y=-20326..7404,z=721..31101
on x=52622..66892,y=-53993..-42390,z=-32377..-15062
on x=32723..65065,y=-77593..-55363,z=-14922..15024
on x=-19522..-12277,y=-93093..-59799,z=23811..37917
on x=71595..92156,y=-44022..-19835,z=16005..34327
on x=-3586..34080,y=2703..21691,z=-92268..-61619
on x=-15386..715,y=73535..97744,z=-8553..16939
on x=-36770..-25752,y=71882..77753,z=-22483..5934
`

func main() {
	tpl, err := scan.ParseTemplate("command", "{{act: string}} x={{x0: int}}..{{x1: int}},y={{y0: int}}..{{y1: int}},z={{z0: int}}..{{z1: int}}")
	if err != nil {
		panic(err)
	}
	funcs := scan.BuiltinFuncs()
	scanner := bufio.NewScanner(bytes.NewBufferString(input))
	for scanner.Scan() {
		ln := strings.TrimSpace(scanner.Text())
		if ln == "" {
			continue
		}
		var cmd Command
		_, err := tpl.Eval(ln, funcs, &cmd.Action, &cmd.X0, &cmd.X1, &cmd.Y0, &cmd.Y1, &cmd.Z0, &cmd.Z1)
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			continue
		}
		fmt.Printf("%s\n", cmd)
	}
}
