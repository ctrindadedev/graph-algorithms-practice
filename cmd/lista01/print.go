package main

import (
	"fmt"
	"strings"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

func printGraph(title string, g graph.Graph) {
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", len([]rune(title))))
	fmt.Print(g)
}

// newStepPrinter exibe cada passo: estado da fila/pilha em Dequeue/Pop (Questões 7 e 8)
// e entrada/retorno indentados pela profundidade em Enter/Leave (Questão 4).
func newStepPrinter() traversal.StepFunc {
	depth := 0
	return func(s traversal.Step) {
		switch s.Kind {
		case traversal.Dequeue:
			fmt.Printf("Retirado da fila: %d\n", s.Vertex)
			field("Vizinhos encontrados", s.Neighbors)
			field("Novos na fila", s.Added)
			field("Fila", s.Frontier)
			field("Visitados", s.Visited)
		case traversal.Pop:
			if s.Skipped {
				fmt.Printf("Retirado do topo: %d (já visitado, ignorado)\n", s.Vertex)
				field("Pilha", s.Frontier)
				return
			}
			fmt.Printf("Retirado do topo: %d\n", s.Vertex)
			field("Vizinhos encontrados", s.Neighbors)
			field("Empilhados", s.Added)
			field("Pilha (topo primeiro)", s.Frontier)
			field("Visitados", s.Visited)
		case traversal.Enter:
			fmt.Printf("%sEntrando no vértice %d\n", strings.Repeat("  ", depth), s.Vertex)
			depth++
		case traversal.Leave:
			depth--
			fmt.Printf("%sRetornando do vértice %d\n", strings.Repeat("  ", depth), s.Vertex)
		}
	}
}

func field(label string, xs []int) {
	fmt.Printf("  %-22s %s\n", label+":", graph.FormatInts(xs))
}

// printOrder exibe a sequência final, ex.: "Percurso: 0 → 1 → 3 → 5 → 4 → 2".
func printOrder(order []int) {
	parts := make([]string, len(order))
	for i, v := range order {
		parts[i] = fmt.Sprint(v)
	}
	fmt.Println("Percurso:", strings.Join(parts, " → "))
}
