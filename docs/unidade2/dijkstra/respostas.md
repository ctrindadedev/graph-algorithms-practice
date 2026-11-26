# Estudo dirigido – Algoritmo de Dijkstra

> Enunciado: [`enunciado.docx`](./enunciado.docx) · Código: [`cmd/unidade2/dijkstra`](../../../cmd/unidade2/dijkstra), [`internal/graph/weighted.go`](../../../internal/graph/weighted.go) e [`internal/shortest/dijkstra.go`](../../../internal/shortest/dijkstra.go)
>
> Linguagem: **Go** (o enunciado sugere Python; a lógica dos pseudocódigos foi preservada). Os vértices A..F viram os índices 0..5, então a ordem crescente dos índices é a ordem alfabética e o desempate “menor letra” é automático.

```bash
go run ./cmd/unidade2/dijkstra   # executa todas as partes, incluindo testes e desafio
```

## Questão 1 – Grafo em lista de adjacência

A tradução foi feita em **Go** (a lógica e os passos do pseudocódigo foram preservados). Os vértices A..F viram os índices 0..5, então a ordem crescente dos índices é a ordem alfabética pedida no enunciado. O grafo fica no tipo `graph.Weighted` (`internal/graph/weighted.go`): `adj[u]` é a lista de pares `Arc{To, Weight}`, o equivalente a `lista[u]` com tuplas (vizinho, peso).

**1. Lista das 9 arestas**

```go
// arestas do grafo-base (Figura 1), não direcionado.
var arestas = []graph.WeightedEdge{
    {U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: B, V: C, Weight: 1},
    {U: B, V: D, Weight: 5}, {U: C, V: D, Weight: 8}, {U: C, V: E, Weight: 10},
    {U: D, V: E, Weight: 2}, {U: D, V: F, Weight: 6}, {U: E, V: F, Weight: 3},
}
```

**2. criar_grafo(vertices, arestas)**

`graph.BuildWeighted(n, directed, arestas)` faz o papel de `criar_grafo`: `NewWeighted` cria uma lista vazia por vértice (PASSO 1) e `AddEdge` registra cada aresta (PASSO 2). Como o grafo não é direcionado (`directed = false`), `(v, peso)` entra em `lista[u]` e `(u, peso)` entra em `lista[v]`.

```go
func NewWeighted(n int, directed bool) *Weighted {
    // PASSO 1 – criação: uma lista vazia de pares (vizinho, peso) por vértice.
    return &Weighted{directed: directed, adj: make([][]Arc, n)}
}

func BuildWeighted(n int, directed bool, edges []WeightedEdge) (*Weighted, error) {
    g := NewWeighted(n, directed)
    for _, e := range edges {
        if err := g.AddEdge(e.U, e.V, e.Weight); err != nil {
            return nil, fmt.Errorf("aresta %d-%d: %w", e.U, e.V, err)
        }
    }
    return g, nil
}

func (g *Weighted) AddEdge(u, v, w int) error {
    if err := validate(g.Order(), u, v); err != nil {
        return err
    }
    // PASSO 2 – inserção: (v, peso) entra na lista de u. Sem direção, (u, peso)
    // entra também na lista de v; com direção, não há arco reverso automático.
    g.adj[u] = insertArc(g.adj[u], Arc{To: v, Weight: w})
    if !g.directed {
        g.adj[v] = insertArc(g.adj[v], Arc{To: u, Weight: w})
    }
    return nil
}
```

**3. Listas de A a F (PASSO 3)**

```text
  A → [(B, 4), (C, 2)]
  B → [(A, 4), (C, 1), (D, 5)]
  C → [(A, 2), (B, 1), (D, 8), (E, 10)]
  D → [(B, 5), (C, 8), (E, 2), (F, 6)]
  E → [(C, 10), (D, 2), (F, 3)]
  F → [(D, 6), (E, 3)]
```

As conexões são bidirecionais: cada uma das 9 arestas aparece duas vezes, uma em cada ponta e com o mesmo peso, num total de 18 entradas. Por exemplo, `A → (B, 4)` e `B → (A, 4)`, ou `E → (F, 3)` e `F → (E, 3)`. As listas ficam ordenadas pelo vizinho, então percorrê-las já segue a ordem alfabética.

**4. E se o grafo fosse direcionado?**

Seria registrada só a ligação de ida: `(v, peso)` em `lista[u]`, sem o retorno `(u, peso)` em `lista[v]`. No código, basta chamar `BuildWeighted` com `directed = true`, e o `if !g.directed` de `AddEdge` deixa de inserir o arco reverso. Cada aresta passa a gerar uma entrada só (9 no total), e `F` fica com a lista vazia. Essa versão é a do desafio da Parte VII.

