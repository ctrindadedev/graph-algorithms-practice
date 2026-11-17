package traversal

import "github.com/ctrindadedev/graph-algorithms-practice/internal/graph"

// DFSRecursive executa a busca em profundidade recursiva (Questão 4). onStep recebe
// um Step Enter ao entrar em cada vértice e um Step Leave ao retornar dele.
func DFSRecursive(g graph.Graph, start int, onStep StepFunc) []int {
	// Inicialmente nenhum vértice foi visitado.
	visited := make([]bool, g.Order())
	var order []int

	notify := func(kind StepKind, u int) {
		if onStep != nil {
			onStep(Step{Kind: kind, Vertex: u, Neighbors: g.Neighbors(u), Visited: append([]int{}, order...)})
		}
	}

	// A pilha está implícita: cada chamada de visit ocupa um quadro na
	// pilha de chamadas do programa, que é desempilhado ao retornar.
	var visit func(u int)
	visit = func(u int) {
		visited[u] = true
		order = append(order, u)
		notify(Enter, u)

		for _, v := range g.Neighbors(u) {
			if !visited[v] {
				// Aprofunda: a chamada para v termina antes de olhar o
				// próximo vizinho de u.
				visit(v)
			}
		}

		// Sem vizinhos não visitados, retorna para a chamada anterior.
		notify(Leave, u)
	}

	visit(start)
	return order
}
