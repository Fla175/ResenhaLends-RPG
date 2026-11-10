# ADR 001: Visão Geral, Nomenclatura e Versionamento do Projeto

## Status
Aprovado

## Contexto
O ADR original (ADRs/adr_projeto_rpg.md) foi modularizado em documentos independentes, um para cada etapa arquitetural do projeto. Esta ADR serve como índice e define a nomenclatura de versionamento padrão adotada por todo o ecossistema ResenhaLends.

O projeto adota padrões de arquitetura de backend comprovados — Clean Architecture, Hexagonal Architecture (Ports and Adapters) e Domain-Driven Design (DDD) — para garantir camadas com regras de dependência claras, definições de interface e fronteiras de teste bem definidas.

## Princípios Arquiteturais

### Clean Architecture
As dependências fluem apenas para dentro. Camadas internas (entidades e casos de uso) não conhecem camadas externas (frameworks, banco de dados, entrega). Cada fronteira de camada é cruzada por uma interface abstrata, permitindo que os casos de uso sejam testados sem banco de dados, UI ou serviços externos.

**Camadas (fluxo de dependência para dentro):**

1. **Entities**: Modelos de negócio centrais, sem imports de frameworks.
2. **Use Cases**: Regras de negócio da aplicação; orquestram entidades.
3. **Interface Adapters**: Controladores, presenters, gateways — traduzem entre casos de uso e formatos externos.
4. **Frameworks & Drivers**: UI, banco de dados, serviços externos — todos no anel mais externo.

**Princípios-chave:**
- Dependências apontam apenas para dentro; camadas internas não sabem nada sobre camadas externas.
- Lógica de negócio é independente de frameworks, bancos de dados e mecanismos de entrega.
- Toda fronteira de camada é cruzada por uma interface abstrata.
- Testável sem UI, banco de dados ou serviços externos.

### Hexagonal Architecture (Ports and Adapters)
O núcleo de domínio fica no centro, livre de frameworks.

**Componentes:**
- **Domain Core**: Lógica de negócio vive aqui, sem frameworks.
- **Ports**: Interfaces abstratas que definem como o núcleo interage com o mundo exterior (ports driving — entrada, e ports driven — saída).
- **Adapters**: Implementações concretas dos ports (repositório em banco, adaptador REST, etc.).

**Benefícios:**
- Trocar implementações sem tocar no núcleo (ex.: substituir persistência sem alterar regras de negócio).
- Usar adapters em memória nos testes — sem Docker necessário.
- Decisões de tecnologia adiadas para as bordas do sistema.

### Domain-Driven Design (DDD)

**Padrões Estratégicos:**
- **Bounded Contexts**: Cada subdomínio tem um modelo coerente e isolado; evitamos compartilhar um único modelo por todo o sistema.
- **Context Mapping**: Define como os contextos se relacionam (Anti-Corruption Layer, Shared Kernel, Open Host Service).
- **Ubiquitous Language**: Cada termo no código corresponde ao termo usado pelos especialistas de domínio.

**Padrões Táticos:**
- **Entidades**: Objetos com identidade estável que mudam ao longo do tempo.
- **Value Objects**: Objetos imutáveis identificados por seus atributos.
- **Aggregates**: Limites de consistência; apenas a raiz é acessível de fora.
- **Repositories**: Persistem e reconstituem aggregates; abstraem o mecanismo de armazenamento.
- **Domain Events**: Capturam coisas que aconteceram dentro do domínio; usados para coordenação entre aggregates.

### Regras de Dependência e Fronteiras
- Dependências apontam apenas para dentro.
- Camadas internas (entidades e casos de uso) não conhecem camadas externas (frameworks, banco de dados, entrega).
- Cada fronteira de camada é cruzada por uma interface abstrata.
- Casos de uso importam apenas do domínio (entidades e interfaces); nunca importam de adaptadores ou infraestrutura.
- Isso evita ciclos de dependência e vazamento de infraestrutura para o domínio.

### Estratégia de Testes
A marca de uma Clean Architecture corretamente aplicada é que todo caso de uso pode ser exercitado em um teste de unidade simples, sem banco de dados real, sem Docker e sem rede:
- Use cases são testados com adapters em memória (in-memory) injetados através das portas abstratas.
- Nenhuma chamada de banco de dados vaza para a camada de aplicação.
- Testes de unidade não requerem infraestrutura externa.

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
* Impede ciclos de dependência entre camadas por meio de interfaces abstratas nas fronteiras.
