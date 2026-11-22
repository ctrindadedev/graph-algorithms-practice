// Comando lista02: menu com os algoritmos da Aula 9 (árvores e árvores geradoras).
// Flag: -start define a raiz das árvores geradoras.
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
	"github.com/ctrindadedev/graph-algorithms-practice/internal/tree"
)

func main() {
	start := flag.Int("start", 0, "vértice inicial das árvores geradoras")
	flag.Parse()

	if *start < 0 || *start >= sample.n {
		fmt.Fprintf(os.Stderr, "erro: vértice inicial %d fora do intervalo [0, %d]\n", *start, sample.n-1)
		os.Exit(1)
	}

	if err := runMenu(bufio.NewReader(os.Stdin), *start); err != nil {
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
	case 1:
		return runIsTree()
	case 2, 3:
		g, err := build(sample)
		if err != nil {
			return err
		}
		if opt == 2 {
			showSpanningTree("Árvore geradora por BFS", g, start, tree.SpanningTreeBFS(g, start))
		} else {
			showSpanningTree("Árvore geradora por DFS", g, start, tree.SpanningTreeDFS(g, start))
		}
	case 4:
		g, err := build(sample)
		if err != nil {
			return err
		}
		compare(g, start)
	case 5:
		return runArborescence()
	case 6:
		g, err := build(sample)
		if err != nil {
			return err
		}
		fmt.Println(sample.name)
		fmt.Print(g)
	default:
		fmt.Printf("Opção %d inválida.\n", opt)
	}
	return nil
}

func build(ex example) (graph.Graph, error) {
	return graph.Build(func(n int) graph.Graph { return graph.NewAdjacencyList(n) }, ex.n, ex.edges)
}

func runIsTree() error {
	for i, ex := range treeExamples {
		g, err := build(ex)
		if err != nil {
			return err
		}
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(ex.name)
		fmt.Printf("  E = %s\n", formatEdges(ex.edges))
		fmt.Printf("  |V| = %d, |E| = %d\n", g.Order(), graph.Size(g))
		fmt.Printf("  É árvore? %s\n", yesNo(tree.IsTree(g)))
	}
	return nil
}

func showSpanningTree(title string, g graph.Graph, start int, t []graph.Edge) tree.Report {
	r := tree.Check(g, start, t)
	n := g.Order()

	fmt.Printf("%s (raiz %d)\n", title, start)
	fmt.Println(strings.Repeat("-", len([]rune(title))+len(fmt.Sprintf(" (raiz %d)", start))))
	fmt.Println("Arestas na ordem de descoberta:")
	for i, e := range t {
		fmt.Printf("  %d. %s\n", i+1, e)
	}
	fmt.Printf("Vértices: %d  Arestas: %d\n", r.Vertices, r.Edges)
	fmt.Printf("  Todos os vértices presentes:    %s\n", yesNo(r.Vertices == n))
	fmt.Printf("  Exatamente V - 1 = %d arestas:   %s\n", n-1, yesNo(r.HasVMinusOne))
	fmt.Printf("  Arestas pertencem ao grafo:     %s\n", yesNo(r.EdgesInGraph))
	fmt.Printf("  Sem ciclos:                     %s\n", yesNo(r.Acyclic))
	fmt.Printf("  Árvore geradora válida:         %s\n", yesNo(r.Valid(n)))
	fmt.Printf("Altura: %d  Largura máxima: %d\n", r.Height, r.MaxWidth)
	return r
}

func compare(g graph.Graph, start int) {
	bfsEdges, dfsEdges := tree.SpanningTreeBFS(g, start), tree.SpanningTreeDFS(g, start)
	bfs := showSpanningTree("Árvore geradora por BFS", g, start, bfsEdges)
	fmt.Println()
	dfs := showSpanningTree("Árvore geradora por DFS", g, start, dfsEdges)
	fmt.Println()

	fmt.Println("Comparação BFS × DFS")
	fmt.Println("--------------------")
	fmt.Printf("%-24s %-22s %s\n", "Aspecto", "BFS", "DFS")
	fmt.Printf("%-24s %-22s %s\n", "Estrutura auxiliar", "Fila (FIFO)", "Pilha (LIFO)")
	fmt.Printf("%-24s %-22d %d\n", "Altura", bfs.Height, dfs.Height)
	fmt.Printf("%-24s %-22d %d\n", "Largura máxima", bfs.MaxWidth, dfs.MaxWidth)
	fmt.Printf("%-24s %-22s %s\n", "Arestas exclusivas", formatEdges(only(bfsEdges, dfsEdges)), formatEdges(only(dfsEdges, bfsEdges)))
	fmt.Printf("Arestas em comum: %s\n", formatEdges(common(bfsEdges, dfsEdges)))

	switch {
	case dfs.Height > bfs.Height:
		fmt.Println("A árvore da DFS ficou mais profunda e a da BFS mais larga.")
	case bfs.Height > dfs.Height:
		fmt.Println("A árvore da BFS ficou mais profunda e a da DFS mais larga.")
	default:
		fmt.Println("As duas árvores têm a mesma altura.")
	}
}

func runArborescence() error {
	d, err := graph.BuildDigraph(digraph.n, digraph.edges)
	if err != nil {
		return err
	}
	fmt.Println(digraph.name, "(arcos que saem de cada vértice)")
	fmt.Print(d)

	for _, root := range []int{0, 3} {
		fmt.Printf("\nArborescência por BFS com raiz %d\n", root)
		a, err := tree.ArborescenceBFS(d, root)
		if err != nil {
			fmt.Println(" ", err)
			continue
		}
		for i, arc := range a {
			fmt.Printf("  %d. %d → %d\n", i+1, arc.U, arc.V)
		}
		fmt.Printf("  Vértices: %d  Arcos: %d\n", d.Order(), len(a))
	}
	return nil
}

func normalize(e graph.Edge) graph.Edge {
	if e.U > e.V {
		return graph.Edge{U: e.V, V: e.U}
	}
	return e
}

func contains(edges []graph.Edge, e graph.Edge) bool {
	for _, x := range edges {
		if normalize(x) == normalize(e) {
			return true
		}
	}
	return false
}

func only(a, b []graph.Edge) []graph.Edge {
	var out []graph.Edge
	for _, e := range a {
		if !contains(b, e) {
			out = append(out, e)
		}
	}
	return out
}

func common(a, b []graph.Edge) []graph.Edge {
	var out []graph.Edge
	for _, e := range a {
		if contains(b, e) {
			out = append(out, e)
		}
	}
	return out
}

func formatEdges(edges []graph.Edge) string {
	parts := make([]string, len(edges))
	for i, e := range edges {
		parts[i] = e.String()
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func yesNo(ok bool) string {
	if ok {
		return "sim"
	}
	return "não"
}

const menu = `=================================
   ÁRVORES E ÁRVORES GERADORAS
=================================
1 - Verificar se é Árvore
2 - Árvore Geradora por BFS
3 - Árvore Geradora por DFS
4 - Comparar BFS × DFS
5 - Arborescência por BFS
6 - Mostrar Grafo
0 - Sair
=================================
`
