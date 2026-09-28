package main

import (
	"bufio"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

type config struct {
	n     int
	edges []graph.Edge
	start int
}

var sampleEdges = []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 2, V: 4}, {U: 3, V: 5}, {U: 4, V: 5}}

func defaultConfig(extra bool) config {
	edges := append([]graph.Edge{}, sampleEdges...)
	if extra {
		edges = append(edges, graph.Edge{U: 2, V: 5})
	}
	return config{n: 6, edges: edges, start: 0}
}

// readConfig lê do usuário o número de vértices, o número de arestas,
// cada aresta e o vértice inicial
func readConfig(in *bufio.Reader) (config, error) {
	// TODO: ler n, m, as m arestas (u v) e o vértice inicial,
	// validando os valores (dica: fmt.Fscan(in, &x)).
	return config{}, nil
}
