# Lista 03 – BFS e DFS em grafos direcionados e ponderados

> Enunciado: [`enunciado.docx`](./enunciado.docx) · Código: [`cmd/unidade1/lista03`](../../../cmd/unidade1/lista03), [`internal/graph`](../../../internal/graph) (`digraph.go`, `weighted.go`), [`internal/traversal`](../../../internal/traversal) e [`internal/shortest`](../../../internal/shortest)
>
> Linguagem: **Go**. Os vértices A..E viram os índices 0..4, então a ordem crescente dos índices é a ordem alfabética pedida no enunciado. Origem: A.

## Construção e buscas no código

- **Estrutura do grafo.** Cada tipo é uma lista de adjacência: `Digraph` guarda só o destino `v` em `out[u]`; `Weighted` guarda pares `Arc{To, Weight}` em `adj[u]` e tem um campo `directed`. As listas ficam sempre ordenadas pelo vizinho, então percorrê-las já segue a ordem alfabética.
- **Inserção.** `Digraph.AddArc(u, v)` grava apenas `v` em `out[u]`. `Weighted.AddEdge(u, v, p)` grava `(v, p)` em `adj[u]` e, se o grafo não for direcionado, também `(u, p)` em `adj[v]`. Laços e vértices fora do intervalo são rejeitados.
- **Buscas.** `BFS`, `DFSRecursive` e `DFSStack` recebem um `graph.Neighborhood` (`Order` + `Neighbors`), então os mesmos algoritmos da Lista 01 servem para os três grafos. `Weighted.Neighbors` descarta os pesos, que é exatamente o `v ← entrada.vertice` do pseudocódigo; em digrafos, `Neighbors` só devolve os arcos de saída.
- **Dijkstra.** `shortest.Dijkstra` é a versão O(V²): a cada rodada fecha o vértice aberto de menor distância e relaxa os arcos que saem dele. Rejeita pesos negativos.

## Etapa 1 – Grafo direcionado

Arcos: `A→B, A→C, B→D, C→D, D→E`.

```
A → [B, C]
B → [D]
C → [D]
D → [E]
E → []
```

**Entradas por aresta:** cada arco `u → v` gera **uma** entrada, `v` na lista de `u`. Não existe a entrada inversa, então `D` não “vê” `B` nem `C`. Total: 5 entradas.

**BFS desde A**

| Processando | Enfileirados | Fila depois |
|-------------|--------------|-------------|
| A | B, C | [B, C] |
| B | D | [C, D] |
| C | – (D já marcado) | [D] |
| D | E | [E] |
| E | – | [] |

Ordem: **A, B, C, D, E**.

**DFS desde A:** entra em A → B → D → E, retorna até A e então entra em C (o único vizinho de C, D, já foi visitado). Ordem: **A, B, D, E, C**.

## Etapa 2 – Grafo ponderado

Arestas: `A—B(4), A—C(2), B—C(1), B—D(5), C—D(8), C—E(10), D—E(2)`.

```
A → [(B, 4), (C, 2)]
B → [(A, 4), (C, 1), (D, 5)]
C → [(A, 2), (B, 1), (D, 8), (E, 10)]
D → [(B, 5), (C, 8), (E, 2)]
E → [(C, 10), (D, 2)]
```

**Entradas por aresta:** cada aresta gera **duas** entradas, `(v, p)` em `u` e `(u, p)` em `v`, com o mesmo peso. Total: 14 entradas. O peso só aumenta o tamanho de cada entrada; não muda quais vizinhos existem.

**BFS desde A**

| Processando | Enfileirados | Fila depois |
|-------------|--------------|-------------|
| A | B, C | [B, C] |
| B | D (A e C já marcados) | [C, D] |
| C | E | [D, E] |
| D | – | [E] |
| E | – | [] |

Ordem: **A, B, C, D, E**.

**DFS desde A:** A → B → C (primeiro vizinho não visitado de B) → D (A e B já visitados) → E, depois retorna até A. Ordem: **A, B, C, D, E**.

## Etapa 3 – Grafo direcionado e ponderado

Arcos: `A→B(4), A→C(2), C→B(1), B→D(5), C→D(8), C→E(10), D→E(2)`.

```
A → [(B, 4), (C, 2)]
B → [(D, 5)]
C → [(B, 1), (D, 8), (E, 10)]
D → [(E, 2)]
E → []
```

**Entradas por aresta:** cada arco gera **uma** entrada, `(v, p)` na lista de `u`. Total: 7 entradas. `C → B` com peso 1 não implica `B → C`, por isso a lista de B só tem D.

**BFS desde A**

| Processando | Enfileirados | Fila depois |
|-------------|--------------|-------------|
| A | B, C | [B, C] |
| B | D | [C, D] |
| C | E (B e D já marcados) | [D, E] |
| D | – | [E] |
| E | – | [] |

Ordem: **A, B, C, D, E**.

**DFS desde A:** A → B → D → E; de B não há arco para C, então a busca retorna até A e só então entra em C. Ordem: **A, B, D, E, C**.

## Comparação dos percursos

Saída de `go run ./cmd/unidade1/lista03`, opção 4:

