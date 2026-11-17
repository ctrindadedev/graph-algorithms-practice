package traversal

import (
	"reflect"
	"testing"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Grafos da Lista 03 com A..E = 0..4; os pesos não alteram a ordem dos percursos.
func TestLista03Traversals(t *testing.T) {
	directed, _ := graph.BuildDigraph(5, []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 2, V: 3}, {U: 3, V: 4}})
	weighted, _ := graph.BuildWeighted(5, false, []graph.WeightedEdge{{U: 0, V: 1, Weight: 4}, {U: 0, V: 2, Weight: 2}, {U: 1, V: 2, Weight: 1}, {U: 1, V: 3, Weight: 5}, {U: 2, V: 3, Weight: 8}, {U: 2, V: 4, Weight: 10}, {U: 3, V: 4, Weight: 2}})
	both, _ := graph.BuildWeighted(5, true, []graph.WeightedEdge{{U: 0, V: 1, Weight: 4}, {U: 0, V: 2, Weight: 2}, {U: 2, V: 1, Weight: 1}, {U: 1, V: 3, Weight: 5}, {U: 2, V: 3, Weight: 8}, {U: 2, V: 4, Weight: 10}, {U: 3, V: 4, Weight: 2}})

	tests := []struct {
		name     string
		g        graph.Neighborhood
		bfs, dfs []int
	}{
		{"direcionado", directed, []int{0, 1, 2, 3, 4}, []int{0, 1, 3, 4, 2}},
		{"ponderado", weighted, []int{0, 1, 2, 3, 4}, []int{0, 1, 2, 3, 4}},
		{"direcionado e ponderado", both, []int{0, 1, 2, 3, 4}, []int{0, 1, 3, 4, 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BFS(tt.g, 0, nil); !reflect.DeepEqual(got, tt.bfs) {
				t.Errorf("BFS = %v, want %v", got, tt.bfs)
			}
			if got := DFSRecursive(tt.g, 0, nil); !reflect.DeepEqual(got, tt.dfs) {
				t.Errorf("DFSRecursive = %v, want %v", got, tt.dfs)
			}
			if got := DFSStack(tt.g, 0, nil); !reflect.DeepEqual(got, tt.dfs) {
				t.Errorf("DFSStack = %v, want %v (deveria reproduzir a recursiva)", got, tt.dfs)
			}
		})
	}
}

func TestLista03UnreachableFromD(t *testing.T) {
	directed, _ := graph.BuildDigraph(5, []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 2, V: 3}, {U: 3, V: 4}})
	if got, want := BFS(directed, 3, nil), []int{3, 4}; !reflect.DeepEqual(got, want) {
		t.Errorf("BFS a partir de D = %v, want %v", got, want)
	}
}
