package graph

import (
	"fmt"
	"slices"
	"strings"
)

// Digraph é um grafo direcionado representado por lista de adjacência:
// out[u] contém apenas os vértices v tais que existe o arco u → v.
type Digraph struct {
	out [][]int
}

func NewDigraph(n int) *Digraph {
	return &Digraph{out: make([][]int, n)}
}

func BuildDigraph(n int, arcs []Edge) (*Digraph, error) {
	d := NewDigraph(n)
	for _, a := range arcs {
		if err := d.AddArc(a.U, a.V); err != nil {
			return nil, fmt.Errorf("arco %d → %d: %w", a.U, a.V, err)
		}
	}
	return d, nil
}

func (d *Digraph) Order() int {
	return len(d.out)
}

// AddArc insere o arco u → v. Diferente do grafo não direcionado, o
// sentido oposto NÃO é registrado.
func (d *Digraph) AddArc(u, v int) error {
	if err := validate(d.Order(), u, v); err != nil {
		return err
	}
	d.out[u] = insertSorted(d.out[u], v)
	return nil
}

func (d *Digraph) HasArc(u, v int) bool {
	if validate(d.Order(), u, v) != nil {
		return false
	}
	_, found := slices.BinarySearch(d.out[u], v)
	return found
}

// Successors retorna, em ordem crescente, os vértices alcançados pelos
// arcos que saem de u.
func (d *Digraph) Successors(u int) []int {
	return slices.Clone(d.out[u])
}

func (d *Digraph) String() string {
	var b strings.Builder
	for i, succ := range d.out {
		fmt.Fprintf(&b, "%d → %s\n", i, FormatInts(succ))
	}
	return b.String()
}
