// Package tree implementa os algoritmos da Aula 9: verificação de árvore,
// árvores geradoras por BFS e por DFS e arborescência por BFS.
package tree

import (
	"fmt"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/ds"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/traversal"
)

// IsTree verifica se o grafo não direcionado g é uma árvore.
func IsTree(g graph.Graph) bool {
	n := g.Order()
	if n == 0 || graph.Size(g) != n-1 {
		return false
	}
	// V − 1 sozinho não basta: um componente pode ter ciclo enquanto outro fica
	// desconectado. Por isso a BFS confere se todos os vértices são alcançados.
	visited := traversal.BFS(g, 0, nil)
	return len(visited) == n
}

// SpanningTreeBFS constrói uma árvore geradora com as arestas de descoberta
// da busca em largura a partir de start, na ordem em que são descobertas.
func SpanningTreeBFS(g graph.Graph, start int) []graph.Edge {
	// visited impede que um vértice seja descoberto duas vezes; T guarda
	// somente as arestas da árvore.
	visited := make([]bool, g.Order())
	var t []graph.Edge

	// Fila FIFO da BFS, iniciada pela raiz já marcada.
	var queue ds.Queue[int]
	visited[start] = true
	queue.Enqueue(start)

	for queue.Len() > 0 {
		u, _ := queue.Dequeue()
		for _, v := range g.Neighbors(u) {
			if !visited[v] {
				visited[v] = true
				// v acabou de ser descoberto a partir de u: {u, v} entra na árvore.
				t = append(t, graph.Edge{U: u, V: v})
				queue.Enqueue(v)
			}
		}
	}
	return t
}

// noParent indica que o vértice é a raiz da busca (o NULO do pseudocódigo).
const noParent = -1

type stackEntry struct {
	vertex, parent int
}

// SpanningTreeDFS constrói uma árvore geradora por DFS com pilha explícita; cada
// entrada da pilha guarda o vértice e o pai pelo qual ele foi alcançado.
func SpanningTreeDFS(g graph.Graph, start int) []graph.Edge {
	// visited garante que cada vértice seja processado uma única vez.
	visited := make([]bool, g.Order())
	var t []graph.Edge

	// Pilha LIFO de pares (vértice, pai); a raiz não tem pai.
	var stack ds.Stack[stackEntry]
	stack.Push(stackEntry{start, noParent})

	for stack.Len() > 0 {
		e, _ := stack.Pop()
		if visited[e.vertex] {
			continue
		}
		visited[e.vertex] = true

		// Só no momento em que o vértice é visitado a aresta (pai, u) vira
		// aresta da árvore; entradas repetidas foram descartadas acima.
		if e.parent != noParent {
			t = append(t, graph.Edge{U: e.parent, V: e.vertex})
		}

		neighbors := g.Neighbors(e.vertex)
		for i := len(neighbors) - 1; i >= 0; i-- {
			if v := neighbors[i]; !visited[v] {
				stack.Push(stackEntry{v, e.vertex})
			}
		}
	}
	return t
}

// UnreachableError indica que a arborescência não existe com a raiz
// escolhida, pois alguns vértices não são alcançáveis a partir dela.
type UnreachableError struct {
	Root        int
	Unreachable []int
}

func (e *UnreachableError) Error() string {
	return fmt.Sprintf("falha: vértices %s inalcançáveis a partir da raiz %d", graph.FormatInts(e.Unreachable), e.Root)
}

// ArborescenceBFS constrói uma arborescência de raiz root seguindo só os arcos de
// saída. Retorna *UnreachableError se algum vértice não for alcançado.
func ArborescenceBFS(d *graph.Digraph, root int) ([]graph.Edge, error) {
	visited := make([]bool, d.Order())
	var a []graph.Edge

	var queue ds.Queue[int]
	visited[root] = true
	queue.Enqueue(root)

	for queue.Len() > 0 {
		u, _ := queue.Dequeue()
		// Successors respeita a orientação: apenas arcos u → v.
		for _, v := range d.Successors(u) {
			if !visited[v] {
				visited[v] = true
				a = append(a, graph.Edge{U: u, V: v})
				queue.Enqueue(v)
			}
		}
	}

	var unreachable []int
	for v, ok := range visited {
		if !ok {
			unreachable = append(unreachable, v)
		}
	}
	if len(unreachable) > 0 {
		return nil, &UnreachableError{Root: root, Unreachable: unreachable}
	}
	return a, nil
}
