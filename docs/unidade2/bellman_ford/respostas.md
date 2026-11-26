# Estudo dirigido – Algoritmo de Bellman-Ford

> Enunciado: [`enunciado.docx`](./enunciado.docx) · Código: [`cmd/unidade2/bellman_ford`](../../../cmd/unidade2/bellman_ford) e [`internal/shortest/bellman_ford.go`](../../../internal/shortest/bellman_ford.go)
>
> Linguagem: **Go** (o enunciado sugere Python; a lógica dos pseudocódigos foi preservada). Vértices A..F = índices 0..5. As arestas são relaxadas na ordem de inserção do enunciado.

```bash
go run ./cmd/unidade2/bellman_ford             # todas as partes e testes
go run ./cmd/unidade2/bellman_ford -start C    # desafio: outra origem
go run ./cmd/unidade2/bellman_ford -reverse    # desafio: arestas em ordem inversa
```

## Questão 1 – Grafo dirigido em lista de adjacência

A tradução foi feita em **Go**. Os vértices A..F são os índices 0..5. Cada aresta é um `graph.WeightedEdge{U, V, Weight}`, o equivalente à tupla (origem, destino, peso).

**1. As nove arestas, na ordem da Figura 1**

```go
// arestas na ordem de inserção da Figura 1. A ordem importa: é nela que cada
// passada relaxa as arestas.
var arestas = []graph.WeightedEdge{
    {U: A, V: B, Weight: 4}, {U: A, V: C, Weight: 2}, {U: C, V: B, Weight: -1},
    {U: B, V: D, Weight: 2}, {U: C, V: D, Weight: 5}, {U: D, V: E, Weight: -2},
    {U: E, V: F, Weight: 3}, {U: B, V: F, Weight: 10}, {U: C, V: E, Weight: 7},
}
```

**2. criar_grafo e listas de A até F**

O pseudocódigo `BellmanFord` percorre a **lista de arestas** (“para cada (u, v, peso) em arestas”), então `shortest.BellmanFord` recebe a lista diretamente, preservando a ordem de inserção. A lista de adjacência é montada para exibição em `printAdjacency`: cria uma lista vazia por vértice (PASSO 1), acrescenta `(v, peso)` só em `lista[u]` (PASSO 2) e escreve as listas (PASSO 3).

```go
// printAdjacency monta a lista na ordem de inserção das arestas.
func printAdjacency(n int, edges []graph.WeightedEdge) {
    adj := make([][]graph.WeightedEdge, n)
    for _, e := range edges {
        adj[e.U] = append(adj[e.U], e) // apenas u → v, sem ligação de retorno
    }
    for u, list := range adj {
        parts := make([]string, len(list))
        for i, e := range list {
            parts[i] = fmt.Sprintf("(%s, %d)", r.Label(e.V), e.Weight)
        }
        fmt.Printf("  %s → [%s]\n", r.Label(u), strings.Join(parts, ", "))
    }
}
```

```text
  A → [(B, 4), (C, 2)]
  B → [(D, 2), (F, 10)]
  C → [(B, -1), (D, 5), (E, 7)]
  D → [(E, -2)]
  E → [(F, 3)]
  F → []
  Arestas negativas: C → B (-1) D → E (-2)
```

O tipo `graph.Weighted` da lista anterior não foi reaproveitado aqui por dois motivos. Ele reordena as listas por vizinho, e a ordem das arestas muda os estados das passadas. E ele rejeita laços u → u, que o teste 6 precisa (G → G).

**3. Por que não inserir a ligação de retorno?**

Porque o grafo é **dirigido**: a aresta C → B(−1) permite ir de C para B, não de B para C. Inserir o retorno criaria caminhos que não existem e mudaria as respostas. Com pesos negativos, o estrago é maior: o retorno B → C(−1) formaria o ciclo C → B → C com peso −2, um ciclo negativo, e nenhuma distância a partir de A seria mais bem definida.

**4. Arestas negativas**

**C → B (−1)** e **D → E (−2)**. Um peso negativo representa um ganho ao percorrer a aresta: um desconto, um bônus ou um crédito (por exemplo, um trecho em que se recebe em vez de pagar). Passar por ela **reduz** o custo acumulado; por isso uma rota com mais arestas pode sair mais barata, e o Dijkstra não serve.

