package traversal

import (
	"slices"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/ds"
	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// DFSStack executa a busca em profundidade com pilha explícita (Questões 5 e 8).
// A cada vértice retirado do topo, onStep recebe um Step Pop.
func DFSStack(g graph.Neighborhood, start int, onStep StepFunc) []int {
	// PASSO 1 – preparação: nenhum vértice visitado.
	visited := make([]bool, g.Order())
	var order []int

	// PASSO 2 – a pilha explícita substitui a pilha de chamadas da versão
	// recursiva. O vértice inicial é o primeiro a entrar.
	var stack ds.Stack[int]
	stack.Push(start)

	for stack.Len() > 0 {
		u, _ := stack.Pop()

		// Um vértice pode ter sido empilhado por caminhos diferentes;
		// se já foi visitado, é apenas descartado.
		if visited[u] {
			if onStep != nil {
				onStep(Step{Kind: Pop, Vertex: u, Frontier: stack.Items(), Visited: append([]int{}, order...), Skipped: true})
			}
			continue
		}

		visited[u] = true
		order = append(order, u)

		// LIFO: empilhar os vizinhos em ordem decrescente faz o menor
		// ficar no topo e ser visitado primeiro.
		neighbors := g.Neighbors(u)
		var added []int
		for _, v := range slices.Backward(neighbors) {
			if !visited[v] {
				stack.Push(v)
				added = append(added, v)
			}
		}

		if onStep != nil {
			onStep(Step{
				Kind:      Pop,
				Vertex:    u,
				Neighbors: neighbors,
				Added:     added,
				Frontier:  stack.Items(),
				Visited:   append([]int{}, order...),
			})
		}
	}
	return order
}
