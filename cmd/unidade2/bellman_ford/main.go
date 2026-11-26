// Comando bellman_ford: estudo dirigido do algoritmo de Bellman-Ford (Unidade II).
// Executa todas as partes em sequência. Flags (desafio): -start escolhe a
// origem (A..F) e -reverse relaxa as arestas na ordem inversa à do enunciado.
package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	r "github.com/ctrindadedev/graph-algorithms-practice/internal/report"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/shortest"
)

// A..F = 0..5, então a ordem crescente dos índices é a alfabética.
const (
	A = iota
	B
	C
	D
	E
	F
	G // vértice isolado dos testes 4 e 6
)

// arestas na ordem de inserção da Figura 1. A ordem importa: é nela que cada
// passada relaxa as arestas.
var arestas = []graph.WeightedEdge{
	{U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: C, V: B, Weight: -1},
	{U: B, V: D, Weight: 2}, {U: C, V: D, Weight: 5}, {U: D, V: E, Weight: -2},
	{U: E, V: F, Weight: 3}, {U: B, V: F, Weight: 10}, {U: C, V: E, Weight: 7},
}

func main() {
	startFlag := flag.String("start", "A", "vértice de origem (A..F)")
	reverse := flag.Bool("reverse", false, "relaxar as arestas na ordem inversa (desafio)")
	flag.Parse()

	s := strings.ToUpper(*startFlag)
	if len(s) != 1 || s[0] < 'A' || s[0] > 'F' {
		fmt.Fprintf(os.Stderr, "erro: origem %q inválida, use uma letra de A a F\n", *startFlag)
		os.Exit(1)
	}
	if err := run(int(s[0]-'A'), *reverse); err != nil {
		fmt.Fprintln(os.Stderr, "erro:", err)
		os.Exit(1)
	}
}

func run(origem int, reverse bool) error {
	edges := arestas
	if reverse {
		edges = slices.Clone(arestas)
		slices.Reverse(edges)
	}

	// PARTE I – criar_grafo: cada arco só no sentido u → v.
	r.Title("PARTE I – Lista de adjacência (direcionada)")
	printAdjacency(F+1, arestas)
	fmt.Print("  Arestas negativas:")
	for _, e := range arestas {
		if e.Weight < 0 {
			fmt.Printf(" %s → %s (%d)", r.Label(e.U), r.Label(e.V), e.Weight)
		}
	}
	fmt.Println()

	// PARTES II e III – passadas.
	title := fmt.Sprintf("PARTES II e III – Bellman-Ford a partir de %s", r.Label(origem))
	if reverse {
		title += " (arestas em ordem inversa)"
	}
	r.Title(title)
	fmt.Println("  Ordem de relaxamento:", edgeOrder(edges))
	res, err := shortest.BellmanFord(F+1, edges, origem)
	if err != nil {
		return err
	}
	printPasses(res)

	// PARTE IV – reconstruir_caminho.
	r.Title("PARTE IV – Reconstrução dos caminhos mínimos")
	for _, dest := range []int{F, B, D, E} {
		printRoute(res, dest)
	}

	if origem != A || reverse {
		return nil // as variações do desafio só repetem as partes II a IV
	}

	r.Title("PARTE VI – Testes")
	return runTests(res)
}