## Questão 2 – Dijkstra com seleção linear

A primeira versão, com seleção linear e sem heap, é `shortest.DijkstraTrace` (`internal/shortest/dijkstra.go`). Os nomes do pseudocódigo aparecem nos comentários: `dist` = `distancia`, `pred` = `anterior`, `closed` = `visitado`, `candidate` = `candidato`, `Inf` (`math.MaxInt`) = infinito e `NoPred` (−1) = NULO. O “Escrever distancia, anterior” do fim de cada iteração é a função `onStep`, que recebe uma cópia dos dois vetores; o programa principal a imprime com `printStep`.

```go
// DijkstraTrace é o Dijkstra com seleção linear, avisando onStep ao fim de cada
// iteração (Questão 2 do estudo dirigido). Em empate de distância, escolhe o
// vértice de menor índice, que é a menor letra.
func DijkstraTrace(g *graph.Weighted, src int, onStep StepFunc) (Result, error) {
    if err := checkNonNegative(g); err != nil {
        return Result{}, err
    }
    n := g.Order()

    // PASSO 1 – inicializar as estruturas de controle:
    // dist = distancia[v] (∞: ainda não conhecemos rota até v),
    // pred = anterior[v] (NULO: não há predecessor registrado),
    // closed = visitado[v] (falso: a distância de v ainda não está definida).
    dist, pred := initVectors(n, src)
    closed := make([]bool, n)

    // PASSO 2 – repetir a escolha do menor custo enquanto existir vértice não
    // visitado com distância finita.
    for {
        // u ← não visitado com menor distancia[u]. Percorrer v em ordem crescente
        // e trocar só com "<" garante o desempate pela menor letra.
        u := -1
        for v := range n {
            if !closed[v] && dist[v] != Inf && (u == -1 || dist[v] < dist[u]) {
                u = v
            }
        }
        if u == -1 {
            break // os vértices restantes são inalcançáveis
        }
        // Com pesos não negativos, nenhum caminho futuro consegue baratear u.
        closed[u] = true

        // PASSO 3 – examinar as arestas que saem de u (relaxamento).
        for _, a := range g.Arcs(u) {
            if closed[a.To] {
                continue
            }
            candidate := dist[u] + a.Weight
            // Só substitui com melhora estrita: em empate, mantém o predecessor antigo.
            if candidate < dist[a.To] {
                dist[a.To] = candidate
                pred[a.To] = u // guardar u para reconstruir o caminho até v
            }
        }
        if onStep != nil {
            onStep(Step{Vertex: u, Dist: slices.Clone(dist), Pred: slices.Clone(pred)})
        }
    }

    // PASSO 4 – devolver os resultados.
    return Result{Source: src, Dist: dist, Pred: pred}, nil
}

// initVectors devolve distancia (∞ em todos, 0 na origem) e anterior (NULO).
func initVectors(n, src int) (dist, pred []int) {
    dist = make([]int, n)
    pred = make([]int, n)
    for v := range n {
        dist[v], pred[v] = Inf, NoPred
    }
    dist[src] = 0 // custo para chegar à origem = zero
    return dist, pred
}
```

**Seleção do próximo vértice**

O enunciado sugere `min()` com chave (distância, rótulo). O laço que escolhe `u` faz o mesmo: percorre os vértices em ordem crescente e só troca o candidato quando acha uma distância **estritamente** menor. Assim, num empate, fica o primeiro encontrado, que é a menor letra. Vértices com distância ∞ são ignorados; quando nenhum aberto tem distância finita, `u` continua −1 e o laço termina. Essa é a condição “enquanto existir vértice não visitado com distância finita”.

**Relaxamento**

Para cada par `(v, peso)` de `lista[u]` ainda não visitado, calcula `candidato = distancia[u] + peso` e só atualiza `distancia[v]` e `anterior[v]` se `candidato < distancia[v]` (melhora estrita).

**Saída de cada iteração**

