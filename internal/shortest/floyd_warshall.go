package shortest

import (
	"fmt"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/graph"
)

// Matrix é uma matriz n × n de inteiros; Inf marca "sem caminho conhecido" em
// dist e NoPred marca NULO em next.
type Matrix [][]int

func newMatrix(n, fill int) Matrix {
	m := make(Matrix, n)
	for i := range m {
		m[i] = make([]int, n)
		for j := range m[i] {
			m[i][j] = fill
		}
	}
	return m
}

// Clone copia a matriz, para registrar o estado de cada etapa.
func (m Matrix) Clone() Matrix {
	c := make(Matrix, len(m))
	for i := range m {
		c[i] = append([]int(nil), m[i]...)
	}
	return c
}

// APSP (all-pairs shortest paths) guarda as matrizes de Floyd-Warshall.
// Dist é dist[i][j]; Next é prox[i][j], o primeiro vértice depois de i numa
// rota mínima até j.
type APSP struct {
	Dist, Next    Matrix
	NegativeCycle bool
}

// NewMatrices é o CriarMatrizes do enunciado: 0 na diagonal, peso da aresta
// quando ela existe e ∞ no resto; prox[u][v] = v para cada aresta direta.
func NewMatrices(n int, edges []graph.WeightedEdge) (dist, next Matrix, err error) {
	// PASSO 1 – inicializar as matrizes.
	dist = newMatrix(n, Inf)
	next = newMatrix(n, NoPred) // ainda não conhecemos uma rota
	for i := range n {
		dist[i][i] = 0
	}
	// PASSO 2 – inserir as arestas orientadas: só dist[u][v], nunca dist[v][u].
	for _, e := range edges {
		if e.U < 0 || e.U >= n || e.V < 0 || e.V >= n {
			return nil, nil, fmt.Errorf("%w: (%d, %d) com n = %d", graph.ErrInvalidVertex, e.U, e.V, n)
		}
		if e.Weight < dist[e.U][e.V] {
			dist[e.U][e.V] = e.Weight // registrar a menor aresta direta
			next[e.U][e.V] = e.V      // o primeiro salto é o próprio destino
		}
	}
	return dist, next, nil
}

// StageFunc recebe a matriz de distâncias depois de cada intermediário k.
type StageFunc func(k int, dist Matrix)

// FloydWarshall roda o algoritmo sobre cópias de dist e next. A cada etapa k,
// dist[i][j] passa a ser o menor custo de i até j usando só os vértices
// 0..k como intermediários (programação dinâmica).
func FloydWarshall(dist, next Matrix, onStage StageFunc) APSP {
	n := len(dist)
	d, nx := dist.Clone(), next.Clone()

	// PASSO 1 – considerar cada vértice intermediário. k precisa ser o laço
	// EXTERNO: a etapa k só pode usar resultados completos da etapa k−1.
	for k := range n {
		// PASSO 2 – analisar todos os pares i → j.
		for i := range n {
			for j := range n {
				// ∞ + algo não é rota; a guarda também evita estouro de MaxInt.
				if d[i][k] == Inf || d[k][j] == Inf {
					continue
				}
				// Testar se passar por k reduz a distância (melhora estrita).
				if candidate := d[i][k] + d[k][j]; candidate < d[i][j] {
					d[i][j] = candidate
					// O primeiro passo de i até j é o primeiro passo de i até k.
					nx[i][j] = nx[i][k]
				}
			}
		}
		if onStage != nil {
			onStage(k, d.Clone())
		}
	}

	// PASSO 3 – verificar ciclos negativos: dist[v][v] começa em 0, e só fica
	// negativo se existir um ciclo saindo de v e voltando a v com soma < 0.
	negative := false
	for v := range n {
		if d[v][v] < 0 {
			negative = true
		}
	}
	return APSP{Dist: d, Next: nx, NegativeCycle: negative}
}

// Path reconstrói a rota de i até j andando pela matriz prox. Retorna nil se
// não houver rota. O limite de n passos protege contra ciclos em prox, que
// surgem quando há ciclo negativo.
func (r APSP) Path(i, j int) []int {
	// PASSO 1 – verificar se há rota.
	if i == j {
		return []int{i} // caminho trivial
	}
	if r.Next[i][j] == NoPred {
		return nil // destino inalcançável
	}
	// PASSO 2 – caminhar usando a matriz prox.
	path := []int{i}
	for cur := i; cur != j; {
		cur = r.Next[cur][j]
		if cur == NoPred || len(path) > len(r.Next) {
			return nil
		}
		path = append(path, cur)
	}
	return path
}

// OnNegativeCycle diz se a rota de i até j pode passar por um vértice v com
// dist[v][v] < 0, isto é, i alcança v e v alcança j. Nesse caso o custo de i
// até j não é bem definido (pode diminuir sem limite).
func (r APSP) OnNegativeCycle(i, j int) bool {
	for v := range r.Dist {
		if r.Dist[v][v] < 0 && r.Dist[i][v] != Inf && r.Dist[v][j] != Inf {
			return true
		}
	}
	return false
}
