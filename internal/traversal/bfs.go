package traversal

import (
	"github.com/ctrindadedev/graph-algorithms-practice/internal/ds"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// BFS executa a busca em largura a partir de start (Questões 3 e 7). A cada vértice
// retirado da fila, onStep (se não for nil) recebe um Step Dequeue.
func BFS(g graph.Graph, start int, onStep StepFunc) []int {
	// Inicialmente nenhum vértice foi visitado.
	visited := make([]bool, g.Order())
	var marked, order []int

	// O vértice inicial é marcado e colocado na fila, que guarda os
	// vértices descobertos mas ainda não processados.
	var queue ds.Queue[int]
	visited[start] = true
	marked = append(marked, start)
	queue.Enqueue(start)

	for queue.Len() > 0 {
		// FIFO: o primeiro a entrar é o primeiro a ser processado, o que
		// faz a busca avançar nível por nível.
		u, _ := queue.Dequeue()
		order = append(order, u)

		neighbors := g.Neighbors(u)
		var added []int
		for _, v := range neighbors {
			if !visited[v] {
				// Marca antes de enfileirar para que v não entre duas vezes.
				visited[v] = true
				marked = append(marked, v)
				added = append(added, v)
				queue.Enqueue(v)
			}
		}

		if onStep != nil {
			onStep(Step{
				Kind:      Dequeue,
				Vertex:    u,
				Neighbors: neighbors,
				Added:     added,
				Frontier:  queue.Items(),
				Visited:   append([]int{}, marked...),
			})
		}
	}
	return order
}
