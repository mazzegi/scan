package scan

import (
	"strings"

	"github.com/pkg/errors"
)

type Item interface{}

type Template struct {
	name  string
	items []Item
}

func (t *Template) Name() string {
	return t.name
}

func (t *Template) Prefix() string {
	if len(t.items) == 0 {
		return ""
	}
	if s, ok := t.items[0].(string); ok {
		return s
	}
	return ""
}

func (t *Template) appendItem(item Item) error {
	if _, ok := item.(*Evaler); ok {
		if len(t.items) > 0 {
			if _, ok := t.items[len(t.items)-1].(*Evaler); ok {
				return errors.Errorf("a match-target cannot immediately follwo a match-target")
			}
		}
	}
	t.items = append(t.items, item)
	return nil
}

//

func ParseTemplate(name string, s string) (*Template, error) {
	items, err := newItemsParser(s).parse()
	if err != nil {
		return nil, err
	}
	//check items
	lastWasEvaler := false
	for _, item := range items {
		switch item.(type) {
		case *Evaler:
			if lastWasEvaler {
				return nil, errors.Errorf("an evaler cannot immediately follow an evaler")
			}
			lastWasEvaler = true
		default:
			lastWasEvaler = false
		}
	}

	return &Template{
		name:  name,
		items: items,
	}, nil
}

//
type itemParseFunc func() (itemParseFunc, error)

type itemsParser struct {
	rs    []rune
	pos   int
	items []Item
}

func newItemsParser(s string) *itemsParser {
	p := &itemsParser{
		rs:    []rune(s),
		pos:   0,
		items: []Item{},
	}
	return p
}

func (p *itemsParser) parse() ([]Item, error) {
	fnc := p.parseText
	var err error
	for fnc != nil {
		fnc, err = fnc()
		if err != nil {
			return nil, err
		}
	}
	return p.items, nil
}

func (p *itemsParser) parseText() (itemParseFunc, error) {
	var text string
	defer func() {
		text = strings.TrimSpace(text)
		if text == "" {
			return
		}
		p.items = append(p.items, text)
	}()

	for {
		if p.pos >= len(p.rs) {
			return nil, nil
		}
		if strings.HasPrefix(string(p.rs[p.pos:]), "{{") {
			p.pos += 2
			return p.parseEvaler, nil
		}
		text += string(p.rs[p.pos])
		p.pos++
	}
}

func (p *itemsParser) parseEvaler() (itemParseFunc, error) {
	idx := strings.Index(string(p.rs[p.pos:]), "}}")
	if idx < 0 {
		return nil, errors.Errorf("no closing }} found")
	}
	sub := p.rs[p.pos : p.pos+idx]

	ev, err := ParseEvaler(string(sub))
	if err != nil {
		return nil, err
	}
	p.items = append(p.items, ev)

	p.pos += idx + 2
	return p.parseText, nil
}