```text
  Processando A (distância 0)
    distancia: {A: 0, B: 4, C: 2, D: ∞, E: ∞, F: ∞}
    anterior:  {A: –, B: A, C: A, D: –, E: –, F: –}
  Processando C (distância 2)
    distancia: {A: 0, B: 3, C: 2, D: 10, E: 12, F: ∞}
    anterior:  {A: –, B: C, C: A, D: C, E: C, F: –}
  Processando B (distância 3)
    distancia: {A: 0, B: 3, C: 2, D: 8, E: 12, F: ∞}
    anterior:  {A: –, B: C, C: A, D: B, E: C, F: –}
  Processando D (distância 8)
    distancia: {A: 0, B: 3, C: 2, D: 8, E: 10, F: 14}
    anterior:  {A: –, B: C, C: A, D: B, E: D, F: D}
  Processando E (distância 10)
    distancia: {A: 0, B: 3, C: 2, D: 8, E: 10, F: 13}
    anterior:  {A: –, B: C, C: A, D: B, E: D, F: E}
  Processando F (distância 13)
    distancia: {A: 0, B: 3, C: 2, D: 8, E: 10, F: 13}
    anterior:  {A: –, B: C, C: A, D: B, E: D, F: E}
```

## Questão 3 – Acompanhamento passo a passo

| Etapa | Vértice escolhido | d(A) | d(B) | d(C) | d(D) | d(E) | d(F) |
|---|---|---|---|---|---|---|---|
| Inicial | – | 0 | ∞ | ∞ | ∞ | ∞ | ∞ |
| 1 | A | 0 | 4 | 2 | ∞ | ∞ | ∞ |
| 2 | C | 0 | 3 | 2 | 10 | 12 | ∞ |
| 3 | B | 0 | 3 | 2 | 8 | 12 | ∞ |
| 4 | D | 0 | 3 | 2 | 8 | 10 | 14 |
| 5 | E | 0 | 3 | 2 | 8 | 10 | 13 |
| 6 | F | 0 | 3 | 2 | 8 | 10 | 13 |

| Momento | anterior[B] | anterior[C] | anterior[D] | anterior[E] | anterior[F] |
|---|---|---|---|---|---|
| Após A | A | A | – | – | – |
| Após C | C | A | C | C | – |
| Final | C | A | B | D | E |

As tabelas do enunciado foram preenchidas acima. Os valores conferem com a saída do programa. Ordem de processamento: **A, C, B, D, E, F**.

**Por que a distância de B muda de 4 para 3?**

Na iteração 1, B recebe 4 pela aresta direta A—B. Na iteração 2, C é fechado com distância 2, e o relaxamento de C—B(1) dá `candidato = 2 + 1 = 3`, que é menor que 4. Então `distancia[B] = 3` e `anterior[B] = C`: a rota A → C → B, com duas arestas, é mais barata que a aresta direta A—B.

**Em qual iteração a distância provisória de F é melhorada?**

F deixa de ser ∞ na **iteração 4**, ao processar D: 8 + 6 = 14, com anterior[F] = D. Na **iteração 5**, ao processar E, F melhora de 14 para 10 + 3 = 13, com anterior[F] = E. Na iteração 6, F é fechado com 13.

**Uma distância já definida chega a ser alterada?**

Não. O código nem tenta: o relaxamento pula vértices com `closed[v] = true`. E a tentativa nunca melhoraria, por causa dos pesos não negativos. Quando u é escolhido, `distancia[u]` é a menor entre os abertos. Qualquer outra rota até u precisa sair do conjunto fechado por algum vértice aberto x, com `distancia[x] ≥ distancia[u]`, e o resto do trajeto só soma pesos ≥ 0, então não fica abaixo de `distancia[u]`. No rastreamento, o valor de cada vértice no momento em que é escolhido (A=0, C=2, B=3, D=8, E=10, F=13) é o valor final.

## Questão 4 – Reconstrução do caminho mínimo

`reconstruir_caminho` é o método `Result.PathTo` (`internal/shortest/dijkstra.go`). Ele segue o pseudocódigo: devolve lista vazia (`nil`) se o destino tiver distância ∞ (PASSO 1), segue os predecessores do destino até a origem com `append` (PASSO 2) e inverte a lista no final, como o `append()` + `reverse()` sugerido (PASSO 3). O laço para ao chegar à origem e dá no máximo n passos, uma proteção que também é usada no Bellman-Ford.

```go
// PathTo reconstrói o caminho mínimo seguindo os predecessores de v até a
// origem (ReconstruirCaminho). Retorna nil se v for inalcançável ou se a cadeia
// de predecessores não chegar à origem, o que pode acontecer quando há ciclo
// negativo. O laço dá no máximo n passos, então nunca fica preso em um ciclo.
func (r Result) PathTo(v int) []int {
    // PASSO 1 – verificar se o destino é alcançável.
    if r.Dist[v] == Inf {
        return nil
    }
    // PASSO 2 – percorrer os predecessores de trás para frente.
    var path []int
    for x, steps := v, 0; x != NoPred && steps <= len(r.Pred); x, steps = r.Pred[x], steps+1 {
        path = append(path, x)
        if x == r.Source {
            // PASSO 3 – inverter para entregar a rota da origem ao destino.
            slices.Reverse(path)
            return path
        }
    }
    return nil // a cadeia não chegou à origem
}
```

