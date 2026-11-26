# Estudo dirigido – Algoritmo de Floyd-Warshall

> Enunciado: [`enunciado.docx`](./enunciado.docx) · Código: [`cmd/unidade2/floyd_warshall`](../../../cmd/unidade2/floyd_warshall) e [`internal/shortest/floyd_warshall.go`](../../../internal/shortest/floyd_warshall.go)
>
> Linguagem: **Go** (o enunciado sugere Python; a lógica dos pseudocódigos foi preservada). Vértices A..F = índices 0..5; linha = origem, coluna = destino.

```bash
go run ./cmd/unidade2/floyd_warshall   # todas as partes, testes e desafio
```

## Questão 1 – Matriz de adjacência ponderada

A tradução foi feita em **Go**. As matrizes são o tipo `shortest.Matrix` (`[][]int`, uma lista de listas). `Inf` (`math.MaxInt`) faz o papel de `float("inf")` e `NoPred` (−1) o de `None`. O dicionário rótulo → índice é dado pelas constantes `A = 0, ..., F = 5` do programa, e `report.Label` faz a conversão inversa para exibir.

**1 e 2. criar_matrizes: matrizes 6×6 e as nove arestas orientadas**

```go
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
```

Só `dist[u][v]` recebe o peso: como o grafo é dirigido, `dist[v][u]` continua ∞. A condição `peso < dist[u][v]` guarda a menor aresta direta, caso houvesse arestas paralelas.

**3. Matriz inicial com rótulos**

```text
  i \ j  A    B    C    D    E    F
  A      0    4    2    ∞    ∞    ∞
  B      ∞    0    ∞    2    ∞    10
  C      ∞    -1   0    5    7    ∞
  D      ∞    ∞    ∞    0    -2   ∞
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0

  Matriz prox inicial:
  i \ j  A    B    C    D    E    F
  A      –    B    C    –    –    –
  B      –    –    –    D    –    F
  C      –    B    –    D    E    –
  D      –    –    –    –    E    –
  E      –    –    –    –    –    F
  F      –    –    –    –    –    –
```

**4. Por que dist[B][A] e dist[A][B] podem ser diferentes?**

Porque o grafo é **dirigido**: a linha é a origem e a coluna o destino. Existe o arco A → B (4), então `dist[A][B] = 4`, mas não existe B → A, então `dist[B][A] = ∞`. Mesmo no final, nenhum vértice volta a A, e `dist[B][A]` continua ∞. Num grafo não dirigido, a matriz seria simétrica.

## Questão 2 – Floyd-Warshall comentado

`shortest.FloydWarshall(dist, prox, onStage)` (`internal/shortest/floyd_warshall.go`) trabalha sobre cópias das matrizes. `onStage` é o “Escrever Intermediário, k, dist” e recebe a matriz depois de cada k.

```go
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
```

- **Ordem k → i → j:** k é o laço externo. A etapa k só pode usar os resultados completos da etapa k − 1; com k por dentro, o algoritmo combinaria caminhos incompletos.
- **Só trechos finitos:** `if d[i][k] == Inf || d[k][j] == Inf { continue }`. Além de ser a regra, isso evita estouro, porque `Inf` é o maior `int`.
- **Melhora estrita:** atualiza `dist[i][j]` e `prox[i][j] = prox[i][k]` só se `candidato < dist[i][j]`.
- **Ciclos negativos:** depois dos laços, qualquer `dist[v][v] < 0` liga `NegativeCycle`.

## Questão 3 – Matrizes por etapas

**Matriz D⁽ᴮ⁾ (após A e B)**

| i \ j | A | B | C | D | E | F |
|---|---|---|---|---|---|---|
| A | 0 | 4 | 2 | 6 | ∞ | 14 |
| B | ∞ | 0 | ∞ | 2 | ∞ | 10 |
| C | ∞ | -1 | 0 | 1 | 7 | 9 |
| D | ∞ | ∞ | ∞ | 0 | -2 | ∞ |
| E | ∞ | ∞ | ∞ | ∞ | 0 | 3 |
| F | ∞ | ∞ | ∞ | ∞ | ∞ | 0 |

**Matriz final**

