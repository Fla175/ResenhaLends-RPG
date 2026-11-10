# ADR 004: Banco de Dados Kekse

## Status
Aprovado (Atualizado C2.0 — Segurança, Robustez e Desempenho)

## Etapa
3 — Banco de Dados

## Contexto
O projeto exige um banco de dados próprio, não-relacional, baseado em grafos, escrito em C99 puro, com suporte a múltiplos sistemas de RPG e persistência em arquivos binários. Além de ser ultra-rápido e enxuto (foco principal), o banco deve ser **seguro**, **robusto** e **otimizado**, resistindo a entradas malformadas, falhas de disco/energia e cargas concorrentes sem corromper dados ou travar o núcleo.

## Decisão

### 4.1 Padrão Arquitetural: Hexagonal (Ports and Adapters) com Clean Architecture
O núcleo do banco de dados **Kekse** seguirá uma arquitetura hexagonal, onde as dependências fluem para dentro. O core de domínio é independente de frameworks, bancos de dados e mecanismos de persistência. Toda interação com o mundo externo ocorre através de portas (interfaces abstratas) e adaptadores concretos.

### 4.2 Camadas e Regras de Dependência

| Camada | Responsabilidades | Pode depender de |
|--------|-------------------|------------------|
| **Domain Core** | Entidades de RPG (Personagem, Item, Atributo), Value Objects, Aggregates, regras de negócio | Nenhuma camada interna |
| **Use Cases** | Orquestração de casos de uso (criar nó, adicionar aresta, consultar grafo) | Domain Core, interfaces de portas |
| **Interface Adapters** | Implementações de portas: repositórios, controladores de plug-ins Lua, presenters, camada de segurança | Domain Core, Use Cases |
| **Frameworks & Drivers** | Runtime Lua, sistema de arquivos, serialização binária, C99, WAL, mmap | Interface Adapters |

**Regra:** Dependências apontam apenas para dentro. O Domain Core não sabe nada sobre Lua, arquivos binários, persistência ou segurança.

### 4.3 Portas e Adaptadores

#### Portas (Driven Ports)
- `IGraphDatabasePort`: operações de grafo (adicionar nó, adicionar aresta, buscar nós, flush, begin/commit/rollback).
- `IPluginPort`: registro e execução de plug-ins de sistemas de RPG.
- `ISecurityPort` (nova): validação de entrada, sanitização de identificadores, controle de acesso e criptografia opcional.
- `IReliabilityPort` (nova): detecção de corrupção, recuperação de falhas, checksums e WAL.

#### Adaptadores
- **BinaryFileAdapter**: implementa `IGraphDatabasePort` sobre persistência em arquivos binários (C99), com páginas alinhadas, WAL e checksums.
- **LuaPluginAdapter**: implementa `IPluginPort`, carregando plug-ins Lua em sandbox (sem `os`/`io`) e arquivos `schema.kekse`.
- **InMemoryGraphAdapter** (testes): implementa `IGraphDatabasePort` em memória, permitindo testar use cases sem banco real.
- **SecurityAdapter**: implementa `ISecurityPort`, aplicando validação de entrada, limites de tamanho e sanitização.
- **ReliabilityAdapter**: implementa `IReliabilityPort`, gerenciando WAL, fsync, checksums CRC64 e recuperação de desastre.

### 4.4 Bounded Contexts (DDD)
Cada sistema de RPG (ex: D&D, Ordem Paranormal) é tratado como um **bounded context** com seu próprio modelo de domínio, linguagem ubíqua e esquema (`schema.kekse`). O núcleo do Kekse atua como um contexto compartilhado de persistência de grafos, enquanto os plug-ins Lua isolam as regras específicas de cada sistema, evitando vazamento de modelo entre contextos.

### 4.5 Suporte a Sistemas de RPG
* **Decisão:** Suporte nativo ao registro de diferentes sistemas de RPG, com foco em compatibilidade de atributos, itens e novas tabelas de dados.
* **Exemplo — Ordem Paranormal:** Deve conter uma nova tabela para tipos de personagens elementais (como sangue e outros), controlada pelos plug-ins.

### 4.6 Plug-ins Lua e Schemas
* **Decisão:** Os diferentes sistemas de RPG serão controlados por plug-ins escritos em **Lua**.
* **Arquivo de Esquema:** Cada plug-in deve vir acompanhado de um arquivo chamado **schema.kekse**.
* **Schema Padrão:** O arquivo `schema.kekse` padrão do projeto já vem embutido no núcleo e é nativamente configurado para **D&D**.
* **Sandbox Lua:** Os plug-ins serão executados em ambiente restrito. Bibliotecas perigosas (`os`, `io`, `package`, `dofile`, `loadfile`) serão removidas ou encapsuladas. A execução será limitada por tempo e memória para evitar negação de serviço.

