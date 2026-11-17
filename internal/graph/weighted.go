package graph

import (
	"fmt"
	"slices"
	"strings"
)

// Arc é uma entrada da lista ponderada: o vizinho To e o peso da ligação.
type Arc struct {
	To, Weight int
}

// Weighted é uma lista de adjacência ponderada: adj[u] guarda pares (v, peso).
// Se directed for falso, cada aresta é registrada nas duas listas.
type Weighted struct {
	directed bool
	adj      [][]Arc
}

var _ Neighborhood = (*Weighted)(nil)

func NewWeighted(n int, directed bool) *Weighted {
	// PASSO 1 – criação: uma lista vazia de pares (vizinho, peso) por vértice.
	return &Weighted{directed: directed, adj: make([][]Arc, n)}
}

func BuildWeighted(n int, directed bool, edges []WeightedEdge) (*Weighted, error) {
	g := NewWeighted(n, directed)
	for _, e := range edges {
		if err := g.AddEdge(e.U, e.V, e.Weight); err != nil {
			return nil, fmt.Errorf("aresta %d-%d: %w", e.U, e.V, err)
		}
	}
	return g, nil
}

// WeightedEdge é a aresta {U, V} (ou o arco U → V) com seu peso.
type WeightedEdge struct {
	U, V, Weight int
}

func (g *Weighted) Order() int {
	return len(g.adj)
}

func (g *Weighted) Directed() bool {
	return g.directed
}

func (g *Weighted) AddEdge(u, v, w int) error {
	if err := validate(g.Order(), u, v); err != nil {
		return err
	}
	// PASSO 2 – inserção: (v, peso) entra na lista de u. Sem direção, (u, peso)
	// entra também na lista de v; com direção, não há arco reverso automático.
	g.adj[u] = insertArc(g.adj[u], Arc{To: v, Weight: w})
	if !g.directed {
		g.adj[v] = insertArc(g.adj[v], Arc{To: u, Weight: w})
	}
	return nil
}

// insertArc mantém a lista ordenada pelo vizinho; uma ligação repetida só
// atualiza o peso.
func insertArc(arcs []Arc, a Arc) []Arc {
	i, found := slices.BinarySearchFunc(arcs, a.To, func(x Arc, to int) int { return x.To - to })
	if found {
		arcs[i] = a
		return arcs
	}
	return slices.Insert(arcs, i, a)
}

// Arcs devolve uma cópia dos pares (vizinho, peso) de u, em ordem crescente.
func (g *Weighted) Arcs(u int) []Arc {
	return slices.Clone(g.adj[u])
}

// Neighbors descarta os pesos: BFS e DFS só olham a conectividade.
func (g *Weighted) Neighbors(u int) []int {
	out := make([]int, len(g.adj[u]))
	for i, a := range g.adj[u] {
		out[i] = a.To
	}
	return out
}

// Weight retorna o peso de u → v e se essa ligação existe.
func (g *Weighted) Weight(u, v int) (int, bool) {
	if validate(g.Order(), u, v) != nil {
		return 0, false
	}
	for _, a := range g.adj[u] {
		if a.To == v {
			return a.Weight, true
		}
	}
	return 0, false
}

func (g *Weighted) String() string {
	// PASSO 3 – exibição no formato "0 → [(1, 4), (2, 2)]".
	var b strings.Builder
	for u, arcs := range g.adj {
		parts := make([]string, len(arcs))
		for i, a := range arcs {
			parts[i] = fmt.Sprintf("(%d, %d)", a.To, a.Weight)
		}
		fmt.Fprintf(&b, "%d → [%s]\n", u, strings.Join(parts, ", "))
	}
	return b.String()
}