| i \ j | A | B | C | D | E | F |
|---|---|---|---|---|---|---|
| A | 0 | 1 | 2 | 3 | 1 | 4 |
| B | ∞ | 0 | ∞ | 2 | 0 | 3 |
| C | ∞ | -1 | 0 | 1 | -1 | 2 |
| D | ∞ | ∞ | ∞ | 0 | -2 | 1 |
| E | ∞ | ∞ | ∞ | ∞ | 0 | 3 |
| F | ∞ | ∞ | ∞ | ∞ | ∞ | 0 |

As duas matrizes do enunciado estão preenchidas acima. Considerar A como intermediário não muda nada, porque nenhum arco chega a A. Logo D⁽ᴬ⁾ = D⁽⁰⁾, e D⁽ᴮ⁾ só melhora A→D, A→F e as entradas de C. A saída do programa, etapa por etapa:

```text
Intermediário A – D(A)
----------------------
  i \ j  A    B    C    D    E    F
  A      0    4    2    ∞    ∞    ∞
  B      ∞    0    ∞    2    ∞    10
  C      ∞    -1   0    5    7    ∞
  D      ∞    ∞    ∞    0    -2   ∞
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0

Intermediário B – D(B)
----------------------
  i \ j  A    B    C    D    E    F
  A      0    4    2    6    ∞    14
  B      ∞    0    ∞    2    ∞    10
  C      ∞    -1   0    1    7    9
  D      ∞    ∞    ∞    0    -2   ∞
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0

Intermediário C – D(C)
----------------------
  i \ j  A    B    C    D    E    F
  A      0    1    2    3    9    11
  B      ∞    0    ∞    2    ∞    10
  C      ∞    -1   0    1    7    9
  D      ∞    ∞    ∞    0    -2   ∞
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0

Intermediário D – D(D)
----------------------
  i \ j  A    B    C    D    E    F
  A      0    1    2    3    1    11
  B      ∞    0    ∞    2    0    10
  C      ∞    -1   0    1    -1   9
  D      ∞    ∞    ∞    0    -2   ∞
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0

Intermediário E – D(E)
----------------------
  i \ j  A    B    C    D    E    F
  A      0    1    2    3    1    4
  B      ∞    0    ∞    2    0    3
  C      ∞    -1   0    1    -1   2
  D      ∞    ∞    ∞    0    -2   1
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0

Intermediário F – D(F)
----------------------
  i \ j  A    B    C    D    E    F
  A      0    1    2    3    1    4
  B      ∞    0    ∞    2    0    3
  C      ∞    -1   0    1    -1   2
  D      ∞    ∞    ∞    0    -2   1
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0
```

**10. Melhoria de A→D com o intermediário B**

Em D⁽ᴬ⁾, `dist[A][D] = ∞`, porque não há arco A → D. Com B liberado: `dist[A][B] + dist[B][D] = 4 + 2 = 6 < ∞`, então `dist[A][D] = 6` e `prox[A][D] = prox[A][B] = B`. Na etapa C, a distância cai de novo: `dist[A][C] + dist[C][D] = 2 + 1 = 3`, pela rota A → C → B → D.

**11. Atualização que usa a aresta negativa D→E**

Na etapa **k = D**: `dist[A][E] = dist[A][D] + dist[D][E] = 3 + (−2) = 1`, menor que o 9 de antes (A → C → E). Na mesma etapa: `dist[B][E] = 2 + (−2) = 0` e `dist[C][E] = 1 + (−2) = −1`.

**12. Custo mínimo A→F e simetria**

**dist[A][F] = 4**, pela rota A → C → B → D → E → F: 2 − 1 + 2 − 2 + 3. A matriz final **não é simétrica**. Por exemplo, `dist[A][F] = 4` e `dist[F][A] = ∞`, ou `dist[C][B] = −1` e `dist[B][C] = ∞`. Isso é esperado num grafo dirigido.

**13. Diagonal**

A diagonal é toda **0**: o melhor custo de v até ele mesmo é ficar parado. Nenhuma etapa conseguiu um `dist[v][v] < 0`, então **não há ciclo negativo**. Neste grafo, aliás, não há ciclo nenhum: nenhum vértice volta a si mesmo.

