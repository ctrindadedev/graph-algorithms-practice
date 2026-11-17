package main

import "github.com/ctrindadedev/graph-algorithms-practice/internal/graph"

// Os vértices A..E viram os índices 0..4; assim a ordem crescente dos
// índices é a mesma ordem alfabética pedida no enunciado.
const (
	A = iota
	B
	C
	D
	E
	order
)

func label(v int) string {
	return string(rune('A' + v))
}

// Etapa 1: arcos sem peso.
var directedArcs = []graph.Edge{{U: A, V: B}, {U: A, V: C}, {U: B, V: D}, {U: C, V: D}, {U: D, V: E}}

// Etapa 2: arestas com peso, válidas nos dois sentidos.
var weightedEdges = []graph.WeightedEdge{
	{U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: B, V: C, Weight: 1}, {U: B, V: D, Weight: 5},
	{U: C, V: D, Weight: 8}, {U: C, V: E, Weight: 10}, {U: D, V: E, Weight: 2},
}

// Etapa 3: mesmos pesos, mas orientados; note C → B no lugar de B — C.
var weightedArcs = []graph.WeightedEdge{
	{U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: C, V: B, Weight: 1}, {U: B, V: D, Weight: 5},
	{U: C, V: D, Weight: 8}, {U: C, V: E, Weight: 10}, {U: D, V: E, Weight: 2},
}

// stage descreve um dos três grafos do estudo dirigido.
type stage struct {
	title string
	build func() (graph.Neighborhood, error)
}

var stages = []stage{
	{"Etapa 1 – Grafo direcionado", func() (graph.Neighborhood, error) {
		return graph.BuildDigraph(order, directedArcs)
	}},
	{"Etapa 2 – Grafo ponderado", func() (graph.Neighborhood, error) {
		return graph.BuildWeighted(order, false, weightedEdges)
	}},
	{"Etapa 3 – Grafo direcionado e ponderado", func() (graph.Neighborhood, error) {
		return graph.BuildWeighted(order, true, weightedArcs)
	}},
}