**Rota até F e demais destinos**

```text
  A → F: A → C → B → D → E → F      custo 13, 5 aresta(s)
  A → B: A → C → B                  custo 3, 2 aresta(s)
  A → D: A → C → B → D              custo 8, 3 aresta(s)
  A → E: A → C → B → D → E          custo 10, 4 aresta(s)
```

**A → C → B → D → E → F, custo 13** (2 + 1 + 5 + 2 + 3).

**Custo × número de arestas**

O custo é a soma dos pesos das arestas da rota; o número de arestas é só o comprimento dela. Dijkstra minimiza o custo. A rota mínima até F tem 5 arestas e custa 13, enquanto A → B → D → F tem só 3 arestas e custa 4 + 5 + 6 = 15. Uma rota com menos arestas pode sair mais cara.

## Questão 5 – Estrutura do programa

Cada TODO do esqueleto virou uma função Go reaproveitável, e o programa principal fica em `cmd/unidade2/dijkstra/main.go`:

| Esqueleto (Python) | Implementação (Go) | Arquivo |
|---|---|---|
| criar_grafo(vertices, arestas) | `graph.BuildWeighted(n, false, arestas)` | internal/graph/weighted.go |
| dijkstra(grafo, origem) | `shortest.DijkstraTrace(g, origem, onStep)` | internal/shortest/dijkstra.go |
| reconstruir_caminho(...) | `Result.PathTo(destino)` | internal/shortest/dijkstra.go |
| programa principal | `run()`: monta o grafo, roda Dijkstra de A e mostra A → F | cmd/unidade2/dijkstra/main.go |

```go
func run() error {
    // PARTE I – criar_grafo: lista de adjacência não direcionada.
    grafo, err := graph.BuildWeighted(F+1, false, arestas)
    if err != nil {
        return err
    }
    r.Title("PARTE I – Lista de adjacência (não direcionada)")
    printAdjacency(grafo)

    // PARTES II e III – Dijkstra com seleção linear e rastreamento.
    r.Title("PARTES II e III – Dijkstra a partir de A (seleção linear)")
    res, err := shortest.DijkstraTrace(grafo, A, printStep)
    if err != nil {
        return err
    }
    printFinal(res)

    // PARTE IV – reconstruir_caminho.
    r.Title("PARTE IV – Reconstrução dos caminhos mínimos")
    for _, dest := range []int{F, B, D, E} {
        printRoute(res, dest)
    }

    // PARTE VI – testes da Questão 6.
    r.Title("PARTE VI – Testes")
    if err := runTests(res); err != nil {
        return err
    }
    if err := compareBFS(grafo, res); err != nil {
        return err
    }

    // PARTE VII – desafio: digrafo e heap.
    return challenge()
}
```

Para executar: `go run ./cmd/unidade2/dijkstra` (ou `make run LISTA=dijkstra`). Os testes automatizados com os valores da tabela ficam em `internal/shortest/unidade2_test.go` (`go test ./...`). Nenhuma biblioteca de caminho mínimo foi usada: só a biblioteca padrão de Go.

## Questão 6 – Testes e análise

| Teste | Rota | Custo | Processados | Justificativa |
|---|---|---|---|---|
| 1 – Original; A → F | A → C → B → D → E → F | 13 | A, C, B, D, E, F | Rota do rastreamento: a via C → B barateia B, e E → F barateia F. |
| 2 – Original; A → A | [A] | 0 | A, C, B, D, E, F | distancia[A] = 0 e anterior[A] = NULO: a reconstrução para na própria origem. |
| 3 – Original; A → D | A → C → B → D | 8 | A, C, B, D, E, F | anterior[D] = B, anterior[B] = C, anterior[C] = A. D recebe 10 via C e cai para 8 via B. |
| 4 – G isolado; A → G | [] (vazio) | ∞ | A, C, B, D, E, F | G nunca recebe distância finita, então nunca é escolhido; PathTo devolve lista vazia. |
| 5 – A—B = 1; A → F | A → B → D → E → F | 11 | A, B, C, D, E, F | A rota muda: B passa a custar 1 pela aresta direta, e C deixa de fazer parte do caminho. |