func runTests(base shortest.BFResult) error {
	r.Subtitle("Teste 1 – Original; A → F")
	printRoute(base, F)

	r.Subtitle("Teste 2 – Original; A → A")
	printRoute(base, A)

	r.Subtitle("Teste 3 – Original; A → D")
	printRoute(base, D)
	fmt.Printf("  Predecessores: anterior[D] = %s, anterior[B] = %s, anterior[C] = %s\n",
		r.Pred(base.Pred[D]), r.Pred(base.Pred[B]), r.Pred(base.Pred[C]))

	// Teste 4: 7 vértices, G sem nenhuma aresta.
	r.Subtitle("Teste 4 – G isolado; A → G")
	res4, err := shortest.BellmanFord(G+1, arestas, A)
	if err != nil {
		return err
	}
	printPasses(res4)
	printRoute(res4, G)

	// Teste 5: E → B (−5) fecha o ciclo B → D → E → B com peso 2 − 2 − 5 = −5.
	r.Subtitle("Teste 5 – Acrescentar E → B com peso −5")
	e5 := append(slices.Clone(arestas), graph.WeightedEdge{U: E, V: B, Weight: -5})
	res5, err := shortest.BellmanFord(F+1, e5, A)
	if err != nil {
		return err
	}
	printPasses(res5)
	printRoute(res5, F)

	// Teste 6: laço G → G (−1) é um ciclo negativo, mas G não é alcançável.
	r.Subtitle("Teste 6 – Acrescentar G → G com peso −1 (G isolado)")
	e6 := append(slices.Clone(arestas), graph.WeightedEdge{U: G, V: G, Weight: -1})
	res6, err := shortest.BellmanFord(G+1, e6, A)
	if err != nil {
		return err
	}
	printPasses(res6)
	printRoute(res6, F)
	printRoute(res6, G)
	return nil
}

// printAdjacency monta a lista na ordem de inserção das arestas.
func printAdjacency(n int, edges []graph.WeightedEdge) {
	adj := make([][]graph.WeightedEdge, n)
	for _, e := range edges {
		adj[e.U] = append(adj[e.U], e) // apenas u → v, sem ligação de retorno
	}
	for u, list := range adj {
		parts := make([]string, len(list))
		for i, e := range list {
			parts[i] = fmt.Sprintf("(%s, %d)", r.Label(e.V), e.Weight)
		}
		fmt.Printf("  %s → [%s]\n", r.Label(u), strings.Join(parts, ", "))
	}
}

func edgeOrder(edges []graph.WeightedEdge) string {
	parts := make([]string, len(edges))
	for i, e := range edges {
		parts[i] = fmt.Sprintf("%s→%s(%d)", r.Label(e.U), r.Label(e.V), e.Weight)
	}
	return strings.Join(parts, ", ")
}

// printPasses mostra o "Escrever i, distancia, anterior" de cada passada.
func printPasses(res shortest.BFResult) {
	n := len(res.Dist)
	fmt.Println()
	r.Header("Etapa", 12, n, 5)
	initial := make([]string, n)
	for v := range initial {
		initial[v] = "∞"
	}
	initial[res.Source] = "0"
	r.Row("Inicial", 12, initial, 5)
	for _, p := range res.Passes {
		r.Row(fmt.Sprintf("Passada %d", p.Index), 12, r.Dists(p.Dist), 5)
	}
	fmt.Println()
	r.Header("anterior", 12, n, 5)
	for _, p := range res.Passes {
		r.Row(fmt.Sprintf("Passada %d", p.Index), 12, r.Preds(p.Pred), 5)
	}
	last := res.Passes[len(res.Passes)-1]
	switch {
	case !last.Changed:
		fmt.Printf("\n  Passada %d não alterou nada: parada antecipada (%d de %d passadas executadas).\n", last.Index, len(res.Passes), n-1)
	default:
		fmt.Printf("\n  Todas as %d passadas alteraram alguma distância.\n", n-1)
	}
	if res.NegativeCycle {
		fmt.Println("  ⚠ Ciclo negativo alcançável detectado: as distâncias não são mínimas bem definidas.")
	} else {
		fmt.Println("  Nenhum ciclo negativo alcançável.")
	}
}

func printRoute(res shortest.BFResult, dest int) {
	src := r.Label(res.Source)
	switch p := res.PathTo(dest); {
	case res.NegativeCycle:
		// Não reconstrói: com ciclo negativo a "rota mínima" não existe.
		fmt.Printf("  %s → %s: não reconstruída (ciclo negativo alcançável)\n", src, r.Label(dest))
	case p == nil:
		fmt.Printf("  %s → %s: inalcançável (distância ∞, caminho [])\n", src, r.Label(dest))
	default:
		fmt.Printf("  %s → %s: %-26s custo %d, %d aresta(s)\n", src, r.Label(dest), r.Path(p), res.Dist[dest], len(p)-1)
	}
}
