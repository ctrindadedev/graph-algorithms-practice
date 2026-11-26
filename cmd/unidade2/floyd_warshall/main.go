// Comando floyd_warshall: estudo dirigido do algoritmo de Floyd-Warshall
// (Unidade II). Executa todas as partes em sequência: matriz inicial, uma
// matriz por intermediário, matriz final, rotas, testes e o desafio
// (rotas a partir de A comparadas com Bellman-Ford).
package main

import (
	"fmt"
	"os"
	"slices"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	r "github.com/ctrindadedev/graph-algorithms-practice/internal/report"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/shortest"
)

// A..F = 0..5: linha = origem, coluna = destino, na ordem alfabética.
const (
	A = iota
	B
	C
	D
	E
	F
	G // vértice isolado do teste 4
)

var arestas = []graph.WeightedEdge{
	{U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: C, V: B, Weight: -1},
	{U: B, V: D, Weight: 2}, {U: C, V: D, Weight: 5}, {U: D, V: E, Weight: -2},
	{U: E, V: F, Weight: 3}, {U: B, V: F, Weight: 10}, {U: C, V: E, Weight: 7},
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func run() error {
	// PARTE I – criar_matrizes.
	dist, prox, err := shortest.NewMatrices(F+1, arestas)
	if err != nil {
		return err
	}
	r.Title("PARTE I – Matriz inicial D(0) (só arestas diretas)")
	r.PrintMatrix(dist, r.Dist)
	fmt.Println("\n  Matriz prox inicial:")
	r.PrintMatrix(prox, r.Pred)

	// PARTES II e III – uma matriz por intermediário k = A, B, ..., F.
	r.Title("PARTES II e III – Floyd-Warshall (k → i → j)")
	res := shortest.FloydWarshall(dist, prox, printStage)
	r.Subtitle("Matriz final de menores custos")
	r.PrintMatrix(res.Dist, r.Dist)
	fmt.Println("\n  Matriz prox final:")
	r.PrintMatrix(res.Next, r.Pred)
	printCycle(res)

	// PARTE IV – reconstruir_caminho.
	r.Title("PARTE IV – Reconstrução dos caminhos mínimos")
	for _, q := range [][2]int{{A, F}, {A, D}, {C, F}, {F, A}} {
		printRoute(res, q[0], q[1])
	}

	r.Title("PARTE VI – Testes")
	if err := runTests(res); err != nil {
		return err
	}
	return challenge(res)
}

func runTests(base shortest.APSP) error {
	r.Subtitle("Teste 1 – Original; todos os pares")
	r.PrintMatrix(base.Dist, r.Dist)
	printRoute(base, A, F)

	r.Subtitle("Teste 2 – Original; A → A")
	printRoute(base, A, A)

	r.Subtitle("Teste 3 – Original; F → A")
	printRoute(base, F, A)

	r.Subtitle("Teste 4 – Acrescentar G isolado")
	res4, err := solve(G+1, arestas)
	if err != nil {
		return err
	}
	r.PrintMatrix(res4.Dist, r.Dist)
	printRoute(res4, A, G)
	printRoute(res4, G, A)

	r.Subtitle("Teste 5 – Acrescentar E → B com peso −5")
	res5, err := solve(F+1, append(slices.Clone(arestas), graph.WeightedEdge{U: E, V: B, Weight: -5}))
	if err != nil {
		return err
	}
	r.PrintMatrix(res5.Dist, r.Dist)
	printCycle(res5)
	printRoute(res5, A, F)
	printRoute(res5, A, C)

	r.Subtitle("Teste 6 – Alterar A → B de 4 para 1")
	e6 := slices.Clone(arestas)
	e6[0].Weight = 1
	res6, err := solve(F+1, e6)
	if err != nil {
		return err
	}
	r.PrintMatrix(res6.Dist, r.Dist)
	printRoute(res6, A, F)
	fmt.Printf("  Comparação A → F: original %s (custo %d) × alterado %s (custo %d)\n",
		r.Path(base.Path(A, F)), base.Dist[A][F], r.Path(res6.Path(A, F)), res6.Dist[A][F])
	return nil
}

// challenge lista todas as rotas mínimas a partir de A e compara com Bellman-Ford.
func challenge(res shortest.APSP) error {
	r.Title("PARTE VII – Desafio: rotas a partir de A × Bellman-Ford")
	bf, err := shortest.BellmanFord(F+1, arestas, A)
	if err != nil {
		return err
	}
	fmt.Printf("  %-8s %-26s %-14s %s\n", "Destino", "Rota (prox)", "Floyd-Warshall", "Bellman-Ford")
	for v := range F + 1 {
		ok := "="
		if res.Dist[A][v] != bf.Dist[v] {
			ok = "≠"
		}
		fmt.Printf("  %-8s %-26s %-14s %s %s\n", r.Label(v), r.Path(res.Path(A, v)), r.Dist(res.Dist[A][v]), r.Dist(bf.Dist[v]), ok)
	}
	return nil
}

func solve(n int, edges []graph.WeightedEdge) (shortest.APSP, error) {
	dist, prox, err := shortest.NewMatrices(n, edges)
	if err != nil {
		return shortest.APSP{}, err
	}
	return shortest.FloydWarshall(dist, prox, nil), nil
}

// printStage é o "Escrever Intermediário, k, dist".
func printStage(k int, dist shortest.Matrix) {
	r.Subtitle(fmt.Sprintf("Intermediário %s – D(%s)", r.Label(k), r.Label(k)))
	r.PrintMatrix(dist, r.Dist)
}

func printCycle(res shortest.APSP) {
	fmt.Print("\n  Diagonal:")
	for v := range res.Dist {
		fmt.Printf(" %s=%s", r.Label(v), r.Dist(res.Dist[v][v]))
	}
	fmt.Println()
	if res.NegativeCycle {
		fmt.Println("  ⚠ Ciclo negativo detectado (há dist[v][v] < 0).")
	} else {
		fmt.Println("  Nenhum ciclo negativo (diagonal toda zero).")
	}
}

func printRoute(res shortest.APSP, i, j int) {
	src, dst := r.Label(i), r.Label(j)
	switch p := res.Path(i, j); {
	case res.OnNegativeCycle(i, j):
		fmt.Printf("  %s → %s: não apresentada como mínima (a rota alcança um ciclo negativo)\n", src, dst)
	case p == nil:
		fmt.Printf("  %s → %s: não existe rota (dist = ∞)\n", src, dst)
	default:
		fmt.Printf("  %s → %s: %-26s dist[%s][%s] = %d\n", src, dst, r.Path(p), src, dst, res.Dist[i][j])
	}
}