### 4.7 Persistência
* **Decisão:** Persistência em arquivos binários.
* **Justificativa:** Garante estrutura ultra-enxuta e fundação pedagógica sobre persistência de dados de baixo nível. A camada de persistência é um adaptador atrás da porta `IGraphDatabasePort`, podendo ser trocada por um adaptador em memória ou outra tecnologia sem alterar o domínio.
* **Melhorias de desempenho:**
  - Armazenamento por **páginas fixas** (4 KiB), alinhadas ao tamanho de página do SO.
  - **Índice em memória** (hash aberto ou B+ leve) para mapeamento rápido de IDs para offsets de página.
  - **Cache LRU de páginas** para minimizar leituras de disco.
  - **mmap opcional** para leitura direta de arquivos grandes, evitando cópias desnecessárias (zero-copy).
  - **Operações em lote (batch):** múltiplas inserções/arestas em uma única transação para reduzir chamadas de fsync.
  - **Alocador de memória pool/arena** para nós e arestas, evitando fragmentação e chamadas frequentes a `malloc`/`free`.
  - Estruturas de dados **cache-friendly** (arrays contíguos, pouca indireção de ponteiros).

### 4.8 Segurança
* **Validação de entrada:** Todos os ponteiros de string, tamanhos e IDs são validados antes de qualquer gravação. Uso exclusivo de `snprintf`, `strncpy_s` (ou equivalente seguro) e funções com limite de tamanho.
* **Sanitização de identificadores:** Nomes de nós, arestas e rótulos são verificados contra caracteres perigosos (`/`, `..`, `\0`, caracteres de controle) para prevenir path traversal e injeção de nomes de arquivo.
* **Controle de acesso:** A porta `ISecurityPort` permite camadas de autenticação/autorização opcionais no adaptador, sem expor essa complexidade ao Domain Core.
* **Criptografia opcional:** O `BinaryFileAdapter` pode ser configurado para criptografar páginas em repouso (ex: AES-256 ou XOR simples para cenários pedagógicos), controlado pelo adaptador, não pelo core.
* **Limites rígidos:** Quantidade máxima de nós, arestas, tamanho de propriedade e profundidade de recursão são limitados para evitar exaustão de memória.

### 4.9 Robustez e Recuperação
* **WAL (Write-Ahead Log):** Antes de modificar páginas principais, alterações são gravadas em um log antecipado. Após commit e fsync, as páginas são aplicadas. Em caso de falha, o WAL é replayado na inicialização.
* **Checksums:** Cada página e cada registro do WAL contém um checksum CRC64. Leitura com checksum inválido é rejeitada e marcada como corrompida.
* **Transações:** Suporte a `begin`, `commit` e `rollback`. Operações não confirmadas são descartadas.
* **Recuperação automática:** Na abertura, se o WAL estiver incompleto ou houver páginas sujas, o banco tenta recuperar ao estado consistente mais recente.
* **Tratamento de erros:** Nenhuma operação usa `abort` ou `exit`. Todas retornam códigos de erro explícitos (`int`/`bool`) para que o chamador decida a estratégia.
* **Gerenciamento de memória:** Uso de pool de memória com detecção de vazamentos em modo debug. Nenhum ponteiro bruto é desreferenciado sem validação de bounds.

### 4.10 Desempenho e Otimização (Foco Principal)
* **Latência mínima:** Caminho crítico (busca por ID) deve ser O(1) ou O(log n) com índice em memória.
* **Throughput:** Inserções em lote devem saturar a largura de banda do disco sequencial, não o acesso aleatório.
* **Concorrência:** Fine-grained locking (mutex por shard de página) ou estruturas lock-free para leituras concorrentes. Leitores não bloqueiam leitores.
* **Alocação zero-cópia:** Quando possível, dados são expostos diretamente do buffer de página sem cópia intermediária.
* **Métricas:** O banco expõe contadores internos (hits/misses de cache, páginas lidas/gravadas, tempo médio de operação) para tuning contínuo.

## Consequências
* **Positivas:** Consumo mínimo de memória, portabilidade extrema, isolamento de sistemas de RPG via Lua, sandboxing possível ao desativar bibliotecas `os` e `io` do Lua, testabilidade de use cases com adaptadores em memória (sem banco real), possibilidade de trocar persistência binária por outra tecnologia sem afetar o domínio, resistência a corrupção via WAL/checksums, proteção contra entradas maliciosas, desempenho previsível e escalável.
* **Negativas:** Gerenciamento manual de memória (C99); necessidade de implementar do zero operações de busca e gravação em grafo; complexidade adicional de manter fronteiras entre camadas; overhead de WAL e checksums em gravações (mitigado por batch e fsync assíncrono opcional).

## Versão
C2.0
