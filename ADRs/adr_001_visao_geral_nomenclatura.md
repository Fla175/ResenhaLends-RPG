# ADR 001: Visão Geral, Nomenclatura e Versionamento do Projeto

## Status
Aprovado

## Contexto
O ADR original (ADRs/adr_projeto_rpg.md) foi modularizado em documentos independentes, um para cada etapa arquitetural do projeto. Esta ADR serve como índice e define a nomenclatura de versionamento padrão adotada por todo o ecossistema ResenhaLends.

O projeto adota padrões de arquitetura de backend comprovados — Clean Architecture, Hexagonal Architecture (Ports and Adapters) e Domain-Driven Design (DDD) — para garantir camadas com regras de dependência claras, definições de interface e fronteiras de teste bem definidas.

## Princípios Arquiteturais

### Clean Architecture
As dependências fluem apenas para dentro. Camadas internas (entidades e casos de uso) não conhecem camadas externas (frameworks, banco de dados, entrega). Cada fronteira de camada é cruzada por uma interface abstrata, permitindo que os casos de uso sejam testados sem banco de dados, UI ou serviços externos.

### Hexagonal Architecture (Ports and Adapters)
O núcleo de domínio fica no centro, livre de frameworks. Ports são interfaces abstratas que definem como o núcleo interage com o mundo exterior (driving e driven). Adapters são implementações concretas desses ports (repositório em banco, adaptador REST, etc.). Isso permite trocar implementações sem tocar no núcleo e usar adapters em memória nos testes.

### Domain-Driven Design (DDD)
* **Bounded Contexts**: cada subdomínio tem um modelo coerente e isolado; evitamos compartilhar um único modelo por todo o sistema.
* **Ubiquitous Language**: cada termo no código corresponde ao termo usado pelos especialistas de domínio.
* **Táticas**: Entidades, Value Objects, Aggregates, Repositories e Domain Events guiam o design tático dentro de cada contexto.

## Nomenclatura de Versionamento

O versionamento segue o padrão:

**C<temporada>.<itens>**

* **<temporada>**: Número da temporada de atualização. Atualmente estamos na **Temporada 1**. Quando o projeto for retomado daqui a alguns anos, passará para a Temporada 2, e assim por diante.
* **<itens>**: Número acumulado de correções ou trabalhos entregues. Exemplo: se foram implementadas 5 funcionalidades e 10 patches sobre elas, o número será 15, resultando em **C1.15**.

### Exemplo
* Primeira entrega da Temporada 1: `C1.1`
* Após 15 entregas/correções na Temporada 1: `C1.15`
* Primeira entrega da Temporada 2: `C2.1`

## Mapa de ADRs Modularizadas

| ADR | Etapa | Foco | Padrão Arquitetural |
|-----|-------|------|---------------------|
| ADR 001 | 0 | Visão Geral, Nomenclatura e Versionamento | Clean Architecture + DDD (índice e bounded contexts) |
| ADR 002 | 1 | Núcleo e Motor Interno (`src/internal/`) + D&D nativo + plug-ins de sistemas de RPG | Clean Architecture + Hexagonal (Ports and Adapters) |
| ADR 003 | 2 | Harness de IA L.I.A.N.D.R.A. + modelo A.L.I.C.E. (70M parâmetros) | Clean Architecture + Bounded Context de IA |
| ADR 004 | 3 | Banco de Dados Kekse (grafo não-relacional) + `schema.kekse` + plug-ins Lua | Hexagonal (Driven Adapter de persistência) |

## Consequências
* Facilita a leitura e manutenção isolada de cada camada.
* Permite versionamento independente por temporada de atualização.
* Mantém o rastreamento de decisões por domínio (núcleo, IA, banco de dados).
* Isola regras de negócio de frameworks e tecnologias externas.
* Permite testar casos de uso com adapters em memória, sem banco de dados nem Docker.
* Define bounded contexts claros, evitando vazamento de modelo entre subdomínios.
