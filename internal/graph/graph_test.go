package graph

import (
	"errors"
	"reflect"
	"testing"
)

// Grafo da lista 01: V = {0..5}.
var sampleEdges = []Edge{{0, 1}, {0, 2}, {1, 3}, {1, 4}, {2, 4}, {3, 5}, {4, 5}}

// constructors permite rodar os mesmos testes nas duas representações.
var constructors = map[string]func(n int) Graph{
	"matriz": func(n int) Graph { return NewAdjacencyMatrix(n) },
	"lista":  func(n int) Graph { return NewAdjacencyList(n) },
}

func TestNeighbors(t *testing.T) {
	want := [][]int{
		{1, 2},
		{0, 3, 4},
		{0, 4},
		{1, 5},
		{1, 2, 5},
		{3, 4},
	}

	for name, newGraph := range constructors {
		t.Run(name, func(t *testing.T) {
			g, err := Build(newGraph, 6, sampleEdges)
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if g == nil {
				t.Fatal("Build retornou nil")
			}
			for v, w := range want {
				if got := g.Neighbors(v); !reflect.DeepEqual(got, w) {
					t.Errorf("Neighbors(%d) = %v, quero %v", v, got, w)
				}
			}
		})
	}
}

// A aresta (2,5) é inserida por último (Questão 10): os vizinhos
// precisam continuar em ordem crescente.
func TestNeighborsSortedAfterLateInsert(t *testing.T) {
	edges := append(append([]Edge{}, sampleEdges...), Edge{2, 5})

	for name, newGraph := range constructors {
		t.Run(name, func(t *testing.T) {
			g, err := Build(newGraph, 6, edges)
			if err != nil || g == nil {
				t.Fatalf("Build: g=%v err=%v", g, err)
			}
			if got, want := g.Neighbors(5), []int{2, 3, 4}; !reflect.DeepEqual(got, want) {
				t.Errorf("Neighbors(5) = %v, quero %v", got, want)
			}
		})
	}
}

func TestHasEdgeIsSymmetric(t *testing.T) {
	for name, newGraph := range constructors {
		t.Run(name, func(t *testing.T) {
			g, err := Build(newGraph, 6, sampleEdges)
			if err != nil || g == nil {
				t.Fatalf("Build: g=%v err=%v", g, err)
			}
			for _, e := range sampleEdges {
				if !g.HasEdge(e.U, e.V) || !g.HasEdge(e.V, e.U) {
					t.Errorf("aresta %v deveria existir nos dois sentidos", e)
				}
			}
			if g.HasEdge(0, 5) {
				t.Error("HasEdge(0, 5) deveria ser false")
			}
		})
	}
}

func TestAddEdgeInvalidVertex(t *testing.T) {
	for name, newGraph := range constructors {
		t.Run(name, func(t *testing.T) {
			g := newGraph(6)
			if err := g.AddEdge(0, 6); !errors.Is(err, ErrInvalidVertex) {
				t.Errorf("AddEdge(0, 6) = %v, quero ErrInvalidVertex", err)
			}
		})
	}
}

func TestMatrixString(t *testing.T) {
	g, err := Build(func(n int) Graph { return NewAdjacencyMatrix(n) }, 6, sampleEdges)
	if err != nil || g == nil {
		t.Fatalf("Build: g=%v err=%v", g, err)
	}
	want := "" +
		"0 1 1 0 0 0\n" +
		"1 0 0 1 1 0\n" +
		"1 0 0 0 1 0\n" +
		"0 1 0 0 0 1\n" +
		"0 1 1 0 0 1\n" +
		"0 0 0 1 1 0\n"
	if got := g.String(); got != want {
		t.Errorf("String() =\n%s\nquero\n%s", got, want)
	}
}

func TestAddEdgeSelfLoop(t *testing.T) {
	for name, newGraph := range constructors {
		t.Run(name, func(t *testing.T) {
			if err := newGraph(3).AddEdge(1, 1); !errors.Is(err, ErrSelfLoop) {
				t.Errorf("AddEdge(1, 1) = %v, quero ErrSelfLoop", err)
			}
		})
	}
}

func TestSizeAndEdges(t *testing.T) {
	for name, newGraph := range constructors {
		t.Run(name, func(t *testing.T) {
			g, err := Build(newGraph, 6, append(append([]Edge{}, sampleEdges...), Edge{1, 0}))
			if err != nil {
				t.Fatalf("Build: %v", err)
			}
			if got := Size(g); got != len(sampleEdges) {
				t.Errorf("Size = %d, quero %d (aresta repetida não conta)", got, len(sampleEdges))
			}
			if got := Edges(g); !reflect.DeepEqual(got, sampleEdges) {
				t.Errorf("Edges = %v, quero %v", got, sampleEdges)
			}
		})
	}
}

func TestDigraphIsDirected(t *testing.T) {
	d, err := BuildDigraph(3, []Edge{{0, 1}, {2, 1}})
	if err != nil {
		t.Fatalf("BuildDigraph: %v", err)
	}
	if !d.HasArc(0, 1) || d.HasArc(1, 0) {
		t.Error("arco 0 → 1 não deveria implicar 1 → 0")
	}
	if got := d.Successors(1); len(got) != 0 {
		t.Errorf("Successors(1) = %v, quero []", got)
	}
}
