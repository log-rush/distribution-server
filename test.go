package main

import "fmt"

type X struct {
	n int
}

type Y struct {
	n int
}

func (x *X) Inc() {
	x.n += 1
}

func (y Y) Inc() {
	y.n += 1
}

type Combi struct {
	a *X
	b X
}

func (c *Combi) IncA() {
	c.a.Inc()
	c.b.Inc()
}

func (c Combi) IncB() {
	c.a.Inc()
	c.b.Inc()
}

func main() {
	a := X{}
	b := &X{}
	e := a
	f := b
	c := Y{}
	d := &Y{}

	a.Inc()
	b.Inc()
	c.Inc()
	d.Inc()
	fmt.Println(a, b, c, d, e, f)

	g := Combi{&X{}, X{}}
	h := &Combi{&X{}, X{}}

	r := g.a
	s := g.b
	t := h.a
	u := h.b
	g.a.Inc()
	g.b.Inc()
	h.a.Inc()
	h.b.Inc()

	fmt.Println(g, h, r, s, t, u)
	g.IncA()
	h.IncA()

	fmt.Println(g, h, r, s, t, u)
	g.IncB()
	h.IncB()
	fmt.Println(g, h, r, s, t, u)
}