## Questão 2 – Bellman-Ford comentado

`shortest.BellmanFord(n, arestas, origem)` (`internal/shortest/bellman_ford.go`) segue os quatro passos do pseudocódigo. `dist` é `distancia`, `pred` é `anterior`, `changed` é `mudou`, `Inf` é infinito e `NoPred` é NULO. O “Escrever i, distancia, anterior” vira um registro `Pass` por passada, que o programa imprime em forma de tabela.

```go
// BellmanFord calcula os menores custos a partir de src aceitando pesos
// negativos. Recebe a LISTA de arestas, e não a lista de adjacência, porque a
// ordem em que as arestas são relaxadas muda os estados intermediários das
// passadas (o enunciado fixa essa ordem). Laços u → u são aceitos, já que um
// laço negativo é um ciclo negativo (teste 6).
func BellmanFord(n int, edges []graph.WeightedEdge, src int) (BFResult, error) {
    for _, e := range edges {
        if e.U < 0 || e.U >= n || e.V < 0 || e.V >= n {
            return BFResult{}, fmt.Errorf("%w: (%d, %d) com n = %d", graph.ErrInvalidVertex, e.U, e.V, n)
        }
    }

    // PASSO 1 – inicializar distâncias (∞; origem 0) e predecessores (NULO).
    dist, pred := initVectors(n, src)
    var passes []Pass

    // PASSO 2 – repetir até |V|−1 passadas por TODAS as arestas. Um caminho
    // mínimo simples tem no máximo |V|−1 arestas; a passada i garante as rotas
    // mínimas de até i arestas.
    for i := 1; i <= n-1; i++ {
        changed := false // mudou ← falso
        for _, e := range edges {
            // Só relaxa a partir de u alcançável: ∞ + peso não é uma rota.
            if dist[e.U] == Inf {
                continue
            }
            candidate := dist[e.U] + e.Weight
            if candidate < dist[e.V] {
                dist[e.V] = candidate // relaxamento da aresta u → v
                pred[e.V] = e.U       // salvar predecessor
                changed = true
            }
        }
        passes = append(passes, Pass{Index: i, Dist: slices.Clone(dist), Pred: slices.Clone(pred), Changed: changed})
        // Parada antecipada: se uma passada inteira não melhorou nada, as
        // seguintes repetiriam exatamente os mesmos testes.
        if !changed {
            break
        }
    }

    // PASSO 3 – testar ciclo negativo alcançável: depois de |V|−1 passadas, as
    // distâncias já são mínimas, a não ser que haja um ciclo negativo. Se alguma
    // aresta com u alcançável ainda melhora v, o ciclo existe.
    negative := false
    for _, e := range edges {
        if dist[e.U] != Inf && dist[e.U]+e.Weight < dist[e.V] {
            negative = true
            break
        }
    }

    // PASSO 4 – devolver resultados.
    return BFResult{Result: Result{Source: src, Dist: dist, Pred: pred}, Passes: passes, NegativeCycle: negative}, nil
}
```

- **Laço externo:** `for i := 1; i <= n-1; i++`, no máximo |V| − 1 passadas completas.
- **Só u alcançável:** `if dist[e.U] == Inf { continue }`; ∞ + peso não representa uma rota.
- **Melhora estrita:** atualiza `dist` e `pred` só se `candidate < dist[e.V]`.
- **Parada antecipada:** se a passada terminou com `changed == false`, interrompe (`break`).
- **Ciclo negativo:** uma passada extra; se alguma aresta com u alcançável ainda melhora v, `NegativeCycle = true`.

## Questão 3 – Acompanhamento das passadas

| Etapa | d(A) | d(B) | d(C) | d(D) | d(E) | d(F) |
|---|---|---|---|---|---|---|
| Inicial | 0 | ∞ | ∞ | ∞ | ∞ | ∞ |
| Passada 1 | 0 | 1 | 2 | 3 | 1 | 4 |
| Passada 2 | 0 | 1 | 2 | 3 | 1 | 4 |
| Passada 3 | 0 | não exec. | não exec. | não exec. | não exec. | não exec. |
| Passada 4 | 0 | não exec. | não exec. | não exec. | não exec. | não exec. |
| Passada 5 | 0 | não exec. | não exec. | não exec. | não exec. | não exec. |