| Tipo de grafo | Ordem BFS a partir de A | Ordem DFS a partir de A |
|---------------|-------------------------|-------------------------|
| Direcionado | A, B, C, D, E | A, B, D, E, C |
| Ponderado | A, B, C, D, E | A, B, C, D, E |
| Direcionado e ponderado | A, B, C, D, E | A, B, D, E, C |

## Perguntas para análise

1. **BFS e DFS visitam necessariamente os vértices na mesma ordem?** Não. A BFS (fila) termina todos os vizinhos de um vértice antes de descer um nível; a DFS (pilha/recursão) segue o primeiro vizinho até o fim antes de voltar. No grafo direcionado, a BFS visita C antes de D, mas a DFS desce A → B → D → E e só visita C no final. As ordens podem coincidir por acaso, como no grafo ponderado, em que as duas dão A, B, C, D, E.
2. **O que muda quando uma aresta passa de não direcionada para direcionada?** A aresta deixa de existir em um dos sentidos, então alguns vértices podem ficar inalcançáveis e a ordem de visita muda. Comparando as etapas 2 e 3: com `B—C`, a DFS vai A → B → C; com `C → B`, B não alcança C e a DFS vai A → B → D → E, deixando C para o fim. Em casos extremos a busca encolhe: a partir de D no digrafo só se alcança E.
3. **Por que o peso não ordena a fila da BFS nem as chamadas da DFS?** Porque BFS e DFS respondem perguntas de conectividade (o que é alcançável e em que ordem de exploração), e a “distância” da BFS é contada em número de arestas. Ordenar a fronteira pelo custo acumulado já seria outro algoritmo, com fila de prioridade: Dijkstra. Além disso, o enunciado fixa a ordem alfabética para que os percursos sejam reproduzíveis.
4. **No grafo ponderado, o caminho com menos arestas é sempre o de menor custo?** Não. De A até B, a aresta direta custa 4, mas A – C – B custa 2 + 1 = 3. De A até E, o caminho com menos arestas é A – C – E (2 arestas, custo 12), mas A – C – B – D – E custa 2 + 1 + 5 + 2 = 10 com 4 arestas.
5. **Como ler do teclado os vértices, as arestas e os pesos?** Ler a quantidade de vértices (ou seus nomes, mapeando cada letra para um índice), se o grafo é direcionado e se é ponderado, a quantidade de arestas e, em seguida, uma linha `u v` ou `u v peso` por aresta. Cada linha é validada (vértices existentes, sem laço, peso inteiro) e repassada para `AddArc`/`AddEdge`. A Lista 01 já faz isso para grafos simples em `readConfig` (`go run ./cmd/unidade1/lista01 -input`); bastaria acrescentar as duas perguntas de tipo e o peso.
6. **Como tratar vértices não alcançáveis a partir da origem?** Eles simplesmente nunca entram na fila/pilha, então ficam com `visitado = FALSO` ao final. O programa compara os visitados com V e lista os que faltaram (`go run ./cmd/unidade1/lista03 -start D` mostra A, B e C como inalcançáveis). Se for preciso percorrer o grafo todo, basta reiniciar a busca a partir de cada vértice ainda não visitado, em ordem alfabética, formando uma floresta de busca. No Dijkstra, esses vértices ficam com distância ∞ e sem predecessor.

## Desafio – Dijkstra no grafo direcionado e ponderado

Rastreamento a partir de A (a cada rodada fecha-se o vértice aberto de menor distância):

| Rodada | Fecha | A | B | C | D | E |
|--------|-------|---|---|---|---|---|
| início | – | 0 | ∞ | ∞ | ∞ | ∞ |
| 1 | A (0) | **0** | 4 (A) | 2 (A) | ∞ | ∞ |
| 2 | C (2) | | 3 (C) | **2** | 10 (C) | 12 (C) |
| 3 | B (3) | | **3** | | 8 (B) | 12 (C) |
| 4 | D (8) | | | | **8** | 10 (D) |
| 5 | E (10) | | | | | **10** |

Vetores finais (opção 5 do programa):

| Vértice | dist | pred |
|---------|------|------|
| A | 0 | – |
| B | 3 | C |
| C | 2 | A |
| D | 8 | B |
| E | 10 | D |

- **Dijkstra:** A → C → B → D → E, 4 arcos, custo **10**.
- **Árvore de busca BFS:** E é descoberto a partir de C, então o caminho é A → C → E, 2 arcos, custo **12**.

A BFS acha o caminho com menos arcos; o Dijkstra acha o de menor custo, que aqui usa o dobro de arcos. Os pesos são todos não negativos, condição necessária para o Dijkstra ser correto.

## Extensão – DFS iterativa com pilha

Para reproduzir a ordem da versão recursiva, os vizinhos de cada vértice são empilhados em ordem alfabética **inversa**, deixando o menor no topo; um vértice desempilhado que já foi visitado é descartado. Ordem de empilhamento no grafo direcionado (opção 6):

| Desempilha | Empilha | Pilha (topo primeiro) |
|------------|---------|-----------------------|
| A | C, B | [B, C] |
| B | D | [D, C] |
| D | E | [E, C] |
| E | – | [C] |
| C | – | [] |

Ordem: A, B, D, E, C, igual à recursiva. Nos três grafos a ordem iterativa coincide com a recursiva; no ponderado aparecem entradas repetidas (E, D e C são empilhados duas vezes) que são descartadas ao sair da pilha.
