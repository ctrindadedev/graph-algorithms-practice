package main

import "github.com/ctrindadedev/graph-algorithms-practice/internal/graph"

type example struct {
	name  string
	n     int
	edges []graph.Edge
}

// sample é o grafo da Lista 01, usado para as árvores geradoras.
var sample = example{
	name:  "Grafo da Lista 01",
	n:     6,
	edges: []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 2, V: 4}, {U: 3, V: 5}, {U: 4, V: 5}},
}

// treeExamples cobrem os casos verdadeiro e falso de EhArvore.
var treeExamples = []example{
	{
		name:  "Árvore (5 arestas, conexo)",
		n:     6,
		edges: []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 2, V: 5}},
	},
	sample,
	{
		name:  "V - 1 arestas, mas com ciclo e desconexo",
		n:     6,
		edges: []graph.Edge{{U: 0, V: 1}, {U: 1, V: 2}, {U: 2, V: 0}, {U: 3, V: 4}, {U: 4, V: 5}},
	},
}

// digraph é o grafo direcionado usado na arborescência. A partir de 0 todos
// os vértices são alcançáveis; a partir de 3, os vértices 0, 2 e 4 não são.
var digraph = example{
	name:  "Digrafo",
	n:     6,
	edges: []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 2, V: 4}, {U: 3, V: 5}, {U: 4, V: 5}, {U: 5, V: 1}},
}
