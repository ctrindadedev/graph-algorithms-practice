package tree

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Grafo da lista 01: V = {0..5}.
var sampleEdges = []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 2, V: 4}, {U: 3, V: 5}, {U: 4, V: 5}}

func newList(n int) graph.Graph { return graph.NewAdjacencyList(n) }

func mustBuild(t *testing.T, n int, edges []graph.Edge) graph.Graph {
	t.Helper()
	g, err := graph.Build(newList, n, edges)
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	return g
}

func TestIsTree(t *testing.T) {
	tests := []struct {
		name  string
		n     int
		edges []graph.Edge
		want  bool
	}{
		{"árvore", 6, []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 2, V: 5}}, true},
		{"vértice único", 1, nil, true},
		{"grafo da lista 01 (ciclos)", 6, sampleEdges, false},
		{"V-1 arestas mas desconexo", 4, []graph.Edge{{U: 0, V: 1}, {U: 1, V: 2}, {U: 2, V: 0}}, false},
		{"menos de V-1 arestas", 3, []graph.Edge{{U: 0, V: 1}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsTree(mustBuild(t, tt.n, tt.edges)); got != tt.want {
				t.Errorf("IsTree = %v, quero %v", got, tt.want)
			}
		})
	}
}

func TestSpanningTrees(t *testing.T) {
	g := mustBuild(t, 6, sampleEdges)

	tests := []struct {
		name       string
		build      func(graph.Graph, int) []graph.Edge
		want       []graph.Edge
		wantHeight int
	}{
		{"BFS", SpanningTreeBFS, []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 3, V: 5}}, 3},
		{"DFS", SpanningTreeDFS, []graph.Edge{{U: 0, V: 1}, {U: 1, V: 3}, {U: 3, V: 5}, {U: 5, V: 4}, {U: 4, V: 2}}, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.build(g, 0)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("arestas = %v, quero %v", got, tt.want)
			}
			r := Check(g, 0, got)
			if !r.Valid(6) {
				t.Errorf("Check = %+v, quero árvore geradora válida", r)
			}
			if r.Height != tt.wantHeight {
				t.Errorf("altura = %d, quero %d", r.Height, tt.wantHeight)
			}
		})
	}
}

func TestCheckRejectsInvalidTrees(t *testing.T) {
	g := mustBuild(t, 6, sampleEdges)

	tests := []struct {
		name string
		t    []graph.Edge
	}{
		{"aresta fora do grafo", []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 0, V: 5}}},
		{"ciclo e vértice isolado", []graph.Edge{{U: 0, V: 1}, {U: 1, V: 4}, {U: 2, V: 4}, {U: 0, V: 2}, {U: 3, V: 5}}},
		{"faltando aresta", []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if r := Check(g, 0, tt.t); r.Valid(6) {
				t.Errorf("Check = %+v, quero inválida", r)
			}
		})
	}
	if r := Check(g, 0, tests[1].t); r.Acyclic {
		t.Error("ciclo 0-1-4-2-0 deveria tornar Acyclic false")
	}
}

func TestArborescenceBFS(t *testing.T) {
	arcs := []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 2, V: 4}, {U: 3, V: 5}, {U: 4, V: 5}, {U: 5, V: 1}}
	d, err := graph.BuildDigraph(6, arcs)
	if err != nil {
		t.Fatalf("BuildDigraph: %v", err)
	}

	got, err := ArborescenceBFS(d, 0)
	if err != nil {
		t.Fatalf("raiz 0: %v", err)
	}
	want := []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 2, V: 4}, {U: 3, V: 5}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("raiz 0: arcos = %v, quero %v", got, want)
	}

	_, err = ArborescenceBFS(d, 3)
	var ue *UnreachableError
	if !errors.As(err, &ue) {
		t.Fatalf("raiz 3: erro = %v, quero *UnreachableError", err)
	}
	if want := []int{0, 2, 4}; !reflect.DeepEqual(ue.Unreachable, want) {
		t.Errorf("raiz 3: inalcançáveis = %v, quero %v", ue.Unreachable, want)
	}
}