**Saída do programa nos testes**

```text
Teste 1 – Original; A → F
-------------------------
  A → F: A → C → B → D → E → F      custo 13, 5 aresta(s)

Teste 2 – Original; A → A
-------------------------
  A → A: A                          custo 0, 0 aresta(s)

Teste 3 – Original; A → D
-------------------------
  A → D: A → C → B → D              custo 8, 3 aresta(s)
  Predecessores intermediários: anterior[D] = B, anterior[B] = C, anterior[C] = A

Teste 4 – Vértice G isolado; A → G
----------------------------------
  Processados: A C B D E F
  A → G: inalcançável (distância ∞, caminho [])

Teste 5 – Peso de A—B trocado de 4 para 1; A → F
------------------------------------------------
  Processando A (distância 0)
    distancia: {A: 0, B: 1, C: 2, D: ∞, E: ∞, F: ∞}
    anterior:  {A: –, B: A, C: A, D: –, E: –, F: –}
  Processando B (distância 1)
    distancia: {A: 0, B: 1, C: 2, D: 6, E: ∞, F: ∞}
    anterior:  {A: –, B: A, C: A, D: B, E: –, F: –}
  Processando C (distância 2)
    distancia: {A: 0, B: 1, C: 2, D: 6, E: 12, F: ∞}
    anterior:  {A: –, B: A, C: A, D: B, E: C, F: –}
  Processando D (distância 6)
    distancia: {A: 0, B: 1, C: 2, D: 6, E: 8, F: 12}
    anterior:  {A: –, B: A, C: A, D: B, E: D, F: D}
  Processando E (distância 8)
    distancia: {A: 0, B: 1, C: 2, D: 6, E: 8, F: 11}
    anterior:  {A: –, B: A, C: A, D: B, E: D, F: E}
  Processando F (distância 11)
    distancia: {A: 0, B: 1, C: 2, D: 6, E: 8, F: 11}
    anterior:  {A: –, B: A, C: A, D: B, E: D, F: E}
  A → F: A → B → D → E → F          custo 11, 4 aresta(s)
```

**BFS × Dijkstra**

```text
  BFS:      A → B → D → F              3 aresta(s), custo 15
  Dijkstra: A → C → B → D → E → F      5 aresta(s), custo 13
```

A BFS minimiza o **número de arestas** e ignora os pesos: chega a F por A → B → D → F (3 arestas, custo 15). Dijkstra minimiza a **soma dos pesos** e encontra A → C → B → D → E → F (5 arestas, custo 13). As duas respostas só coincidem quando todos os pesos são iguais.

**Por que não usar Dijkstra com pesos negativos?**

O algoritmo é guloso: ao fechar u, supõe que nenhum caminho futuro consegue baratear `distancia[u]`, e isso só vale se as arestas restantes não subtraem custo. Com uma aresta negativa, um vértice já fechado poderia melhorar depois, e o resultado sai errado. Exemplo: A→B(2), A→C(3), C→B(−2). Dijkstra fecha B com 2, mas o custo real é 3 − 2 = 1. Com um ciclo negativo, nem existe custo mínimo. Por isso `Dijkstra` recusa o grafo e devolve `ErrNegativeWeight`; para pesos negativos, o certo é Bellman-Ford.

**Complexidade**

Na versão com seleção linear, o laço principal roda no máximo |V| vezes, e cada escolha varre os |V| vértices: O(|V|²). Somando todas as iterações, os relaxamentos percorrem cada lista uma vez: O(|E|). Total **O(|V|² + |E|) = O(|V|²)**. Com uma fila de prioridade (heap binário), a escolha do mínimo cai para O(log |V|), e cada relaxamento que melhora faz uma inserção O(log |V|). O total fica **O((|V| + |E|) log |V|)**, bem melhor em grafos esparsos. Em grafos densos (|E| ≈ |V|²), a versão linear empata ou ganha. A versão com heap está em `shortest.DijkstraHeap` (desafio).

## Parte VII – Desafio

**Grafo direcionado:** só a construção mudou. As mesmas 9 arestas, passadas a `BuildWeighted` com `directed = true`, viram arcos u → v sem a ligação de retorno. O resultado muda porque C → B não existe mais, só B → C. B fica com 4 pela aresta direta, e a rota até F passa a ser **A → B → D → E → F, custo 14**.