## Questão 4 – Reconstrução dos caminhos mínimos

```go
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
```

**14.** `APSP.Path(i, j)` usa índices de modo consistente e converte para rótulos só na exibição. Segue o pseudocódigo: devolve `[i]` quando i = j, vazio quando `prox[i][j]` é NULO, e então anda `atual = prox[atual][destino]` até o destino. O limite de n passos sugerido no enunciado também está lá. O programa só apresenta uma rota se ela não passa por ciclo negativo (`APSP.OnNegativeCycle`).

```text
  A → F: A → C → B → D → E → F      dist[A][F] = 4
  A → D: A → C → B → D              dist[A][D] = 3
  C → F: C → B → D → E → F          dist[C][F] = 2
  F → A: não existe rota (dist = ∞)
```

**15.** A → F: A → C → B → D → E → F, com `dist[A][F] = 4`. **16.** A → D: A → C → B → D (3); C → F: C → B → D → E → F (2); F → A: **não existe rota**, porque `prox[F][A]` é NULO e `dist[F][A] = ∞` (F não tem arcos de saída).

**17. anterior[v] (Bellman-Ford) × prox[i][j] (Floyd-Warshall)**

|  | anterior[v] | prox[i][j] |
|---|---|---|
| Escopo | uma origem fixa | todos os pares |
| Guarda | o vértice **antes** de v na rota origem → v | o vértice **depois** de i na rota i → j |
| Estrutura | vetor de n posições | matriz n × n |
| Reconstrução | de trás para frente (destino → origem) e depois inverte | de frente para trás (origem → destino), sem inverter |
| Atualização | `anterior[v] = u` | `prox[i][j] = prox[i][k]` |

## Questão 5 – Estrutura do programa

| Esqueleto (Python) | Implementação (Go) | Arquivo |
|---|---|---|
| criar_matrizes(vertices, arestas) | `shortest.NewMatrices(n, arestas)` | internal/shortest/floyd_warshall.go |
| floyd_warshall(vertices, dist, prox) | `shortest.FloydWarshall(dist, prox, onStage)` | internal/shortest/floyd_warshall.go |
| reconstruir_caminho(prox, origem, destino, vertices) | `APSP.Path(origem, destino)` | internal/shortest/floyd_warshall.go |
| programa principal | `run()` | cmd/unidade2/floyd_warshall/main.go |

```go
func run() error {
    // PARTE I – criar_matrizes.
    dist, prox, err := shortest.NewMatrices(F+1, arestas)
    if err != nil {
        return err
    }
    r.Title("PARTE I – Matriz inicial D(0) (só arestas diretas)")
    r.PrintMatrix(dist, r.Dist)
    fmt.Println("\n  Matriz prox inicial:")
    r.PrintMatrix(prox, r.Pred)

    // PARTES II e III – uma matriz por intermediário k = A, B, ..., F.
    r.Title("PARTES II e III – Floyd-Warshall (k → i → j)")
    res := shortest.FloydWarshall(dist, prox, printStage)
    r.Subtitle("Matriz final de menores custos")
    r.PrintMatrix(res.Dist, r.Dist)
    fmt.Println("\n  Matriz prox final:")
    r.PrintMatrix(res.Next, r.Pred)
    printCycle(res)

    // PARTE IV – reconstruir_caminho.
    r.Title("PARTE IV – Reconstrução dos caminhos mínimos")
    for _, q := range [][2]int{{A, F}, {A, D}, {C, F}, {F, A}} {
        printRoute(res, q[0], q[1])
    }

    r.Title("PARTE VI – Testes")
    if err := runTests(res); err != nil {
        return err
    }
    return challenge(res)
}
```

Execução: `go run ./cmd/unidade2/floyd_warshall`. Os testes automatizados ficam em `internal/shortest/unidade2_test.go`. Nenhuma biblioteca pronta foi usada.

## Questão 6 – Testes e análise

