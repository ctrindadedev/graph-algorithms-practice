package main

import (
	"fmt"
	"strings"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

func title(s string) {
	fmt.Println(s)
	fmt.Println(strings.Repeat("-", len([]rune(s))))
}

func labels(vs []int) string {
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = label(v)
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

func path(vs []int) string {
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = label(v)
	}
	return strings.Join(parts, " → ")
}

// printAdjacency mostra a lista com letras: só vizinhos no digrafo e pares
// (vizinho, peso) nos ponderados.
func printAdjacency(g graph.Neighborhood) {
	for u := range g.Order() {
		w, weighted := g.(*graph.Weighted)
		if !weighted {
			fmt.Printf("  %s → %s\n", label(u), labels(g.Neighbors(u)))
			continue
		}
		var parts []string
		for _, a := range w.Arcs(u) {
			parts = append(parts, fmt.Sprintf("(%s, %d)", label(a.To), a.Weight))
		}
		fmt.Printf("  %s → [%s]\n", label(u), strings.Join(parts, ", "))
	}
}

// newStepPrinter reproduz os "Escrever" dos pseudocódigos: vértice processado
// e fila na BFS, entrada e retorno indentados pela profundidade na DFS.
func newStepPrinter() traversal.StepFunc {
	depth := 0
	return func(s traversal.Step) {
		switch s.Kind {
		case traversal.Dequeue:
			fmt.Printf("  Processando: %s  Enfileirados: %-9s Fila: %s\n", label(s.Vertex), labels(s.Added), labels(s.Frontier))
		case traversal.Pop:
			if s.Skipped {
				fmt.Printf("  Desempilhado: %s (já visitado, descartado)  Pilha: %s\n", label(s.Vertex), labels(s.Frontier))
				return
			}
			fmt.Printf("  Desempilhado: %s  Empilhados: %-9s Pilha (topo primeiro): %s\n", label(s.Vertex), labels(s.Added), labels(s.Frontier))
		case traversal.Enter:
			fmt.Printf("  %sEntrando no vértice: %s\n", strings.Repeat("  ", depth), label(s.Vertex))
			depth++
		case traversal.Leave:
			depth--
			fmt.Printf("  %sRetornando do vértice: %s\n", strings.Repeat("  ", depth), label(s.Vertex))
		}
	}
}
