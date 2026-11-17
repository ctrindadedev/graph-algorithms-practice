package ds

import (
	"reflect"
	"testing"
)

func TestQueueFIFO(t *testing.T) {
	var q Queue[int]
	for _, x := range []int{1, 2, 3} {
		q.Enqueue(x)
	}
	if got, want := q.Items(), []int{1, 2, 3}; !reflect.DeepEqual(got, want) {
		t.Errorf("Items() = %v, quero %v", got, want)
	}
	for _, want := range []int{1, 2, 3} {
		if got, ok := q.Dequeue(); !ok || got != want {
			t.Errorf("Dequeue() = %d, %v; quero %d, true", got, ok, want)
		}
	}
	if _, ok := q.Dequeue(); ok {
		t.Error("Dequeue() em fila vazia deveria retornar ok=false")
	}
}

func TestStackLIFO(t *testing.T) {
	var s Stack[int]
	for _, x := range []int{1, 2, 3} {
		s.Push(x)
	}
	if got, want := s.Items(), []int{3, 2, 1}; !reflect.DeepEqual(got, want) {
		t.Errorf("Items() = %v, quero %v (topo primeiro)", got, want)
	}
	for _, want := range []int{3, 2, 1} {
		if got, ok := s.Pop(); !ok || got != want {
			t.Errorf("Pop() = %d, %v; quero %d, true", got, ok, want)
		}
	}
	if _, ok := s.Pop(); ok {
		t.Error("Pop() em pilha vazia deveria retornar ok=false")
	}
}