```text
Lista de adjacência (direcionada):
  A → [(B, 4), (C, 2)]
  B → [(C, 1), (D, 5)]
  C → [(D, 8), (E, 10)]
  D → [(E, 2), (F, 6)]
  E → [(F, 3)]
  F → []
```

**Versão com heap:** `shortest.DijkstraHeap` usa `container/heap`, o equivalente em Go ao `heapq`. O heap guarda pares (distância, vértice), ordenados por distância e, em empate, pelo menor vértice, a mesma regra da versão linear. A remoção é “preguiçosa”: um vértice pode entrar mais de uma vez no heap, e as entradas desatualizadas são descartadas ao sair.

```go
// DijkstraHeap é a versão com fila de prioridade (desafio): em vez de varrer os
// n vértices para achar o menor, retira o topo de um heap binário. Custo
// O((V + E) log V). Usa "inserção preguiçosa": um vértice pode entrar mais de
// uma vez no heap, e as entradas desatualizadas são descartadas ao sair.
func DijkstraHeap(g *graph.Weighted, src int, onStep StepFunc) (Result, error) {
    if err := checkNonNegative(g); err != nil {
        return Result{}, err
    }
    n := g.Order()
    dist, pred := initVectors(n, src)
    closed := make([]bool, n)

    pq := &minHeap{{vertex: src, dist: 0}}
    for pq.Len() > 0 {
        top := heap.Pop(pq).(entry)
        u := top.vertex
        if closed[u] || top.dist != dist[u] {
            continue // entrada antiga: u já foi fechado com uma distância menor
        }
        closed[u] = true

        for _, a := range g.Arcs(u) {
            if closed[a.To] {
                continue
            }
            if candidate := dist[u] + a.Weight; candidate < dist[a.To] {
                dist[a.To] = candidate
                pred[a.To] = u
                heap.Push(pq, entry{vertex: a.To, dist: candidate})
            }
        }
        if onStep != nil {
            onStep(Step{Vertex: u, Dist: slices.Clone(dist), Pred: slices.Clone(pred)})
        }
    }
    return Result{Source: src, Dist: dist, Pred: pred}, nil
}
```

**Comparação:** nos dois grafos, as duas versões processam os vértices na mesma ordem e chegam às mesmas distâncias e predecessores; o teste `TestDijkstraHeapMatchesLinear` confere isso. A diferença está só no custo de achar o próximo vértice: varredura O(|V|) contra `heap.Pop` O(log |V|).

```text
Fila de prioridade (heap)
-------------------------
  Processando A (distância 0)
    distancia: {A: 0, B: 4, C: 2, D: ∞, E: ∞, F: ∞}
    anterior:  {A: –, B: A, C: A, D: –, E: –, F: –}
  Processando C (distância 2)
    distancia: {A: 0, B: 4, C: 2, D: 10, E: 12, F: ∞}
    anterior:  {A: –, B: A, C: A, D: C, E: C, F: –}
  Processando B (distância 4)
    distancia: {A: 0, B: 4, C: 2, D: 9, E: 12, F: ∞}
    anterior:  {A: –, B: A, C: A, D: B, E: C, F: –}
  Processando D (distância 9)
    distancia: {A: 0, B: 4, C: 2, D: 9, E: 11, F: 15}
    anterior:  {A: –, B: A, C: A, D: B, E: D, F: D}
  Processando E (distância 11)
    distancia: {A: 0, B: 4, C: 2, D: 9, E: 11, F: 14}
    anterior:  {A: –, B: A, C: A, D: B, E: D, F: E}
  Processando F (distância 14)
    distancia: {A: 0, B: 4, C: 2, D: 9, E: 11, F: 14}
    anterior:  {A: –, B: A, C: A, D: B, E: D, F: E}
  A → F: A → B → D → E → F          custo 14, 4 aresta(s)
```

## Entrega

- **Código:** `cmd/unidade2/dijkstra/main.go` (programa principal), `internal/graph/weighted.go` (criar_grafo), `internal/shortest/dijkstra.go` (dijkstra, versão com heap e reconstruir_caminho).
- **Comentários:** os passos 1 a 4 do pseudocódigo estão marcados no código (inicialização, seleção, relaxamento e reconstrução).
- **Registro de execução:** `go run ./cmd/unidade2/dijkstra`; os trechos estão nas respostas acima.
- **Testes automatizados:** `go test ./internal/shortest/` confere a tabela de rastreamento, as rotas e os testes da Questão 6.
- Nenhuma biblioteca pronta de Dijkstra foi usada; o código usa só a biblioteca padrão de Go.
