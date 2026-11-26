// Comando dijkstra: estudo dirigido do algoritmo de Dijkstra (Unidade II).
// Executa todas as partes em sequência: lista de adjacência, rastreamento,
// reconstrução de caminhos, testes da Questão 6 e o desafio (digrafo + heap).
package main

import (
	"fmt"
	"os"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	r "github.com/ctrindadedev/graph-algorithms-practice/internal/report"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/shortest"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

// Os vértices A..F viram os índices 0..5: a ordem crescente dos índices é a
// ordem alfabética, então o desempate "menor letra" é o desempate "menor índice".
const (
	A = iota
	B
	C
	D
	E
	F
	G // vértice isolado do teste 4
)

// arestas do grafo-base (Figura 1), não direcionado.
var arestas = []graph.WeightedEdge{
	{U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: B, V: C, Weight: 1},
	{U: B, V: D, Weight: 5}, {U: C, V: D, Weight: 8}, {U: C, V: E, Weight: 10},
	{U: D, V: E, Weight: 2}, {U: D, V: F, Weight: 6}, {U: E, V: F, Weight: 3},
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func run() error {
	// PARTE I – criar_grafo: lista de adjacência não direcionada.
	grafo, err := graph.BuildWeighted(F+1, false, arestas)
	if err != nil {
		return err
	}
	r.Title("PARTE I – Lista de adjacência (não direcionada)")
	printAdjacency(grafo)

	// PARTES II e III – Dijkstra com seleção linear e rastreamento.
	r.Title("PARTES II e III – Dijkstra a partir de A (seleção linear)")
	res, err := shortest.DijkstraTrace(grafo, A, printStep)
	if err != nil {
		return err
	}
	printFinal(res)

	// PARTE IV – reconstruir_caminho.
	r.Title("PARTE IV – Reconstrução dos caminhos mínimos")
	for _, dest := range []int{F, B, D, E} {
		printRoute(res, dest)
	}

	// PARTE VI – testes da Questão 6.
	r.Title("PARTE VI – Testes")
	if err := runTests(res); err != nil {
		return err
	}
	if err := compareBFS(grafo, res); err != nil {
		return err
	}

	// PARTE VII – desafio: digrafo e heap.
	return challenge()
}

func runTests(base shortest.Result) error {
	r.Subtitle("Teste 1 – Original; A → F")
	printRoute(base, F)

	r.Subtitle("Teste 2 – Original; A → A")
	printRoute(base, A)

	r.Subtitle("Teste 3 – Original; A → D")
	printRoute(base, D)
	fmt.Printf("  Predecessores intermediários: anterior[D] = %s, anterior[B] = %s, anterior[C] = %s\n",
		r.Pred(base.Pred[D]), r.Pred(base.Pred[B]), r.Pred(base.Pred[C]))

	// Teste 4: mesma lista de arestas, mas com 7 vértices; G não recebe arestas.
	r.Subtitle("Teste 4 – Vértice G isolado; A → G")
	g4, err := graph.BuildWeighted(G+1, false, arestas)
	if err != nil {
		return err
	}
	res4, err := shortest.DijkstraTrace(g4, A, printOrder())
	if err != nil {
		return err
	}
	fmt.Println()
	printRoute(res4, G)

	// Teste 5: cópia das arestas com A—B passando de 4 para 1.
	r.Subtitle("Teste 5 – Peso de A—B trocado de 4 para 1; A → F")
	copia := append([]graph.WeightedEdge(nil), arestas...)
	copia[0].Weight = 1
	g5, err := graph.BuildWeighted(F+1, false, copia)
	if err != nil {
		return err
	}
	res5, err := shortest.DijkstraTrace(g5, A, printStep)
	if err != nil {
		return err
	}
	printRoute(res5, F)
	return nil
}

// compareBFS mostra a rota da árvore BFS até F, que minimiza o número de
// arestas e ignora os pesos.
func compareBFS(g *graph.Weighted, res shortest.Result) error {
	r.Subtitle("BFS × Dijkstra até F")
	pred := make([]int, g.Order())
	for v := range pred {
		pred[v] = shortest.NoPred
	}
	traversal.BFS(g, A, func(s traversal.Step) {
		for _, v := range s.Added {
			pred[v] = s.Vertex // o pai de v na árvore BFS é quem o enfileirou
		}
	})
	bfs := shortest.Result{Source: A, Dist: make([]int, g.Order()), Pred: pred}
	path := bfs.PathTo(F)
	cost := 0
	for i := 1; i < len(path); i++ {
		w, _ := g.Weight(path[i-1], path[i])
		cost += w
	}
	fmt.Printf("  BFS:      %-26s %d aresta(s), custo %d\n", r.Path(path), len(path)-1, cost)
	dp := res.PathTo(F)
	fmt.Printf("  Dijkstra: %-26s %d aresta(s), custo %d\n", r.Path(dp), len(dp)-1, res.Dist[F])
	return nil
}

func challenge() error {
	r.Title("PARTE VII – Desafio: grafo DIRECIONADO e versão com heapq")
	// Só a construção muda: cada aresta vira o arco u → v, sem retorno.
	dg, err := graph.BuildWeighted(F+1, true, arestas)
	if err != nil {
		return err
	}
	fmt.Println("Lista de adjacência (direcionada):")
	printAdjacency(dg)

	r.Subtitle("Seleção linear")
	lin, err := shortest.DijkstraTrace(dg, A, printStep)
	if err != nil {
		return err
	}
	printRoute(lin, F)

	r.Subtitle("Fila de prioridade (heap)")
	hp, err := shortest.DijkstraHeap(dg, A, printStep)
	if err != nil {
		return err
	}
	printRoute(hp, F)

	r.Subtitle("Heap no grafo não direcionado original")
	ug, _ := graph.BuildWeighted(F+1, false, arestas)
	uh, err := shortest.DijkstraHeap(ug, A, printOrder())
	if err != nil {
		return err
	}
	fmt.Println()
	printFinal(uh)
	return nil
}

// printAdjacency é o PASSO 3 do CriarGrafoPonderado: v e seus pares (vizinho, peso).
func printAdjacency(g *graph.Weighted) {
	for u := range g.Order() {
		fmt.Printf("  %s → [", r.Label(u))
		for i, a := range g.Arcs(u) {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("(%s, %d)", r.Label(a.To), a.Weight)
		}
		fmt.Println("]")
	}
}

// printStep reproduz "Escrever Processando, u, distancia[u]" e o estado dos
// dicionários ao fim da iteração.
func printStep(s shortest.Step) {
	fmt.Printf("  Processando %s (distância %s)\n", r.Label(s.Vertex), r.Dist(s.Dist[s.Vertex]))
	fmt.Printf("    distancia: %s\n", vector(s.Dist, r.Dist))
	fmt.Printf("    anterior:  %s\n", vector(s.Pred, r.Pred))
}

// printOrder só lista os vértices processados, em ordem.
func printOrder() shortest.StepFunc {
	first := true
	return func(s shortest.Step) {
		if first {
			fmt.Print("  Processados:")
			first = false
		}
		fmt.Printf(" %s", r.Label(s.Vertex))
	}
}

func vector(xs []int, cell func(int) string) string {
	s := "{"
	for v, x := range xs {
		if v > 0 {
			s += ", "
		}
		s += r.Label(v) + ": " + cell(x)
	}
	return s + "}"
}

func printFinal(res shortest.Result) {
	fmt.Println("\n  Resultado final:")
	r.Header("", 12, len(res.Dist), 4)
	r.Row("distancia", 12, r.Dists(res.Dist), 4)
	r.Row("anterior", 12, r.Preds(res.Pred), 4)
}

func printRoute(res shortest.Result, dest int) {
	p := res.PathTo(dest)
	if p == nil {
		fmt.Printf("  %s → %s: inalcançável (distância ∞, caminho [])\n", r.Label(res.Source), r.Label(dest))
		return
	}
	fmt.Printf("  %s → %s: %-26s custo %d, %d aresta(s)\n", r.Label(res.Source), r.Label(dest), r.Path(p), res.Dist[dest], len(p)-1)
}
