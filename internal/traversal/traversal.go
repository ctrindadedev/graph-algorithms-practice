// Package traversal implementa BFS e DFS. Os algoritmos não imprimem nada: retornam
// a ordem de visitação e notificam cada passo via StepFunc; a exibição fica em cmd/.
package traversal

type StepKind int

const (
	Dequeue StepKind = iota // BFS: vértice retirado da fila (Questão 7)
	Pop                     // DFS com pilha: vértice retirado do topo (Questão 8)
	Enter                   // DFS recursiva: "Entrando no vértice..." (Questão 4)
	Leave                   // DFS recursiva: "Retornando do vértice..." (Questão 4)
)

// Step é uma fotografia do estado do algoritmo em uma iteração.
// Campos que não fazem sentido para um tipo de passo ficam vazios.
type Step struct {
	Kind      StepKind
	Vertex    int   // vértice sendo processado
	Neighbors []int // vizinhos encontrados
	Added     []int // novos vértices enfileirados/empilhados
	Frontier  []int // conteúdo atual da fila (início primeiro) ou da pilha (topo primeiro)
	Visited   []int // vértices já visitados, na ordem de visitação
	Skipped   bool  // DFS com pilha: vértice retirado já tinha sido visitado
}

// StepFunc recebe cada passo do percurso. Pode ser nil quando não se
// deseja acompanhar a execução.
type StepFunc func(Step)
