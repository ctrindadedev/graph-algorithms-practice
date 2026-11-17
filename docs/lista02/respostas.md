# Lista 02 – Árvores e Árvores Geradoras (Aula 9)

> Enunciado: [`enunciado.docx`](./enunciado.docx) · Código: [`cmd/lista02`](../../cmd/lista02) e [`internal/tree`](../../internal/tree)

## Pseudocódigos comentados

### EhArvore

```
Algoritmo EhArvore(G)
    // Uma árvore com |V| vértices tem exatamente |V| - 1 arestas;
    // se a contagem não bate, nem é preciso percorrer o grafo.
    se |E| ≠ |V| - 1 então retornar FALSO
    // Qualquer vértice serve como ponto de partida do teste de conectividade.
    s ← um vértice qualquer de V
    // A busca devolve todos os vértices alcançáveis a partir de s.
    visitados ← BFS(G, s)
    // Conexo + |V| - 1 arestas ⇒ acíclico ⇒ árvore.
    retornar quantidade(visitados) = |V|
```

### ArvoreGeradoraBFS

```
Algoritmo ArvoreGeradoraBFS(G, s)
    // Nenhum vértice foi descoberto ainda.
    para cada v em V: visitado[v] ← FALSO
    // Q é a fila FIFO; T guardará só as arestas da árvore.
    criar fila Q; T ← ∅
    // A raiz é marcada e inicia a busca.
    visitado[s] ← VERDADEIRO; enfileirar(Q, s)
    enquanto Q não estiver vazia:
        // Próximo vértice a explorar, na ordem de descoberta (por níveis).
        u ← desenfileirar(Q)
        para cada v em Adj[u]:
            // Só um vértice ainda não descoberto gera aresta da árvore.
            se NÃO visitado[v]:
                visitado[v] ← VERDADEIRO
                // {u, v} é a aresta que descobriu v: entra em T.
                T ← T ∪ {(u, v)}
                enfileirar(Q, v)
    // T tem as |V| - 1 arestas de descoberta (se G for conexo).
    retornar T
```

### ArvoreGeradoraDFS

```
Algoritmo ArvoreGeradoraDFS(G, s)
    para cada v em V: visitado[v] ← FALSO
    // P é a pilha LIFO; cada entrada guarda o vértice e seu pai.
    criar pilha P; T ← ∅
    // A raiz não tem pai.
    empilhar(P, (s, NULO))
    enquanto P não estiver vazia:
        (u, pai) ← desempilhar(P)
        // u pode ter sido empilhado por vários pais; só a primeira
        // retirada conta, as demais são descartadas.
        se NÃO visitado[u]:
            visitado[u] ← VERDADEIRO
            // A aresta (pai, u) entra em T no momento em que u é visitado.
            se pai ≠ NULO: T ← T ∪ {(pai, u)}
            // Ordem inversa: o menor vizinho fica no topo e sai primeiro.
            para cada v em Adj[u] em ordem inversa:
                se NÃO visitado[v]: empilhar(P, (v, u))
    retornar T
```

### ArborescenciaBFS

```
Algoritmo ArborescenciaBFS(G, r)
    // A guardará os arcos da arborescência.
    marcar todos como NÃO visitados; A ← ∅
    visitado[r] ← VERDADEIRO; enfileirar(Q, r)
    enquanto Q não estiver vazia:
        u ← desenfileirar(Q)
        // Grafo direcionado: só os arcos que SAEM de u podem ser usados.
        para cada arco (u, v):
            se NÃO visitado[v]:
                visitado[v] ← VERDADEIRO
                // O arco u → v liga o novo vértice à arborescência.
                A ← A ∪ {(u, v)}
                enfileirar(Q, v)
    // Algum vértice inalcançável a partir de r ⇒ não há arborescência com essa raiz.
    se existe vértice não visitado: FALHA
    retornar A
```

## Parte I – Conceitos

1. **Por que um grafo conexo com V − 1 arestas é uma árvore?** Um grafo conexo precisa de pelo menos V − 1 arestas para ligar todos os vértices. Se ele tivesse um ciclo, daria para remover uma aresta do ciclo e continuar conexo, ficando com V − 2 arestas, o que é impossível. Logo não há ciclo, e um grafo conexo e acíclico é uma árvore.
2. **Árvore, floresta e grafo com ciclo.** Uma *árvore* é conexa e acíclica. Uma *floresta* é acíclica mas pode ser desconexa: é um conjunto de árvores (com c componentes, tem V − c arestas). Um *grafo com ciclo* tem pelo menos um caminho fechado sem repetir arestas, logo não é árvore nem floresta.
3. **Por que a árvore geradora preserva todos os vértices, mas não todas as arestas?** "Geradora" significa que ela alcança (gera) todos os vértices do grafo original. Para ser árvore, porém, não pode ter ciclos, então as arestas que fechariam ciclos são descartadas, sobrando exatamente V − 1.
4. **Árvore geradora × árvore geradora mínima.** Qualquer subgrafo que seja árvore e contenha todos os vértices é uma árvore geradora; um grafo costuma ter várias. A árvore geradora *mínima* só faz sentido em grafos ponderados: é a árvore geradora cuja soma dos pesos das arestas é a menor possível (obtida, por exemplo, por Kruskal ou Prim).

## Atividade 1 – `EhArvore`

Saída de `go run ./cmd/lista02`, opção 1:

