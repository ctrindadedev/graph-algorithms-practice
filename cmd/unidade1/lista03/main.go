// Comando lista03: menu com BFS e DFS em grafos direcionados, ponderados e
// direcionados ponderados, mais o desafio de Dijkstra. Flag: -start (A..E).
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/cli"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/shortest"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

func main() {
	startFlag := flag.String("start", "A", "vértice de origem (A..E)")
	flag.Parse()

	s := strings.ToUpper(*startFlag)
	if len(s) != 1 || s[0] < 'A' || int(s[0]-'A') >= order {
		fmt.Fprintf(os.Stderr, "erro: origem %q inválida, use uma letra de A a %s\n", *startFlag, label(order-1))
		os.Exit(1)
	}

	if err := runMenu(bufio.NewReader(os.Stdin), int(s[0]-'A')); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func runMenu(in *bufio.Reader, start int) error {
	for {
		fmt.Print(menu)
		opt, err := cli.ReadInt(in, "Opção: ")
		if errors.Is(err, io.EOF) {
			return nil
		}
		fmt.Println()
		if err != nil {
			fmt.Println("Opção inválida:", err)
			fmt.Println()
			continue
		}
		if opt == 0 {
			return nil
		}
		if err := runOption(opt, start); err != nil {
			return err
		}
		fmt.Println()
	}
}

func runOption(opt, start int) error {
	switch opt {
	case 1, 2, 3:
		return runStage(stages[opt-1], start)
	case 4:
		return compareAll(start)
	case 5:
		return runDijkstra(start)
	case 6:
		return runIterativeDFS(start)
	default:
		fmt.Printf("Opção %d inválida.\n", opt)
	}
	return nil
}

func runStage(st stage, start int) error {
	g, err := st.build()
	if err != nil {
		return err
	}
	title(st.title)
	fmt.Println("Lista de adjacência:")
	printAdjacency(g)

	fmt.Printf("\nBFS a partir de %s:\n", label(start))
	bfs := traversal.BFS(g, start, newStepPrinter())
	fmt.Println("  Ordem BFS:", path(bfs))

	fmt.Printf("\nDFS recursiva a partir de %s:\n", label(start))
	dfs := traversal.DFSRecursive(g, start, newStepPrinter())
	fmt.Println("  Ordem DFS:", path(dfs))

	reportUnreachable(g, bfs)
	return nil
}

// reportUnreachable avisa quais vértices ficaram fora do percurso (Pergunta 6).
func reportUnreachable(g graph.Neighborhood, visited []int) {
	seen := make([]bool, g.Order())
	for _, v := range visited {
		seen[v] = true
	}
	var missing []int
	for v, ok := range seen {
		if !ok {
			missing = append(missing, v)
		}
	}
	if len(missing) > 0 {
		fmt.Println("  Inalcançáveis a partir da origem:", labels(missing))
	}
}

func compareAll(start int) error {
	title(fmt.Sprintf("Comparação dos percursos a partir de %s", label(start)))
	fmt.Printf("%-26s %-22s %s\n", "Tipo de grafo", "Ordem BFS", "Ordem DFS")
	names := []string{"Direcionado", "Ponderado", "Direcionado e ponderado"}
	for i, st := range stages {
		g, err := st.build()
		if err != nil {
			return err
		}
		fmt.Printf("%-26s %-22s %s\n", names[i], path(traversal.BFS(g, start, nil)), path(traversal.DFSRecursive(g, start, nil)))
	}
	return nil
}

func runDijkstra(start int) error {
	g, err := graph.BuildWeighted(order, true, weightedArcs)
	if err != nil {
		return err
	}
	r, err := shortest.Dijkstra(g, start)
	if err != nil {
		return err
	}

	title(fmt.Sprintf("Dijkstra no grafo direcionado e ponderado (origem %s)", label(start)))
	fmt.Printf("%-8s %-10s %s\n", "Vértice", "dist", "pred")
	for v := range order {
		dist, pred := "∞", "–"
		if r.Dist[v] != shortest.Inf {
			dist = fmt.Sprint(r.Dist[v])
		}
		if r.Pred[v] != shortest.NoPred {
			pred = label(r.Pred[v])
		}
		fmt.Printf("%-8s %-10s %s\n", label(v), dist, pred)
	}

	target := E
	fmt.Printf("\nMenor custo de %s até %s:\n", label(start), label(target))
	if p := r.PathTo(target); p != nil {
		fmt.Printf("  Dijkstra: %s  (%d arco(s), custo %d)\n", path(p), len(p)-1, r.Dist[target])
	} else {
		fmt.Printf("  Dijkstra: %s é inalcançável\n", label(target))
	}

	// Caminho da árvore de busca BFS: o pai de v é o vértice que o enfileirou.
	pred := make([]int, order)
	for v := range pred {
		pred[v] = shortest.NoPred
	}
	traversal.BFS(g, start, func(s traversal.Step) {
		for _, v := range s.Added {
			pred[v] = s.Vertex
		}
	})
	if target != start && pred[target] == shortest.NoPred {
		fmt.Printf("  BFS:      %s é inalcançável\n", label(target))
		return nil
	}
	var bfsPath []int
	cost := 0
	for v := target; v != shortest.NoPred; v = pred[v] {
		bfsPath = append([]int{v}, bfsPath...)
		if pred[v] != shortest.NoPred {
			w, _ := g.Weight(pred[v], v)
			cost += w
		}
	}
	fmt.Printf("  BFS:      %s  (%d arco(s), custo %d)\n", path(bfsPath), len(bfsPath)-1, cost)
	return nil
}

func runIterativeDFS(start int) error {
	title(fmt.Sprintf("Extensão – DFS iterativa com pilha (origem %s)", label(start)))
	fmt.Println("Os vizinhos são empilhados em ordem alfabética INVERSA para o menor ficar no topo.")
	for _, st := range stages {
		g, err := st.build()
		if err != nil {
			return err
		}
		fmt.Printf("\n%s\n", st.title)
		iter := traversal.DFSStack(g, start, newStepPrinter())
		rec := traversal.DFSRecursive(g, start, nil)
		fmt.Printf("  Ordem iterativa: %s  (recursiva: %s)\n", path(iter), path(rec))
	}
	return nil
}

const menu = `=========================================
  BFS E DFS: DIRECIONADO E PONDERADO
=========================================
1 - Etapa 1: Grafo direcionado
2 - Etapa 2: Grafo ponderado
3 - Etapa 3: Grafo direcionado e ponderado
4 - Tabela comparativa BFS × DFS
5 - Desafio: Dijkstra × BFS
6 - Extensão: DFS iterativa com pilha
0 - Sair
=========================================
`