| Momento | ant(B) | ant(C) | ant(D) | ant(E) | ant(F) |
|---|---|---|---|---|---|
| Inicial | – | – | – | – | – |
| Após 1ª passada | C | A | B | D | E |
| Final | C | A | B | D | E |

As tabelas do enunciado estão preenchidas acima. As passadas 3, 4 e 5 **não foram executadas** por causa da parada antecipada. Se fossem, repetiriam 0, 1, 2, 3, 1, 4, porque nenhuma aresta melhora mais nada. Detalhe da passada 1, com as arestas na ordem do enunciado:

| Aresta | Teste | Resultado |
|---|---|---|
| A→B(4) | 0 + 4 = 4 < ∞ | B = 4, ant(B) = A |
| A→C(2) | 0 + 2 = 2 < ∞ | C = 2, ant(C) = A |
| C→B(−1) | 2 − 1 = 1 < 4 | **B = 1, ant(B) = C** |
| B→D(2) | 1 + 2 = 3 < ∞ | D = 3, ant(D) = B |
| C→D(5) | 2 + 5 = 7, não é < 3 | sem mudança |
| D→E(−2) | 3 − 2 = 1 < ∞ | E = 1, ant(E) = D |
| E→F(3) | 1 + 3 = 4 < ∞ | F = 4, ant(F) = E |
| B→F(10) | 1 + 10 = 11, não é < 4 | sem mudança |
| C→E(7) | 2 + 7 = 9, não é < 1 | sem mudança |

**Saída do programa**

```text
  Ordem de relaxamento: A→B(4), A→C(2), C→B(-1), B→D(2), C→D(5), D→E(-2), E→F(3), B→F(10), C→E(7)

  Etapa       A    B    C    D    E    F
  Inicial     0    ∞    ∞    ∞    ∞    ∞
  Passada 1   0    1    2    3    1    4
  Passada 2   0    1    2    3    1    4

  anterior    A    B    C    D    E    F
  Passada 1   –    C    A    B    D    E
  Passada 2   –    C    A    B    D    E

  Passada 2 não alterou nada: parada antecipada (2 de 5 passadas executadas).
  Nenhum ciclo negativo alcançável.
```

**10. Como C→B(−1) modifica o custo de B?**

A→B(4) atribui B = 4 pela aresta direta. Logo depois, na mesma passada, A→C(2) define C = 2, e C→B(−1) dá 2 + (−1) = 1 < 4. Então B cai para **1** e anterior[B] passa de A para **C**. Como a atualização acontece durante a passada, o novo valor de B já é usado em B→D, e a melhoria se propaga para D, E e F na mesma rodada.

**11. Em qual passada o algoritmo deixa de melhorar?**

Na **passada 2**. A passada 1 já chega aos valores finais (0, 1, 2, 3, 1, 4), porque a ordem das arestas segue o sentido das rotas mínimas: cada aresta da rota A→C→B→D→E→F é relaxada depois da anterior. Na passada 2, nenhum dos 9 testes encontra `candidato < distancia[v]`, e `mudou` termina falso.

**12. Por que dá para encerrar antes das cinco passadas?**

Uma passada só depende das distâncias atuais. Se uma passada inteira não mudou nada, a próxima faria exatamente os mesmos testes, com os mesmos valores, e também não mudaria nada. As distâncias chegaram a um ponto fixo em que toda aresta satisfaz `distancia[v] ≤ distancia[u] + peso`, e isso só acontece com os custos mínimos. |V| − 1 é apenas o limite do pior caso, quando cada passada fixa uma aresta a mais da rota.

## Questão 4 – Reconstrução do caminho mínimo

`reconstruir_caminho` é `Result.PathTo` (compartilhado com o Dijkstra), já na versão deste enunciado: devolve vazio se o destino é ∞ (PASSO 1), recua pelos predecessores e **retorna assim que chega à origem** (PASSO 2). Se a cadeia não chega à origem, devolve vazio. O laço tem limite de n passos, então nunca fica preso quando os predecessores formam um ciclo, que é o que acontece com ciclo negativo. Mesmo assim, o programa só reconstrói rotas se `NegativeCycle` for falso.

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

```text
  A → F: A → C → B → D → E → F      custo 4, 5 aresta(s)
  A → B: A → C → B                  custo 1, 2 aresta(s)
  A → D: A → C → B → D              custo 3, 3 aresta(s)
  A → E: A → C → B → D → E          custo 1, 4 aresta(s)
```

