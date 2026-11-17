// Package ds contém estruturas de dados auxiliares reaproveitadas pelos algoritmos de grafos.
package ds

// Queue é uma fila genérica FIFO – First In, First Out (usada pela BFS).
type Queue[T any] struct {
	items []T
}

func (q *Queue[T]) Enqueue(x T) {
	q.items = append(q.items, x)
}

// Dequeue remove e retorna o primeiro elemento. ok é false se a fila estiver vazia.
func (q *Queue[T]) Dequeue() (x T, ok bool) {
	if len(q.items) == 0 {
		return x, false
	}
	x = q.items[0]
	q.items = q.items[1:]
	return x, true
}

func (q *Queue[T]) Len() int {
	return len(q.items)
}

// Items retorna uma cópia do conteúdo, do primeiro ao último, para exibição do estado da fila (Questão 7).
func (q *Queue[T]) Items() []T {
	return append([]T{}, q.items...)
}
