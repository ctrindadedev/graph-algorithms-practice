package graph

import (
	"reflect"
	"testing"
)

func TestWeightedUndirectedStoresBothSides(t *testing.T) {
	g, err := BuildWeighted(3, false, []WeightedEdge{{0, 2, 7}, {0, 1, 4}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []Arc{{1, 4}, {2, 7}}; !reflect.DeepEqual(g.Arcs(0), want) {
		t.Errorf("Arcs(0) = %v, want %v", g.Arcs(0), want)
	}
	if w, ok := g.Weight(2, 0); !ok || w != 7 {
		t.Errorf("Weight(2, 0) = %d, %v; want 7, true", w, ok)
	}
	if want := []int{1, 2}; !reflect.DeepEqual(g.Neighbors(0), want) {
		t.Errorf("Neighbors(0) = %v, want %v", g.Neighbors(0), want)
	}
}

func TestWeightedDirectedHasNoReverseArc(t *testing.T) {
	g, err := BuildWeighted(3, true, []WeightedEdge{{2, 1, 1}})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := g.Weight(1, 2); ok {
		t.Error("arco 1 → 2 não deveria existir")
	}
	if len(g.Arcs(1)) != 0 {
		t.Errorf("Arcs(1) = %v, want vazio", g.Arcs(1))
	}
}

func TestWeightedRejectsInvalidEdges(t *testing.T) {
	g := NewWeighted(2, false)
	if g.AddEdge(0, 0, 1) == nil {
		t.Error("laço deveria ser rejeitado")
	}
	if g.AddEdge(0, 5, 1) == nil {
		t.Error("vértice fora do intervalo deveria ser rejeitado")
	}
}