| Teste | Resultado | Verificação |
|---|---|---|
| 1 – Original; todos os pares | matriz final acima; A → F = 4 por A → C → B → D → E → F | Igual ao rastreamento manual. |
| 2 – Original; A → A | custo 0, rota [A] | Caminho trivial (origem = destino). |
| 3 – Original; F → A | ∞, caminho vazio | F não tem arcos de saída. |
| 4 – G isolado | linha e coluna G com ∞, exceto dist[G][G] = 0 | G não alcança nem é alcançado. |
| 5 – E → B (−5) | diagonal B = −5, D = −5, E = −10; alerta de ciclo | A → F não é apresentada como mínima; A → C (2) continua válida. |
| 6 – A → B = 1 | A → F continua **4**, agora por A → B → D → E → F | Empate: 1 + 2 − 2 + 3 = 4 = 2 − 1 + 2 − 2 + 3. |

**18. Matrizes e rotas de cada teste**

```text
Teste 1 – Original; todos os pares
----------------------------------
  i \ j  A    B    C    D    E    F
  A      0    1    2    3    1    4
  B      ∞    0    ∞    2    0    3
  C      ∞    -1   0    1    -1   2
  D      ∞    ∞    ∞    0    -2   1
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0
  A → F: A → C → B → D → E → F      dist[A][F] = 4

Teste 2 – Original; A → A
-------------------------
  A → A: A                          dist[A][A] = 0

Teste 3 – Original; F → A
-------------------------
  F → A: não existe rota (dist = ∞)

Teste 4 – Acrescentar G isolado
-------------------------------
  i \ j  A    B    C    D    E    F    G
  A      0    1    2    3    1    4    ∞
  B      ∞    0    ∞    2    0    3    ∞
  C      ∞    -1   0    1    -1   2    ∞
  D      ∞    ∞    ∞    0    -2   1    ∞
  E      ∞    ∞    ∞    ∞    0    3    ∞
  F      ∞    ∞    ∞    ∞    ∞    0    ∞
  G      ∞    ∞    ∞    ∞    ∞    ∞    0
  A → G: não existe rota (dist = ∞)
  G → A: não existe rota (dist = ∞)

Teste 5 – Acrescentar E → B com peso −5
---------------------------------------
  i \ j  A    B    C    D    E    F
  A      0    -4   2    -2   -4   -1
  B      ∞    -5   ∞    -3   -5   -2
  C      ∞    -6   0    -4   -6   -3
  D      ∞    -7   ∞    -5   -7   -4
  E      ∞    -10  ∞    -8   -10  -7
  F      ∞    ∞    ∞    ∞    ∞    0

  Diagonal: A=0 B=-5 C=0 D=-5 E=-10 F=0
  ⚠ Ciclo negativo detectado (há dist[v][v] < 0).
  A → F: não apresentada como mínima (a rota alcança um ciclo negativo)
  A → C: A → C                      dist[A][C] = 2

Teste 6 – Alterar A → B de 4 para 1
-----------------------------------
  i \ j  A    B    C    D    E    F
  A      0    1    2    3    1    4
  B      ∞    0    ∞    2    0    3
  C      ∞    -1   0    1    -1   2
  D      ∞    ∞    ∞    0    -2   1
  E      ∞    ∞    ∞    ∞    0    3
  F      ∞    ∞    ∞    ∞    ∞    0
  A → F: A → B → D → E → F          dist[A][F] = 4
  Comparação A → F: original A → C → B → D → E → F (custo 4) × alterado A → B → D → E → F (custo 4)
```

No teste 5, os valores da diagonal (−5 e até −10) mostram que o ciclo B → D → E → B (2 − 2 − 5 = −5) pode ser repetido: durante a própria etapa, o algoritmo reaproveita valores já reduzidos e dá mais de uma volta. Por isso nenhum custo que alcance o ciclo é bem definido. No teste 6, a matriz não muda, mas a rota A → F muda. Com `A→B = 1`, a aresta direta já empata com A → C → B, que também custa 1. Como a atualização exige melhora **estrita**, `prox[A][B]` continua B.

**19. Floyd-Warshall × Bellman-Ford**

