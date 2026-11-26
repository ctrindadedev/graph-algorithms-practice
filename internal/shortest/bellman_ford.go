package shortest

import (
	"fmt"
	"slices"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Pass é o "Escrever i, distancia, anterior" do fim de cada passada do
// Bellman-Ford. Changed é a variável mudou do pseudocódigo.
type Pass struct {
	Index   int // número da passada, a partir de 1
	Dist    []int
	Pred    []int
	Changed bool
}

// BFResult acrescenta ao Result as passadas executadas e o teste de ciclo negativo.
type BFResult struct {
	Result
	Passes        []Pass
	NegativeCycle bool
}

// BellmanFord calcula os menores custos a partir de src aceitando pesos
// negativos. Recebe a LISTA de arestas, e não a lista de adjacência, porque a
// ordem em que as arestas são relaxadas muda os estados intermediários das
// passadas (o enunciado fixa essa ordem). Laços u → u são aceitos, já que um
// laço negativo é um ciclo negativo (teste 6).
func BellmanFord(n int, edges []graph.WeightedEdge, src int) (BFResult, error) {
	for _, e := range edges {
		if e.U < 0 || e.U >= n || e.V < 0 || e.V >= n {
			return BFResult{}, fmt.Errorf("%w: (%d, %d) com n = %d", graph.ErrInvalidVertex, e.U, e.V, n)
		}
	}

	// PASSO 1 – inicializar distâncias (∞; origem 0) e predecessores (NULO).
	dist, pred := initVectors(n, src)
	var passes []Pass

	// PASSO 2 – repetir até |V|−1 passadas por TODAS as arestas. Um caminho
	// mínimo simples tem no máximo |V|−1 arestas; a passada i garante as rotas
	// mínimas de até i arestas.
	for i := 1; i <= n-1; i++ {
		changed := false // mudou ← falso
		for _, e := range edges {
			// Só relaxa a partir de u alcançável: ∞ + peso não é uma rota.
			if dist[e.U] == Inf {
				continue
			}
			candidate := dist[e.U] + e.Weight
			if candidate < dist[e.V] {
				dist[e.V] = candidate // relaxamento da aresta u → v
				pred[e.V] = e.U       // salvar predecessor
				changed = true
			}
		}
		passes = append(passes, Pass{Index: i, Dist: slices.Clone(dist), Pred: slices.Clone(pred), Changed: changed})
		// Parada antecipada: se uma passada inteira não melhorou nada, as
		// seguintes repetiriam exatamente os mesmos testes.
		if !changed {
			break
		}
	}

	// PASSO 3 – testar ciclo negativo alcançável: depois de |V|−1 passadas, as
	// distâncias já são mínimas, a não ser que haja um ciclo negativo. Se alguma
	// aresta com u alcançável ainda melhora v, o ciclo existe.
	negative := false
	for _, e := range edges {
		if dist[e.U] != Inf && dist[e.U]+e.Weight < dist[e.V] {
			negative = true
			break
		}
	}

	// PASSO 4 – devolver resultados.
	return BFResult{Result: Result{Source: src, Dist: dist, Pred: pred}, Passes: passes, NegativeCycle: negative}, nil
}
