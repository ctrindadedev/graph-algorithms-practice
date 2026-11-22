// Comando matriz_lista_bfs_dfs: menu com os algoritmos da Lista 01 (Questão 9).
// Flags: -extra adiciona a aresta (2,5) (Questão 10); -input lê o grafo do usuário (Questão 13).
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
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

func main() {
	extra := flag.Bool("extra", false, "acrescenta a aresta (2,5) ao grafo (Questão 10)")
	input := flag.Bool("input", false, "lê o grafo e o vértice inicial da entrada padrão (Questão 13)")
	flag.Parse()

	in := bufio.NewReader(os.Stdin)

	cfg := defaultConfig(*extra)
	if *input {
		var err error
		if cfg, err = readConfig(in); err != nil {
			fmt.Fprintln(os.Stderr, "erro:", err)
			os.Exit(1)
		}
	}

	if err := runMenu(in, cfg); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func runMenu(in *bufio.Reader, cfg config) error {
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
		if err := runOption(opt, cfg); err != nil {
			return err
		}
		fmt.Println()
	}
}

func runOption(opt int, cfg config) error {
	newMatrix := func(n int) graph.Graph { return graph.NewAdjacencyMatrix(n) }
	newList := func(n int) graph.Graph { return graph.NewAdjacencyList(n) }

	switch opt {
	case 1:
		g, err := graph.Build(newMatrix, cfg.n, cfg.edges)
		if err != nil {
			return err
		}
		printGraph("Matriz de Adjacência", g)
	case 2:
		g, err := graph.Build(newList, cfg.n, cfg.edges)
		if err != nil {
			return err
		}
		printGraph("Lista de Adjacência", g)
	case 3, 4, 5:
		g, err := graph.Build(newList, cfg.n, cfg.edges)
		if err != nil {
			return err
		}
		runTraversal(opt, g, cfg.start)
	case 6:
		g, err := graph.Build(newList, cfg.n, cfg.edges)
		if err != nil {
			return err
		}
		showGraph(g, cfg.start)
	default:
		fmt.Printf("Opção %d inválida.\n", opt)
	}
	return nil
}

func runTraversal(opt int, g graph.Graph, start int) {
	var order []int
	switch opt {
	case 3:
		fmt.Printf("BFS a partir do vértice %d\n\n", start)
		order = traversal.BFS(g, start, newStepPrinter())
	case 4:
		fmt.Printf("DFS recursiva a partir do vértice %d\n\n", start)
		order = traversal.DFSRecursive(g, start, newStepPrinter())
	case 5:
		fmt.Printf("DFS com pilha a partir do vértice %d\n", start)
		fmt.Printf("Estado inicial: Pilha: [%d]  Visitados: []\n\n", start)
		order = traversal.DFSStack(g, start, newStepPrinter())
	}
	fmt.Println()
	printOrder(order)
}

func showGraph(g graph.Graph, start int) {
	vertices := make([]string, g.Order())
	for i := range vertices {
		vertices[i] = fmt.Sprint(i)
	}
	var edges []string
	for _, e := range graph.Edges(g) {
		edges = append(edges, e.String())
	}
	fmt.Printf("V = {%s}\n", strings.Join(vertices, ", "))
	fmt.Printf("E = {%s}\n", strings.Join(edges, ", "))
	fmt.Printf("|V| = %d, |E| = %d, vértice inicial = %d\n", g.Order(), graph.Size(g), start)
}

const menu = `=================================
        ALGORITMOS DE GRAFOS
=================================
1 - Criar Matriz de Adjacência
2 - Criar Lista de Adjacência
3 - Executar BFS
4 - Executar DFS Recursiva
5 - Executar DFS com Pilha
6 - Mostrar Grafo
0 - Sair
=================================
`