Bellman-Ford resolve o caminho mínimo de **uma origem** para todos os destinos, em O(|V|·|E|). Floyd-Warshall resolve **todos os pares** de uma vez, em O(|V|³). Para ter todos os pares com Bellman-Ford, seria preciso rodá-lo |V| vezes, O(|V|²·|E|), o que em grafo denso chega a O(|V|⁴). Os dois aceitam pesos negativos e detectam ciclos negativos. Bellman-Ford só enxerga os alcançáveis a partir da origem; Floyd-Warshall vê qualquer ciclo, pela diagonal. A linha A da matriz final é exatamente o resultado do Bellman-Ford a partir de A (ver desafio).

**20. Por que dist[v][v] < 0 indica ciclo negativo?**

`dist[v][v]` começa em 0, o caminho vazio. Ele só muda se, em alguma etapa k, `dist[v][k] + dist[k][v] < 0`, ou seja, se existe um passeio que sai de v, passa por k e volta a v com soma negativa. Um passeio fechado de peso negativo contém um ciclo de peso negativo. Na outra direção, se v está num ciclo negativo, quando k for o último vértice do ciclo a ser considerado, os dois trechos v → k e k → v já estarão calculados, e a soma deles deixa `dist[v][v]` negativo.

**21. Complexidade**

**Tempo Θ(|V|³)**: três laços aninhados de |V| iterações (aqui, 6³ = 216 testes de relaxamento), independente do número de arestas. **Espaço Θ(|V|²)** para as matrizes `dist` e `prox`. A versão do programa guarda uma cópia da matriz por etapa só para exibir o rastreamento; o algoritmo em si precisa só das duas matrizes.

## Parte VII – Desafio

`challenge` lista as rotas mínimas a partir de A com `APSP.Path(A, v)` para cada destino e compara `dist[A][v]` com o `distancia[v]` de uma execução de Bellman-Ford a partir de A no mesmo grafo.

```go
// challenge lista todas as rotas mínimas a partir de A e compara com Bellman-Ford.
func challenge(res shortest.APSP) error {
    r.Title("PARTE VII – Desafio: rotas a partir de A × Bellman-Ford")
    bf, err := shortest.BellmanFord(F+1, arestas, A)
    if err != nil {
        return err
    }
    fmt.Printf("  %-8s %-26s %-14s %s\n", "Destino", "Rota (prox)", "Floyd-Warshall", "Bellman-Ford")
    for v := range F + 1 {
        ok := "="
        if res.Dist[A][v] != bf.Dist[v] {
            ok = "≠"
        }
        fmt.Printf("  %-8s %-26s %-14s %s %s\n", r.Label(v), r.Path(res.Path(A, v)), r.Dist(res.Dist[A][v]), r.Dist(bf.Dist[v]), ok)
    }
    return nil
}
```

```text
  Destino  Rota (prox)                Floyd-Warshall Bellman-Ford
  A        A                          0              0 =
  B        A → C → B                  1              1 =
  C        A → C                      2              2 =
  D        A → C → B → D              3              3 =
  E        A → C → B → D → E          1              1 =
  F        A → C → B → D → E → F      4              4 =
```

**Os custos coincidem** porque os dois algoritmos calculam a mesma grandeza: o peso do caminho mínimo δ(A, v). Sem ciclo negativo, esse valor existe e é único. Bellman-Ford chega a ele relaxando arestas até o ponto fixo; Floyd-Warshall, liberando intermediários. A linha A da matriz final é o vetor de distâncias do Bellman-Ford a partir de A. O teste `TestFloydWarshallEstudo` confere essa igualdade.

## Entrega

- **Código:** `cmd/unidade2/floyd_warshall/main.go` (programa principal) e `internal/shortest/floyd_warshall.go` (criar_matrizes, floyd_warshall e reconstruir_caminho).
- **Comentários:** inicialização, intermediários, relaxamento, detecção de ciclo e reconstrução estão marcados no código.
- **Matrizes:** a inicial, a de cada intermediário e a final estão na saída de `go run ./cmd/unidade2/floyd_warshall`; D⁽ᴮ⁾ e a final estão preenchidas no enunciado.
- **Testes automatizados:** `go test ./internal/shortest/` confere D⁽ᴮ⁾, a matriz final, as rotas e o ciclo negativo.
- Nenhuma implementação pronta de Floyd-Warshall foi usada.