| Grafo                                         | \|V\| | \|E\| | É árvore? |
|-----------------------------------------------|-----|-----|-----------|
| `{(0,1), (0,2), (1,3), (1,4), (2,5)}`         | 6   | 5   | sim       |
| Grafo da Lista 01                             | 6   | 7   | não (falha na contagem) |
| `{(0,1), (1,2), (2,0), (3,4), (4,5)}`         | 6   | 5   | não (falha na conectividade) |

O terceiro caso mostra por que a contagem V − 1 sozinha não basta: o grafo tem 5 arestas, mas o triângulo 0-1-2 é um ciclo e a BFS a partir de 0 só alcança 3 dos 6 vértices.

## Atividades 2 e 3 – Árvores geradoras (grafo da Lista 01, raiz 0)

| Ordem | BFS (fila)  | DFS (pilha) |
|-------|-------------|-------------|
| 1     | `(0, 1)`    | `(0, 1)`    |
| 2     | `(0, 2)`    | `(1, 3)`    |
| 3     | `(1, 3)`    | `(3, 5)`    |
| 4     | `(1, 4)`    | `(5, 4)`    |
| 5     | `(3, 5)`    | `(4, 2)`    |

As duas têm 6 vértices e 5 = V − 1 arestas, todas pertencentes ao grafo original e sem ciclos.

## Parte V – Comparação entre BFS e DFS

| Aspecto            | BFS                    | DFS                    |
|--------------------|------------------------|------------------------|
| Estrutura auxiliar | Fila (FIFO)            | Pilha (LIFO)           |
| Altura             | 3                      | 5                      |
| Largura máxima     | 2                      | 1                      |
| Arestas exclusivas | `(0,2)`, `(1,4)`       | `(5,4)`, `(4,2)`       |
| Arestas em comum   | `(0,1)`, `(1,3)`, `(3,5)` | `(0,1)`, `(1,3)`, `(3,5)` |

- A árvore da **DFS ficou mais profunda**: é um caminho único `0 – 1 – 3 – 5 – 4 – 2`, de altura 5.
- A árvore da **BFS ficou mais larga**: tem altura 3, com dois vértices em cada nível (`{1, 2}` e `{3, 4}`).
- **Por que duas árvores diferentes podem ser válidas?** Um grafo com ciclos admite várias árvores geradoras: em cada ciclo, qualquer uma das arestas pode ser a descartada. A BFS e a DFS apenas escolhem as arestas em ordens diferentes, e as duas atendem a todos os critérios (todos os vértices, V − 1 arestas do grafo, sem ciclos).

## Atividade 4 – Arborescência por BFS

Digrafo: `0 → 1, 0 → 2, 1 → 3, 2 → 4, 3 → 5, 4 → 5, 5 → 1`.

- **Raiz 0:** arcos `0 → 1, 0 → 2, 1 → 3, 2 → 4, 3 → 5` (6 vértices, 5 arcos).
- **Raiz 3:** falha, pois `0`, `2` e `4` não são alcançáveis seguindo os arcos (de 3 só se chega a 5 e 1).

## Questões de reflexão

1. **Por que a aresta é adicionada à árvore apenas quando um novo vértice é descoberto?** Porque cada vértice (exceto a raiz) deve entrar na árvore por exatamente uma aresta, a que liga ao seu pai. Isso dá V − 1 arestas e garante que nenhuma aresta liga dois vértices que já estavam conectados na árvore.
2. **E se todas as arestas exploradas fossem adicionadas?** A "árvore" seria o próprio grafo (ou quase): arestas que levam a vértices já visitados fechariam ciclos, e o número de arestas passaria de V − 1. No grafo da Lista 01 apareceriam as 7 arestas e os ciclos `0-1-4-2-0` e `1-3-5-4-1`.
3. **Por que BFS e DFS podem produzir árvores diferentes?** Porque a ordem de exploração é diferente. A BFS (fila) termina todos os vizinhos de um vértice antes de descer; a DFS (pilha) segue o primeiro vizinho até o fim antes de voltar. Assim, um mesmo vértice pode ser descoberto por pais diferentes, como o 4: na BFS ele vem do 1, na DFS vem do 5.
4. **Por que a BFS preserva as menores distâncias a partir da raiz?** Porque a fila processa os vértices em ordem não decrescente de distância: todos os vértices a distância k saem da fila antes de qualquer vértice a distância k + 1. Assim, a primeira vez que um vértice é descoberto é por um caminho de menor número de arestas, e é essa aresta que entra na árvore.
5. **Por que V − 1 arestas não bastam para afirmar que um grafo é árvore?** Porque as arestas podem estar "mal distribuídas": um componente pode ter um ciclo (arestas sobrando) enquanto outro fica desconectado (arestas faltando). É o terceiro exemplo da Atividade 1. É preciso testar também a conectividade (ou a ausência de ciclos).
6. **Em que a arborescência difere de uma árvore geradora não direcionada?** Na arborescência os arcos têm sentido: todos saem da raiz em direção às folhas e cada vértice (exceto a raiz) tem exatamente um arco de entrada. A busca só pode seguir os arcos de saída, então a escolha da raiz importa: com raiz 3 o digrafo do exemplo não tem arborescência, embora tenha com raiz 0. Em um grafo conexo não direcionado, qualquer vértice pode ser raiz de uma árvore geradora.
