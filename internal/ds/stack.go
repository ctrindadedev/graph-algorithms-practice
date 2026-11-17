package ds

// Stack é uma pilha genérica LIFO – Last In, First Out (usada pela DFS iterativa).
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(x T) {
	s.items = append(s.items, x)
}

// Pop remove e retorna o elemento do topo. ok é false se a pilha
// estiver vazia.
func (s *Stack[T]) Pop() (x T, ok bool) {
	if len(s.items) == 0 {
		return x, false
	}
	top := len(s.items) - 1
	x = s.items[top]
	s.items = s.items[:top]
	return x, true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

// Items retorna uma cópia do conteúdo com o TOPO PRIMEIRO, no mesmo
// formato usado no acompanhamento do enunciado: Pilha: [3, 4, 2].
func (s *Stack[T]) Items() []T {
	out := make([]T, len(s.items))
	for i, x := range s.items {
		out[len(s.items)-1-i] = x
	}
	return out
}
