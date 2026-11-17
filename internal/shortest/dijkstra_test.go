package shortest

import (
	"errors"
	"reflect"
	"testing"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Digrafo ponderado da Lista 03 com A..E = 0..4.
var lista03 = []graph.WeightedEdge{{U: 0, V: 1, Weight: 4}, {U: 0, V: 2, Weight: 2}, {U: 2, V: 1, Weight: 1}, {U: 1, V: 3, Weight: 5}, {U: 2, V: 3, Weight: 8}, {U: 2, V: 4, Weight: 10}, {U: 3, V: 4, Weight: 2}}

func TestDijkstra(t *testing.T) {
	g, err := graph.BuildWeighted(5, true, lista03)
	if err != nil {
		t.Fatal(err)
	}
	r, err := Dijkstra(g, 0)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int{0, 3, 2, 8, 10}; !reflect.DeepEqual(r.Dist, want) {
		t.Errorf("Dist = %v, want %v", r.Dist, want)
	}
	if want := []int{NoPred, 2, 0, 1, 3}; !reflect.DeepEqual(r.Pred, want) {
		t.Errorf("Pred = %v, want %v", r.Pred, want)
	}
	if want := []int{0, 2, 1, 3, 4}; !reflect.DeepEqual(r.PathTo(4), want) {
		t.Errorf("PathTo(E) = %v, want %v", r.PathTo(4), want)
	}
}

func TestDijkstraUnreachable(t *testing.T) {
	g, _ := graph.BuildWeighted(5, true, lista03)
	r, _ := Dijkstra(g, 3)
	if r.Dist[0] != Inf || r.PathTo(0) != nil {
		t.Errorf("A deveria ser inalcançável a partir de D: dist = %d", r.Dist[0])
	}
	if r.Dist[4] != 2 {
		t.Errorf("Dist[E] = %d, want 2", r.Dist[4])
	}
}

func TestDijkstraNegativeWeight(t *testing.T) {
	g, _ := graph.BuildWeighted(2, true, []graph.WeightedEdge{{U: 0, V: 1, Weight: -1}})
	if _, err := Dijkstra(g, 0); !errors.Is(err, ErrNegativeWeight) {
		t.Errorf("err = %v, want ErrNegativeWeight", err)
	}
}
