// Package shortest implementa caminho mínimo em grafos ponderados (Dijkstra).
package shortest

import (
	"errors"
	"fmt"
	"math"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Inf marca distância infinita: vértice ainda não (ou nunca) alcançado.
const Inf = math.MaxInt

// NoPred indica que o vértice não tem predecessor (a origem ou um inalcançável).
const NoPred = -1

var ErrNegativeWeight = errors.New("shortest: Dijkstra não aceita pesos negativos")

// Result guarda o vetor de distâncias e o de predecessores a partir de Source.
type Result struct {
	Source int
	Dist   []int
	Pred   []int
}

// Dijkstra calcula o menor custo de src até cada vértice. Usa a versão O(V²),
// que escolhe a cada rodada o vértice aberto de menor distância.
func Dijkstra(g *graph.Weighted, src int) (Result, error) {
	n := g.Order()
	for u := range n {
		for _, a := range g.Arcs(u) {
			if a.Weight < 0 {
				return Result{}, fmt.Errorf("%w: %d → %d com peso %d", ErrNegativeWeight, u, a.To, a.Weight)
			}
		}
	}

	// PASSO 1 – todas as distâncias começam em ∞, exceto a origem (0).
	dist := make([]int, n)
	pred := make([]int, n)
	closed := make([]bool, n)
	for v := range n {
		dist[v], pred[v] = Inf, NoPred
	}
	dist[src] = 0

	for range n {
		// PASSO 2 – fecha o vértice aberto de menor distância; com pesos não
		// negativos, nenhum caminho futuro consegue melhorá-la.
		u := -1
		for v := range n {
			if !closed[v] && dist[v] != Inf && (u == -1 || dist[v] < dist[u]) {
				u = v
			}
		}
		if u == -1 {
			break // os vértices restantes são inalcançáveis
		}
		closed[u] = true

		// PASSO 3 – relaxamento: passar por u pode baratear o caminho até v.
		for _, a := range g.Arcs(u) {
			if !closed[a.To] && dist[u]+a.Weight < dist[a.To] {
				dist[a.To] = dist[u] + a.Weight
				pred[a.To] = u
			}
		}
	}
	return Result{Source: src, Dist: dist, Pred: pred}, nil
}

// PathTo reconstrói o caminho mínimo seguindo os predecessores de v até a
// origem. Retorna nil se v for inalcançável.
func (r Result) PathTo(v int) []int {
	if r.Dist[v] == Inf {
		return nil
	}
	var path []int
	for x := v; x != NoPred; x = r.Pred[x] {
		path = append([]int{x}, path...)
	}
	return path
}