**14.** A → F: **A → C → B → D → E → F, custo mínimo 4** (2 − 1 + 2 − 2 + 3). **15.** B: A → C → B (1); D: A → C → B → D (3); E: A → C → B → D → E (1).

**16. Por que o número de arestas não determina o custo?**

O custo é a soma dos pesos, e com pesos negativos uma aresta a mais pode **diminuir** essa soma. A rota mínima até F tem 5 arestas e custa 4, enquanto A → B → F tem 2 arestas e custa 14. Até E, a rota de 4 arestas custa 1, contra 9 de A → C → E, com 2 arestas.

## Questão 5 – Estrutura do programa

| Esqueleto (Python) | Implementação (Go) | Arquivo |
|---|---|---|
| criar_grafo(vertices, arestas) | `printAdjacency(n, arestas)` (exibição) + lista `arestas` | cmd/unidade2/bellman_ford/main.go |
| bellman_ford(vertices, arestas, origem) | `shortest.BellmanFord(n, arestas, origem)` | internal/shortest/bellman_ford.go |
| reconstruir_caminho(...) | `Result.PathTo(destino)` | internal/shortest/dijkstra.go |
| programa principal | `run(origem, reverse)` | cmd/unidade2/bellman_ford/main.go |

```go
func run(origem int, reverse bool) error {
    edges := arestas
    if reverse {
        edges = slices.Clone(arestas)
        slices.Reverse(edges)
    }

    // PARTE I – criar_grafo: cada arco só no sentido u → v.
    r.Title("PARTE I – Lista de adjacência (direcionada)")
    printAdjacency(F+1, arestas)
    fmt.Print("  Arestas negativas:")
    for _, e := range arestas {
        if e.Weight < 0 {
            fmt.Printf(" %s → %s (%d)", r.Label(e.U), r.Label(e.V), e.Weight)
        }
    }
    fmt.Println()

    // PARTES II e III – passadas.
    title := fmt.Sprintf("PARTES II e III – Bellman-Ford a partir de %s", r.Label(origem))
    if reverse {
        title += " (arestas em ordem inversa)"
    }
    r.Title(title)
    fmt.Println("  Ordem de relaxamento:", edgeOrder(edges))
    res, err := shortest.BellmanFord(F+1, edges, origem)
    if err != nil {
        return err
    }
    printPasses(res)

    // PARTE IV – reconstruir_caminho.
    r.Title("PARTE IV – Reconstrução dos caminhos mínimos")
    for _, dest := range []int{F, B, D, E} {
        printRoute(res, dest)
    }

    if origem != A || reverse {
        return nil // as variações do desafio só repetem as partes II a IV
    }

    r.Title("PARTE VI – Testes")
    return runTests(res)
}
```

Execução: `go run ./cmd/unidade2/bellman_ford` (flags do desafio: `-start C`, `-reverse`). Os testes automatizados ficam em `internal/shortest/unidade2_test.go`. Nenhuma biblioteca de caminhos mínimos foi usada.

## Questão 6 – Testes e análise

| Teste | Resultado | Passadas | Verificação |
|---|---|---|---|
| 1 – Original; A → F | A → C → B → D → E → F, custo 4 | 2 (parada antecipada) | Rota e custo mínimos. |
| 2 – Original; A → A | [A], custo 0 | 2 | anterior[A] = NULO; a rota é só a origem. |
| 3 – Original; A → D | A → C → B → D, custo 3 | 2 | anterior[D] = B, anterior[B] = C, anterior[C] = A. |
| 4 – G isolado; A → G | distância ∞, caminho [] | 2 (de 6 possíveis) | Nenhuma aresta chega a G. |
| 5 – E → B (−5) | **ciclo negativo detectado** | 5 (todas mudaram) | Ciclo B → D → E → B = 2 − 2 − 5 = −5; as rotas não são reconstruídas. |
| 6 – G → G (−1), G isolado | sem alerta; A → F = 4; G = ∞ | 2 | O laço negativo existe, mas G é inalcançável desde A. |

**17. Saída e estados das passadas em cada teste**

