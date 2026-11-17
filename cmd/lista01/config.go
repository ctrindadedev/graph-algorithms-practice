package main

import (
	"bufio"
	"fmt"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/cli"
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
// cada aresta e o vértice inicial (Questão 13).
func readConfig(in *bufio.Reader) (config, error) {
	n, err := cli.ReadInt(in, "Número de vértices: ")
	if err != nil {
		return config{}, err
	}
	if n <= 0 {
		return config{}, fmt.Errorf("o número de vértices deve ser positivo (recebi %d)", n)
	}

	m, err := cli.ReadInt(in, "Número de arestas: ")
	if err != nil {
		return config{}, err
	}
	if m < 0 {
		return config{}, fmt.Errorf("o número de arestas não pode ser negativo (recebi %d)", m)
	}

	edges := make([]graph.Edge, 0, m)
	for k := range m {
		u, err := cli.ReadInt(in, fmt.Sprintf("Aresta %d (u v): ", k+1))
		if err != nil {
			return config{}, err
		}
		v, err := cli.ReadInt(in, "")
		if err != nil {
			return config{}, err
		}
		if u < 0 || u >= n || v < 0 || v >= n || u == v {
			return config{}, fmt.Errorf("aresta (%d, %d) inválida: use vértices distintos entre 0 e %d", u, v, n-1)
		}
		edges = append(edges, graph.Edge{U: u, V: v})
	}

	start, err := cli.ReadInt(in, "Vértice inicial: ")
	if err != nil {
		return config{}, err
	}
	if start < 0 || start >= n {
		return config{}, fmt.Errorf("vértice inicial %d fora do intervalo [0, %d]", start, n-1)
	}
	fmt.Println()

	return config{n: n, edges: edges, start: start}, nil
}
