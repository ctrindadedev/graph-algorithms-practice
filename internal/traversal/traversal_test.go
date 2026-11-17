package traversal

import (
	"reflect"
	"testing"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

var sampleEdges = []graph.Edge{{U: 0, V: 1}, {U: 0, V: 2}, {U: 1, V: 3}, {U: 1, V: 4}, {U: 2, V: 4}, {U: 3, V: 5}, {U: 4, V: 5}}

// Grafo da Questão 10: original + aresta (2,5).
var extendedEdges = append(append([]graph.Edge{}, sampleEdges...), graph.Edge{U: 2, V: 5})

func newList(n int) graph.Graph { return graph.NewAdjacencyList(n) }

func TestTraversals(t *testing.T) {
	algorithms := map[string]func(graph.Neighborhood, int, StepFunc) []int{
		"BFS":          BFS,
		"DFSRecursive": DFSRecursive,
		"DFSStack":     DFSStack,
	}

	tests := []struct {
		name  string
		edges []graph.Edge
		want  map[string][]int
	}{
		{
			name:  "grafo original",
			edges: sampleEdges,
			want: map[string][]int{
				"BFS":          {0, 1, 2, 3, 4, 5},
				"DFSRecursive": {0, 1, 3, 5, 4, 2},
				"DFSStack":     {0, 1, 3, 5, 4, 2},
			},
		},
		{
			name:  "com aresta (2,5)",
			edges: extendedEdges,
			want: map[string][]int{
				"BFS":          {0, 1, 2, 3, 4, 5},
				"DFSRecursive": {0, 1, 3, 5, 2, 4},
				"DFSStack":     {0, 1, 3, 5, 2, 4},
			},
		},
	}

	for _, tt := range tests {
		g, err := graph.Build(newList, 6, tt.edges)
		if err != nil || g == nil {
			t.Fatalf("%s: Build: g=%v err=%v", tt.name, g, err)
		}
		for algo, run := range algorithms {
			t.Run(tt.name+"/"+algo, func(t *testing.T) {
				if got := run(g, 0, nil); !reflect.DeepEqual(got, tt.want[algo]) {
					t.Errorf("ordem = %v, quero %v", got, tt.want[algo])
				}
			})
		}
	}
}

// Confere o acompanhamento passo a passo do enunciado (Questão 5).
func TestDFSStackSteps(t *testing.T) {
	g, err := graph.Build(newList, 6, sampleEdges)
	if err != nil || g == nil {
		t.Fatalf("Build: g=%v err=%v", g, err)
	}

	var frontiers [][]int
	DFSStack(g, 0, func(s Step) {
		frontiers = append(frontiers, s.Frontier)
	})

	// Estado da pilha (topo primeiro) após processar 0, 1, 3, 5, 4 e 2.
	want := [][]int{{1, 2}, {3, 4, 2}, {5, 4, 2}, {4, 4, 2}, {2, 4, 2}, {4, 2}}
	if len(frontiers) < len(want) {
		t.Fatalf("recebi %d passos, quero pelo menos %d", len(frontiers), len(want))
	}
	for i, w := range want {
		if !reflect.DeepEqual(frontiers[i], w) {
			t.Errorf("passo %d: pilha = %v, quero %v", i, frontiers[i], w)
		}
	}
}