```text
Teste 1 – Original; A → F
-------------------------
  A → F: A → C → B → D → E → F      custo 4, 5 aresta(s)

Teste 2 – Original; A → A
-------------------------
  A → A: A                          custo 0, 0 aresta(s)

Teste 3 – Original; A → D
-------------------------
  A → D: A → C → B → D              custo 3, 3 aresta(s)
  Predecessores: anterior[D] = B, anterior[B] = C, anterior[C] = A

Teste 4 – G isolado; A → G
--------------------------

  Etapa       A    B    C    D    E    F    G
  Inicial     0    ∞    ∞    ∞    ∞    ∞    ∞
  Passada 1   0    1    2    3    1    4    ∞
  Passada 2   0    1    2    3    1    4    ∞

  anterior    A    B    C    D    E    F    G
  Passada 1   –    C    A    B    D    E    –
  Passada 2   –    C    A    B    D    E    –

  Passada 2 não alterou nada: parada antecipada (2 de 6 passadas executadas).
  Nenhum ciclo negativo alcançável.
  A → G: inalcançável (distância ∞, caminho [])

Teste 5 – Acrescentar E → B com peso −5
---------------------------------------

  Etapa       A    B    C    D    E    F
  Inicial     0    ∞    ∞    ∞    ∞    ∞
  Passada 1   0    -4   2    3    1    4
  Passada 2   0    -9   2    -2   -4   -1
  Passada 3   0    -14  2    -7   -9   -6
  Passada 4   0    -19  2    -12  -14  -11
  Passada 5   0    -24  2    -17  -19  -16

  anterior    A    B    C    D    E    F
  Passada 1   –    E    A    B    D    E
  Passada 2   –    E    A    B    D    E
  Passada 3   –    E    A    B    D    E
  Passada 4   –    E    A    B    D    E
  Passada 5   –    E    A    B    D    E

  Todas as 5 passadas alteraram alguma distância.
  ⚠ Ciclo negativo alcançável detectado: as distâncias não são mínimas bem definidas.
  A → F: não reconstruída (ciclo negativo alcançável)

Teste 6 – Acrescentar G → G com peso −1 (G isolado)
---------------------------------------------------

  Etapa       A    B    C    D    E    F    G
  Inicial     0    ∞    ∞    ∞    ∞    ∞    ∞
  Passada 1   0    1    2    3    1    4    ∞
  Passada 2   0    1    2    3    1    4    ∞

  anterior    A    B    C    D    E    F    G
  Passada 1   –    C    A    B    D    E    –
  Passada 2   –    C    A    B    D    E    –

  Passada 2 não alterou nada: parada antecipada (2 de 6 passadas executadas).
  Nenhum ciclo negativo alcançável.
  A → F: A → C → B → D → E → F      custo 4, 5 aresta(s)
  A → G: inalcançável (distância ∞, caminho [])
```

**18. Bellman-Ford × Dijkstra**

Os dois resolvem o caminho mínimo a partir de **uma** origem. Dijkstra é mais rápido, O(|V|²) ou O((|V| + |E|) log |V|), mas exige pesos **não negativos**, porque fecha cada vértice de vez. Bellman-Ford custa O(|V|·|E|), mas aceita pesos negativos e **detecta ciclos negativos alcançáveis**. É necessário quando há arestas negativas, como neste grafo (C→B e D→E), ou quando é preciso saber se existe um ciclo negativo. Sem pesos negativos, os dois dão os mesmos custos; o teste `TestBellmanFordMatchesDijkstraWithoutNegatives` confere isso no grafo do estudo de Dijkstra.

**19. Condição de ciclo negativo alcançável**

Depois das |V| − 1 passadas, o algoritmo testa cada aresta mais uma vez: se `distancia[u] ≠ ∞` e `distancia[u] + peso < distancia[v]`, há ciclo negativo. **Justificativa:** sem ciclo negativo, toda rota mínima pode ser simples, com no máximo |V| − 1 arestas, e já foi encontrada; nenhuma aresta melhora mais nada. Se alguma ainda melhora, existe uma rota mais barata com |V| ou mais arestas. Essa rota repete um vértice, e o trecho repetido é um ciclo de peso negativo. A condição `distancia[u] ≠ ∞` restringe o teste a arestas que partem de vértices **alcançáveis** a partir da origem. Por isso o laço G → G do teste 6 não é sinalizado.

**20. Complexidade**

