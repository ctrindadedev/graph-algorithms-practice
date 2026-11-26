package shortest

import (
	"reflect"
	"testing"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Vértices A..F = 0..5 dos estudos dirigidos da Unidade II.
const (
	vA = iota
	vB
	vC
	vD
	vE
	vF
)

const inf = Inf

// Grafo não direcionado do estudo de Dijkstra.
var dijkstraEdges = []graph.WeightedEdge{
	{U: vA, V: vB, Weight: 4}, {U: vA, V: vC, Weight: 2}, {U: vB, V: vC, Weight: 1},
	{U: vB, V: vD, Weight: 5}, {U: vC, V: vD, Weight: 8}, {U: vC, V: vE, Weight: 10},
	{U: vD, V: vE, Weight: 2}, {U: vD, V: vF, Weight: 6}, {U: vE, V: vF, Weight: 3},
}

// Digrafo de Bellman-Ford e Floyd-Warshall, na ordem do enunciado.
var negEdges = []graph.WeightedEdge{
	{U: vA, V: vB, Weight: 4}, {U: vA, V: vC, Weight: 2}, {U: vC, V: vB, Weight: -1},
	{U: vB, V: vD, Weight: 2}, {U: vC, V: vD, Weight: 5}, {U: vD, V: vE, Weight: -2},
	{U: vE, V: vF, Weight: 3}, {U: vB, V: vF, Weight: 10}, {U: vC, V: vE, Weight: 7},
}

func TestDijkstraTraceEstudo(t *testing.T) {
	g, err := graph.BuildWeighted(6, false, dijkstraEdges)
	if err != nil {
		t.Fatal(err)
	}
	var order []int
	var rows [][]int
	r, err := DijkstraTrace(g, vA, func(s Step) {
		order = append(order, s.Vertex)
		rows = append(rows, s.Dist)
	})
	if err != nil {
		t.Fatal(err)
	}
	if want := []int{vA, vC, vB, vD, vE, vF}; !reflect.DeepEqual(order, want) {
		t.Errorf("ordem = %v, want %v", order, want)
	}
	wantRows := [][]int{
		{0, 4, 2, inf, inf, inf},
		{0, 3, 2, 10, 12, inf},
		{0, 3, 2, 8, 12, inf},
		{0, 3, 2, 8, 10, 14},
		{0, 3, 2, 8, 10, 13},
		{0, 3, 2, 8, 10, 13},
	}
	if !reflect.DeepEqual(rows, wantRows) {
		t.Errorf("tabela = %v, want %v", rows, wantRows)
	}
	if want := []int{NoPred, vC, vA, vB, vD, vE}; !reflect.DeepEqual(r.Pred, want) {
		t.Errorf("Pred = %v, want %v", r.Pred, want)
	}
	if want := []int{vA, vC, vB, vD, vE, vF}; !reflect.DeepEqual(r.PathTo(vF), want) {
		t.Errorf("PathTo(F) = %v, want %v", r.PathTo(vF), want)
	}
	if p := r.PathTo(vA); !reflect.DeepEqual(p, []int{vA}) {
		t.Errorf("PathTo(A) = %v, want [A]", p)
	}
}

func TestDijkstraHeapMatchesLinear(t *testing.T) {
	g, _ := graph.BuildWeighted(6, false, dijkstraEdges)
	lin, _ := Dijkstra(g, vA)
	hp, err := DijkstraHeap(g, vA, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(lin, hp) {
		t.Errorf("heap = %+v, linear = %+v", hp, lin)
	}
}

func TestDijkstraLighterAB(t *testing.T) {
	edges := append([]graph.WeightedEdge(nil), dijkstraEdges...)
	edges[0].Weight = 1
	g, _ := graph.BuildWeighted(6, false, edges)
	r, _ := Dijkstra(g, vA)
	if r.Dist[vF] != 11 {
		t.Errorf("Dist[F] = %d, want 11", r.Dist[vF])
	}
	if want := []int{vA, vB, vD, vE, vF}; !reflect.DeepEqual(r.PathTo(vF), want) {
		t.Errorf("PathTo(F) = %v, want %v", r.PathTo(vF), want)
	}
}

func TestBellmanFordEstudo(t *testing.T) {
	r, err := BellmanFord(6, negEdges, vA)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Passes) != 2 || !r.Passes[0].Changed || r.Passes[1].Changed {
		t.Fatalf("passadas = %+v, want 2 (a 2ª sem mudança)", r.Passes)
	}
	if want := []int{0, 1, 2, 3, 1, 4}; !reflect.DeepEqual(r.Passes[0].Dist, want) {
		t.Errorf("passada 1 = %v, want %v", r.Passes[0].Dist, want)
	}
	if want := []int{NoPred, vC, vA, vB, vD, vE}; !reflect.DeepEqual(r.Pred, want) {
		t.Errorf("Pred = %v, want %v", r.Pred, want)
	}
	if r.NegativeCycle {
		t.Error("não há ciclo negativo no grafo-base")
	}
	if want := []int{vA, vC, vB, vD, vE, vF}; !reflect.DeepEqual(r.PathTo(vF), want) {
		t.Errorf("PathTo(F) = %v, want %v", r.PathTo(vF), want)
	}
}

func TestBellmanFordNegativeCycle(t *testing.T) {
	edges := append(append([]graph.WeightedEdge(nil), negEdges...), graph.WeightedEdge{U: vE, V: vB, Weight: -5})
	r, _ := BellmanFord(6, edges, vA)
	if !r.NegativeCycle {
		t.Error("ciclo B → D → E → B (peso −5) deveria ser detectado")
	}
	for v := range 6 {
		_ = r.PathTo(v) // não pode travar com predecessores em ciclo
	}
}

func TestBellmanFordUnreachableCycle(t *testing.T) {
	// G isolado com laço G → G de peso −1: inalcançável desde A, não sinaliza.
	edges := append(append([]graph.WeightedEdge(nil), negEdges...), graph.WeightedEdge{U: 6, V: 6, Weight: -1})
	r, err := BellmanFord(7, edges, vA)
	if err != nil {
		t.Fatal(err)
	}
	if r.NegativeCycle {
		t.Error("ciclo inalcançável não deveria ser sinalizado")
	}
	if r.Dist[6] != Inf || r.PathTo(6) != nil {
		t.Errorf("G deveria ser inalcançável: dist = %d", r.Dist[6])
	}
}

func TestBellmanFordMatchesDijkstraWithoutNegatives(t *testing.T) {
	// Grafo não direcionado = cada aresta nos dois sentidos.
	var arcs []graph.WeightedEdge
	for _, e := range dijkstraEdges {
		arcs = append(arcs, e, graph.WeightedEdge{U: e.V, V: e.U, Weight: e.Weight})
	}
	g, _ := graph.BuildWeighted(6, false, dijkstraEdges)
	d, _ := Dijkstra(g, vA)
	bf, _ := BellmanFord(6, arcs, vA)
	if !reflect.DeepEqual(d.Dist, bf.Dist) {
		t.Errorf("Bellman-Ford = %v, Dijkstra = %v", bf.Dist, d.Dist)
	}
}

func TestFloydWarshallEstudo(t *testing.T) {
	dist, next, err := NewMatrices(6, negEdges)
	if err != nil {
		t.Fatal(err)
	}
	var stages []Matrix
	r := FloydWarshall(dist, next, func(_ int, d Matrix) { stages = append(stages, d) })

	wantB := Matrix{
		{0, 4, 2, 6, inf, 14},
		{inf, 0, inf, 2, inf, 10},
		{inf, -1, 0, 1, 7, 9},
		{inf, inf, inf, 0, -2, inf},
		{inf, inf, inf, inf, 0, 3},
		{inf, inf, inf, inf, inf, 0},
	}
	if !reflect.DeepEqual(stages[vB], wantB) {
		t.Errorf("D(B) = %v, want %v", stages[vB], wantB)
	}
	wantFinal := Matrix{
		{0, 1, 2, 3, 1, 4},
		{inf, 0, inf, 2, 0, 3},
		{inf, -1, 0, 1, -1, 2},
		{inf, inf, inf, 0, -2, 1},
		{inf, inf, inf, inf, 0, 3},
		{inf, inf, inf, inf, inf, 0},
	}
	if !reflect.DeepEqual(r.Dist, wantFinal) {
		t.Errorf("final = %v, want %v", r.Dist, wantFinal)
	}
	if r.NegativeCycle {
		t.Error("não há ciclo negativo no grafo-base")
	}
	if want := []int{vA, vC, vB, vD, vE, vF}; !reflect.DeepEqual(r.Path(vA, vF), want) {
		t.Errorf("Path(A, F) = %v, want %v", r.Path(vA, vF), want)
	}
	if r.Path(vF, vA) != nil {
		t.Errorf("Path(F, A) = %v, want nil", r.Path(vF, vA))
	}
	// Os custos a partir de A coincidem com os do Bellman-Ford (desafio).
	bf, _ := BellmanFord(6, negEdges, vA)
	if !reflect.DeepEqual(r.Dist[vA], bf.Dist) {
		t.Errorf("linha A = %v, Bellman-Ford = %v", r.Dist[vA], bf.Dist)
	}
}

func TestFloydWarshallNegativeCycle(t *testing.T) {
	edges := append(append([]graph.WeightedEdge(nil), negEdges...), graph.WeightedEdge{U: vE, V: vB, Weight: -5})
	dist, next, _ := NewMatrices(6, edges)
	r := FloydWarshall(dist, next, nil)
	if !r.NegativeCycle {
		t.Fatal("ciclo negativo deveria ser detectado")
	}
	for _, v := range []int{vB, vD, vE} {
		if r.Dist[v][v] >= 0 {
			t.Errorf("dist[%d][%d] = %d, want < 0", v, v, r.Dist[v][v])
		}
	}
	if !r.OnNegativeCycle(vA, vF) {
		t.Error("A → F passa pelo ciclo")
	}
	for i := range 6 {
		for j := range 6 {
			_ = r.Path(i, j) // não pode travar
		}
	}
}
