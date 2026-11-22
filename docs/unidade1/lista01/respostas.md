# Lista 01 – Respostas teóricas

> Enunciado: [`enunciado.docx`](./enunciado.docx) · Código: [`cmd/unidade1/lista01`](../../../cmd/unidade1/lista01)

## Questão 6 – DFS recursiva × DFS com pilha

1. **Qual das duas versões parece mais simples de compreender?** A recursiva. O código é curto e segue a definição da DFS: visita o vértice e, para cada vizinho não visitado, repete o processo. O retorno para o vértice anterior acontece sozinho, quando a chamada termina. Na versão com pilha é preciso controlar a pilha, empilhar os vizinhos em ordem decrescente e descartar os vértices já visitados ao desempilhar.
2. **Qual versão utiliza explicitamente uma pilha?** A iterativa (`DFSStack`), que usa `ds.Stack[int]`.
3. **Na versão recursiva, onde a pilha está presente?** Na pilha de chamadas do programa. Cada chamada de `visit(u)` empilha um quadro com o vértice atual e a posição no laço de vizinhos. Quando a chamada termina, o quadro é desempilhado e a execução volta ao vértice anterior (o "Retornando do vértice X").
4. **Por que a ordem de inserção dos vizinhos na pilha modifica a ordem de visitação?** Porque a pilha é LIFO: o último vizinho empilhado é o primeiro a ser explorado. Empilhando `1, 2` o vértice 2 sai primeiro; empilhando `2, 1` sai o 1. Por isso os vizinhos são empilhados em ordem decrescente, para que o menor fique no topo.
5. **As duas implementações precisam produzir a mesma sequência?** Não necessariamente. As duas são DFS válidas, mas a ordem depende de como os vizinhos são escolhidos. A recursiva escolhe o próximo vizinho somente depois de voltar da chamada anterior, enquanto a iterativa empilha todos os vizinhos de uma vez. Com os vizinhos em ordem crescente e empilhados em ordem decrescente, os resultados coincidem neste grafo (`0 → 1 → 3 → 5 → 4 → 2`). Empilhando na ordem crescente, a iterativa visitaria `0 → 2 → 4 → 5 → 3 → 1`.

## Questão 10 – Novo grafo (aresta 2–5)

Execução: `go run ./cmd/unidade1/lista01` (original) e `go run ./cmd/unidade1/lista01 -extra` (com a aresta (2,5)).

**Matriz de adjacência.** Passam a valer 1 as posições `[2][5]` e `[5][2]`. A matriz continua simétrica.

```
Original           Com (2,5)
0 1 1 0 0 0        0 1 1 0 0 0
1 0 0 1 1 0        1 0 0 1 1 0
1 0 0 0 1 0        1 0 0 0 1 1
0 1 0 0 0 1        0 1 0 0 0 1
0 1 1 0 0 1        0 1 1 0 0 1
0 0 0 1 1 0        0 0 1 1 1 0
```

**Lista de adjacência.** Mudam apenas os vértices 2 e 5. O vértice 5 fica com `[2, 3, 4]` (e não `[3, 4, 2]`) porque os vizinhos são mantidos em ordem crescente.

```
Original           Com (2,5)
2 → [0, 4]         2 → [0, 4, 5]
5 → [3, 4]         5 → [2, 3, 4]
```

**Percursos a partir do vértice 0.**

| Algoritmo     | Original                    | Com (2,5)                   |
|---------------|-----------------------------|-----------------------------|
| BFS           | `0 → 1 → 2 → 3 → 4 → 5`     | `0 → 1 → 2 → 3 → 4 → 5`     |
| DFS recursiva | `0 → 1 → 3 → 5 → 4 → 2`     | `0 → 1 → 3 → 5 → 2 → 4`     |
| DFS com pilha | `0 → 1 → 3 → 5 → 4 → 2`     | `0 → 1 → 3 → 5 → 2 → 4`     |

**Alterações observadas.**

- **BFS:** a ordem não muda, pois os níveis continuam os mesmos (nível 1: 1, 2; nível 2: 3, 4; nível 3: 5). O que muda é *quem descobre* o vértice 5: antes ele entrava na fila ao processar o 3 (fila `[4, 5]`), agora entra já ao processar o 2 (fila `[3, 4, 5]`). A distância de 0 até 5 cai de 3 para 2 arestas.
- **DFS (as duas versões):** ao chegar no vértice 5, o menor vizinho não visitado passa a ser o 2 em vez do 4. A busca então segue `5 → 2 → 4` em vez de `5 → 4 → 2`.
- **DFS com pilha:** ao processar o 5 são empilhados `4` e `2`, e a pilha cresce para `[2, 4, 4, 2]`. Há mais entradas repetidas, que são descartadas no final (três retiradas ignoradas, contra duas no grafo original).

## Questão 11 – Matriz × Lista de Adjacência

1. **Qual representação pareceu mais simples de implementar?** A matriz. Inserir e consultar uma aresta é só acessar `cells[u][v]`, e os vizinhos já saem em ordem crescente ao percorrer a linha. Na lista foi preciso cuidar da inserção ordenada e de arestas repetidas.
2. **Quantas posições possui a matriz de um grafo com seis vértices?** 6 × 6 = **36** posições, das quais só 14 (2 × 7 arestas) valem 1 no grafo da lista.
3. **Por que a matriz de um grafo não direcionado é simétrica?** Porque a aresta {u, v} liga os dois vértices nos dois sentidos. Toda inserção grava `matriz[u][v] = 1` e `matriz[v][u] = 1`, então `matriz[i][j] = matriz[j][i]` para todo par.
4. **Qual estrutura tende a utilizar menos memória em grafos esparsos?** A lista de adjacência, que usa O(V + E) de memória, contra O(V²) da matriz. Quando E é bem menor que V², a lista guarda muito menos.
5. **A lista armazena posições correspondentes a arestas inexistentes?** Não. Cada lista contém apenas os vizinhos reais do vértice; a ausência de aresta é representada simplesmente por o vértice não aparecer na lista. Já a matriz reserva uma posição (com 0) para cada par, exista a aresta ou não.

## Questão 12 – BFS × DFS

```
BFS → FILA  → FIFO → percurso por NÍVEIS
DFS → PILHA → LIFO → percurso em PROFUNDIDADE
```

1. **Qual estrutura é utilizada pela BFS?** A fila.
2. **Qual estrutura é utilizada pela DFS iterativa?** A pilha explícita.
3. **Qual estrutura está implicitamente envolvida na DFS recursiva?** A pilha de chamadas (pilha de execução) do programa.
4. **Qual algoritmo percorre o grafo por níveis?** A BFS.
5. **Qual algoritmo procura aprofundar um caminho antes de retornar?** A DFS.
