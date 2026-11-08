# ADR 004: Banco de Dados Kekse

## Status
Aprovado

## Etapa
3 — Banco de Dados

## Contexto
O projeto exige um banco de dados próprio, não-relacional, baseado em grafos, escrito em C99 puro, com suporte a múltiplos sistemas de RPG e persistência em arquivos binários.

## Decisão

### 4.1 Padrão Arquitetural: Hexagonal (Ports and Adapters) com Clean Architecture
O núcleo do banco de dados **Kekse** seguirá uma arquitetura hexagonal, onde as dependências fluem para dentro. O core de domínio é independente de frameworks, bancos de dados e mecanismos de persistência. Toda interação com o mundo externo ocorre através de portas (interfaces abstratas) e adaptadores concretos.

### 4.2 Camadas e Regras de Dependência

| Camada | Responsabilidades | Pode depender de |
|--------|-------------------|------------------|
| **Domain Core** | Entidades de RPG (Personagem, Item, Atributo), Value Objects, Aggregates, regras de negócio | Nenhuma camada interna |
| **Use Cases** | Orquestração de casos de uso (criar nó, adicionar aresta, consultar grafo) | Domain Core, interfaces de portas |
| **Interface Adapters** | Implementações de portas: repositórios, controladores de plug-ins Lua, presenters | Domain Core, Use Cases |
| **Frameworks & Drivers** | Runtime Lua, sistema de arquivos, serialização binária, C99 | Interface Adapters |

**Regra:** Dependências apontam apenas para dentro. O Domain Core não sabe nada sobre Lua, arquivos binários ou persistência.

### 4.3 Portas e Adaptadores

#### Portas (Driven Ports)
- `IGraphDatabasePort`: operações de grafo (adicionar nó, adicionar aresta, buscar nós, flush).
- `IPluginPort`: registro e execução de plug-ins de sistemas de RPG.

#### Adaptadores
- **BinaryFileAdapter**: implementa `IGraphDatabasePort` sobre persistência em arquivos binários (C99).
- **LuaPluginAdapter**: implementa `IPluginPort`, carregando plug-ins Lua e arquivos `schema.kekse`.
- **InMemoryGraphAdapter** (testes): implementa `IGraphDatabasePort` em memória, permitindo testar use cases sem banco real.

### 4.4 Bounded Contexts (DDD)
Cada sistema de RPG (ex: D&D, Ordem Paranormal) é tratado como um **bounded context** com seu próprio modelo de domínio, linguagem ubíqua e esquema (`schema.kekse`). O núcleo do Kekse atua como um contexto compartilhado de persistência de grafos, enquanto os plug-ins Lua isolam as regras específicas de cada sistema, evitando vazamento de modelo entre contextos.

### 4.5 Suporte a Sistemas de RPG
* **Decisão:** Suporte nativo ao registro de diferentes sistemas de RPG, com foco em compatibilidade de atributos, itens e novas tabelas de dados.
* **Exemplo — Ordem Paranormal:** Deve conter uma nova tabela para tipos de personagens elementais (como sangue e outros), controlada pelos plug-ins.

### 4.6 Plug-ins Lua e Schemas
* **Decisão:** Os diferentes sistemas de RPG serão controlados por plug-ins escritos em **Lua**.
* **Arquivo de Esquema:** Cada plug-in deve vir acompanhado de um arquivo chamado **schema.kekse**.
* **Schema Padrão:** O arquivo `schema.kekse` padrão do projeto já vem embutido no núcleo e é nativamente configurado para **D&D**.

### 4.7 Persistência
* **Decisão:** Persistência em arquivos binários.
* **Justificativa:** Garante estrutura ultra-enxuta e fundação pedagógica sobre persistência de dados de baixo nível. A camada de persistência é um adaptador atrás da porta `IGraphDatabasePort`, podendo ser trocada por um adaptador em memória ou outra tecnologia sem alterar o core.

## Consequências
* **Positivas:** Consumo mínimo de memória, portabilidade extrema, isolamento de sistemas de RPG via Lua, sandboxing possível ao desativar bibliotecas `os` e `io` do Lua, testabilidade de use cases com adaptadores em memória (sem banco real), possibilidade de trocar persistência binária por outra tecnologia sem afetar o domínio.
* **Negativas:** Gerenciamento manual de memória (C99); necessidade de implementar do zero operações de busca e gravação em grafo; complexidade adicional de manter fronteiras entre camadas.

## Versão
C1.0