**O(|V|·|E|)** no tempo: até |V| − 1 passadas, cada uma testando as |E| arestas, mais uma passada de verificação O(|E|). Com a parada antecipada, o melhor caso é O(|E|): aqui foram 2 passadas, ou 18 relaxamentos, em vez de 5 × 9 = 45. O espaço extra é O(|V|) para `distancia` e `anterior`, além da lista de arestas.

## Parte VII – Desafio

O programa aceita `-start` (qualquer origem de A a F) e `-reverse` (relaxa as arestas na ordem inversa à do enunciado). O grafo continua o mesmo.

**Origem C (go run ./cmd/unidade2/bellman_ford -start C)**

```text
  Etapa       A    B    C    D    E    F
  Inicial     ∞    ∞    0    ∞    ∞    ∞
  Passada 1   ∞    -1   0    1    -1   2
  Passada 2   ∞    -1   0    1    -1   2

  anterior    A    B    C    D    E    F
  Passada 1   –    C    –    B    D    E
  Passada 2   –    C    –    B    D    E

  Passada 2 não alterou nada: parada antecipada (2 de 5 passadas executadas).
  Nenhum ciclo negativo alcançável.

  C → F: C → B → D → E → F          custo 2, 4 aresta(s)
  C → B: C → B                      custo -1, 1 aresta(s)
  C → D: C → B → D                  custo 1, 2 aresta(s)
  C → E: C → B → D → E              custo -1, 3 aresta(s)
```

**Arestas em ordem inversa (go run ./cmd/unidade2/bellman_ford -reverse)**

```text
  Ordem de relaxamento: C→E(7), B→F(10), E→F(3), D→E(-2), C→D(5), B→D(2), C→B(-1), A→C(2), A→B(4)

  Etapa       A    B    C    D    E    F
  Inicial     0    ∞    ∞    ∞    ∞    ∞
  Passada 1   0    4    2    ∞    ∞    ∞
  Passada 2   0    1    2    6    9    12
  Passada 3   0    1    2    3    4    11
  Passada 4   0    1    2    3    1    7
  Passada 5   0    1    2    3    1    4

  anterior    A    B    C    D    E    F
  Passada 1   –    A    A    –    –    –
  Passada 2   –    C    A    B    C    E
  Passada 3   –    C    A    B    D    B
  Passada 4   –    C    A    B    D    E
  Passada 5   –    C    A    B    D    E

  Todas as 5 passadas alteraram alguma distância.
  Nenhum ciclo negativo alcançável.
```

| Ordem das arestas | Com parada antecipada | Sem parada antecipada |
|---|---|---|
| Enunciado (A→B, A→C, C→B, ...) | 2 passadas (a 2ª só confirma) | 5 passadas |
| Inversa (C→E, B→F, ..., A→B) | 5 passadas (todas mudaram algo) | 5 passadas |

**Por que a ordem muda os estados intermediários, mas não os custos finais?** Na ordem do enunciado, cada aresta da rota mínima é relaxada depois da anterior, e a melhoria percorre a rota inteira numa só passada. Na ordem inversa, cada aresta é relaxada **antes** de a anterior ter sido atualizada, e cada passada avança só uma aresta da rota. Por isso a passada 1 conhece apenas B e C, e F só chega a 4 na passada 5. Em qualquer ordem, porém, depois da passada i estão corretas todas as distâncias cujas rotas mínimas têm até i arestas. Sem ciclo negativo alcançável, |V| − 1 passadas cobrem qualquer rota mínima, e o resultado é o único ponto fixo, que são os custos mínimos. Só os caminhos para chegar lá mudam.

## Entrega

- **Código:** `cmd/unidade2/bellman_ford/main.go` (programa principal e criar_grafo), `internal/shortest/bellman_ford.go` (bellman_ford) e `internal/shortest/dijkstra.go` (reconstruir_caminho).
- **Comentários:** inicialização, relaxamento, parada antecipada, teste de ciclo e reconstrução estão marcados no código.
- **Registros de execução:** `go run ./cmd/unidade2/bellman_ford`; tabela de passadas preenchida e trechos da saída nas respostas acima.
- **Testes automatizados:** `go test ./internal/shortest/`, com a passada 1, o ciclo negativo alcançável e o laço inalcançável.
- Nenhuma implementação pronta de Bellman-Ford foi usada; o código usa só a biblioteca padrão de Go.
