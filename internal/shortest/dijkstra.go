// Package shortest implementa caminhos mínimos em grafos ponderados: Dijkstra
// (pesos não negativos), Bellman-Ford (pesos negativos, uma origem) e
// Floyd-Warshall (todos os pares).
package shortest

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"slices"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Inf marca distância infinita: vértice ainda não (ou nunca) alcançado.
const Inf = math.MaxInt

// NoPred indica que o vértice não tem predecessor (a origem ou um inalcançável).
const NoPred = -1

var ErrNegativeWeight = errors.New("shortest: Dijkstra não aceita pesos negativos")

// Result guarda o vetor de distâncias e o de predecessores a partir de Source.
// Nos estudos dirigidos, Dist é distancia[v] e Pred é anterior[v].
type Result struct {
	Source int
	Dist   []int
	Pred   []int
}

// Step é o "Escrever distancia, anterior" do fim de cada iteração do Dijkstra:
// o vértice fechado na iteração e uma cópia dos dois vetores naquele momento.
type Step struct {
	Vertex int
	Dist   []int
	Pred   []int
}

// StepFunc recebe cada iteração; pode ser nil quando o rastreamento não interessa.
type StepFunc func(Step)

// Dijkstra calcula o menor custo de src até cada vértice. Usa a versão O(V²),
// que escolhe a cada rodada o vértice aberto de menor distância.
func Dijkstra(g *graph.Weighted, src int) (Result, error) {
	return DijkstraTrace(g, src, nil)
}

// DijkstraTrace é o Dijkstra com seleção linear, avisando onStep ao fim de cada
// iteração (Questão 2 do estudo dirigido). Em empate de distância, escolhe o
// vértice de menor índice, que é a menor letra.
func DijkstraTrace(g *graph.Weighted, src int, onStep StepFunc) (Result, error) {
	if err := checkNonNegative(g); err != nil {
		return Result{}, err
	}
	n := g.Order()

	// PASSO 1 – inicializar as estruturas de controle:
	// dist = distancia[v] (∞: ainda não conhecemos rota até v),
	// pred = anterior[v] (NULO: não há predecessor registrado),
	// closed = visitado[v] (falso: a distância de v ainda não está definida).
	dist, pred := initVectors(n, src)
	closed := make([]bool, n)

	// PASSO 2 – repetir a escolha do menor custo enquanto existir vértice não
	// visitado com distância finita.
	for {
		// u ← não visitado com menor distancia[u]. Percorrer v em ordem crescente
		// e trocar só com "<" garante o desempate pela menor letra.
		u := -1
		for v := range n {
			if !closed[v] && dist[v] != Inf && (u == -1 || dist[v] < dist[u]) {
				u = v
			}
		}
		if u == -1 {
			break // os vértices restantes são inalcançáveis
		}
		// Com pesos não negativos, nenhum caminho futuro consegue baratear u.
		closed[u] = true

		// PASSO 3 – examinar as arestas que saem de u (relaxamento).
		for _, a := range g.Arcs(u) {
			if closed[a.To] {
				continue
			}
			candidate := dist[u] + a.Weight
			// Só substitui com melhora estrita: em empate, mantém o predecessor antigo.
			if candidate < dist[a.To] {
				dist[a.To] = candidate
				pred[a.To] = u // guardar u para reconstruir o caminho até v
			}
		}
		if onStep != nil {
			onStep(Step{Vertex: u, Dist: slices.Clone(dist), Pred: slices.Clone(pred)})
		}
	}

	// PASSO 4 – devolver os resultados.
	return Result{Source: src, Dist: dist, Pred: pred}, nil
}

// DijkstraHeap é a versão com fila de prioridade (desafio): em vez de varrer os
// n vértices para achar o menor, retira o topo de um heap binário. Custo
// O((V + E) log V). Usa "inserção preguiçosa": um vértice pode entrar mais de
// uma vez no heap, e as entradas desatualizadas são descartadas ao sair.
func DijkstraHeap(g *graph.Weighted, src int, onStep StepFunc) (Result, error) {
	if err := checkNonNegative(g); err != nil {
		return Result{}, err
	}
	n := g.Order()
	dist, pred := initVectors(n, src)
	closed := make([]bool, n)

	pq := &minHeap{{vertex: src, dist: 0}}
	for pq.Len() > 0 {
		top := heap.Pop(pq).(entry)
		u := top.vertex
		if closed[u] || top.dist != dist[u] {
			continue // entrada antiga: u já foi fechado com uma distância menor
		}
		closed[u] = true

		for _, a := range g.Arcs(u) {
			if closed[a.To] {
				continue
			}
			if candidate := dist[u] + a.Weight; candidate < dist[a.To] {
				dist[a.To] = candidate
				pred[a.To] = u
				heap.Push(pq, entry{vertex: a.To, dist: candidate})
			}
		}
		if onStep != nil {
			onStep(Step{Vertex: u, Dist: slices.Clone(dist), Pred: slices.Clone(pred)})
		}
	}
	return Result{Source: src, Dist: dist, Pred: pred}, nil
}

func checkNonNegative(g *graph.Weighted) error {
	for u := range g.Order() {
		for _, a := range g.Arcs(u) {
			if a.Weight < 0 {
				return fmt.Errorf("%w: %d → %d com peso %d", ErrNegativeWeight, u, a.To, a.Weight)
			}
		}
	}
	return nil
}

// initVectors devolve distancia (∞ em todos, 0 na origem) e anterior (NULO).
func initVectors(n, src int) (dist, pred []int) {
	dist = make([]int, n)
	pred = make([]int, n)
	for v := range n {
		dist[v], pred[v] = Inf, NoPred
	}
	dist[src] = 0 // custo para chegar à origem = zero
	return dist, pred
}

// PathTo reconstrói o caminho mínimo seguindo os predecessores de v até a
// origem (ReconstruirCaminho). Retorna nil se v for inalcançável ou se a cadeia
// de predecessores não chegar à origem, o que pode acontecer quando há ciclo
// negativo. O laço dá no máximo n passos, então nunca fica preso em um ciclo.
func (r Result) PathTo(v int) []int {
	// PASSO 1 – verificar se o destino é alcançável.
	if r.Dist[v] == Inf {
		return nil
	}
	// PASSO 2 – percorrer os predecessores de trás para frente.
	var path []int
	for x, steps := v, 0; x != NoPred && steps <= len(r.Pred); x, steps = r.Pred[x], steps+1 {
		path = append(path, x)
		if x == r.Source {
			// PASSO 3 – inverter para entregar a rota da origem ao destino.
			slices.Reverse(path)
			return path
		}
	}
	return nil // a cadeia não chegou à origem
}

// entry é um par (vértice, distância no momento da inserção) do heap.
type entry struct {
	vertex, dist int
}

// minHeap ordena por distância e, em empate, pelo menor vértice (menor letra).
type minHeap []entry

func (h minHeap) Len() int { return len(h) }
func (h minHeap) Less(i, j int) bool {
	if h[i].dist != h[j].dist {
		return h[i].dist < h[j].dist
	}
	return h[i].vertex < h[j].vertex
}
func (h minHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x any)   { *h = append(*h, x.(entry)) }
func (h *minHeap) Pop() any {
	old := *h
	x := old[len(old)-1]
	*h = old[:len(old)-1]
	return x
}
