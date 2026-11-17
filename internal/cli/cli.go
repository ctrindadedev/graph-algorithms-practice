// Package cli reúne utilitários de entrada usados pelos comandos em cmd/.
package cli

import (
	"bufio"
	"fmt"
	"strconv"
)

// ReadInt exibe prompt e lê o próximo inteiro de in. Tokens não numéricos
// são consumidos e geram erro, evitando que o menu entre em laço infinito.
func ReadInt(in *bufio.Reader, prompt string) (int, error) {
	fmt.Print(prompt)
	var tok string
	if _, err := fmt.Fscan(in, &tok); err != nil {
		return 0, err
	}
	n, err := strconv.Atoi(tok)
	if err != nil {
		return 0, fmt.Errorf("%q não é um número inteiro", tok)
	}
	return n, nil
}
