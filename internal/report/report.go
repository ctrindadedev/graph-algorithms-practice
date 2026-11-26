// Package report formata a saída dos estudos dirigidos de caminhos mínimos:
// vértices como letras (0 = A), ∞ para distâncias infinitas e – para NULO.
package report

import (
	"fmt"
	"strings"

	"github.com/ctrindadedev/graph-algorithms-practice/internal/shortest"
)

// Label converte o índice no rótulo do enunciado: 0 → A, 1 → B, ...
func Label(v int) string {
	return string(rune('A' + v))
}

// Pred mostra o predecessor (ou o próximo passo) e – para NULO.
func Pred(v int) string {
	if v == shortest.NoPred {
		return "–"
	}
	return Label(v)
}

// Dist mostra a distância e ∞ quando o vértice não foi alcançado.
func Dist(d int) string {
	if d == shortest.Inf {
		return "∞"
	}
	return fmt.Sprint(d)
}

// Path junta a rota no formato A → C → B. Rota vazia vira "(sem caminho)".
func Path(vs []int) string {
	if len(vs) == 0 {
		return "(sem caminho)"
	}
	parts := make([]string, len(vs))
	for i, v := range vs {
		parts[i] = Label(v)
	}
	return strings.Join(parts, " → ")
}

// Title escreve um título sublinhado.
func Title(s string) {
	fmt.Println()
	fmt.Println(s)
	fmt.Println(strings.Repeat("=", len([]rune(s))))
}

// Subtitle escreve um subtítulo sublinhado com hífens.
func Subtitle(s string) {
	fmt.Println()
	fmt.Println(s)
	fmt.Println(strings.Repeat("-", len([]rune(s))))
}

// pad alinha à esquerda contando runas (∞, → e – ocupam uma coluna).
func pad(s string, width int) string {
	if n := len([]rune(s)); n < width {
		return s + strings.Repeat(" ", width-n)
	}
	return s
}

// Row escreve uma linha de tabela com colunas de largura fixa.
func Row(first string, firstWidth int, cells []string, width int) {
	var b strings.Builder
	b.WriteString("  " + pad(first, firstWidth))
	for _, c := range cells {
		b.WriteString(pad(c, width))
	}
	fmt.Println(strings.TrimRight(b.String(), " "))
}

// Header escreve o cabeçalho A B C ... de uma tabela com n colunas.
func Header(first string, firstWidth, n, width int) {
	cols := make([]string, n)
	for v := range n {
		cols[v] = Label(v)
	}
	Row(first, firstWidth, cols, width)
}

// Dists converte um vetor de distâncias em células.
func Dists(ds []int) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = Dist(d)
	}
	return out
}

// Preds converte um vetor de predecessores em células.
func Preds(ps []int) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = Pred(p)
	}
	return out
}

// PrintMatrix exibe uma matriz com rótulos de linha e coluna (i \ j).
func PrintMatrix(m shortest.Matrix, cell func(int) string) {
	Header("i \\ j", 7, len(m), 5)
	for i, row := range m {
		cells := make([]string, len(row))
		for j, x := range row {
			cells[j] = cell(x)
		}
		Row(Label(i), 7, cells, 5)
	}
}
