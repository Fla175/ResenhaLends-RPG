# ADR 002: Núcleo e Motor Interno

## Status
Aprovado

## Etapa
1 — Núcleo e motor interno

## Contexto
O projeto requer um núcleo de jogo de RPG portátil, escrito em C++ Moderno, que venha nativamente configurado para o sistema D&D e suporte a importação de diferentes sistemas de RPG por meio de plug-ins. O núcleo deve seguir padrões de arquitetura limpa, hexagonal e orientada a domínio, com interconexão compacta e estável com o sistema de plug-ins, sem dependência excessiva de rede ou processamento externo.

## Decisões Tecnológicas e Arquiteturais

### 1.1 Clean Architecture
* **Decisão:** O núcleo adota Clean Architecture com dependências fluindo exclusivamente para dentro.
* **Camadas:**
  * **Entities (Domínio):** Modelos de negócio do RPG (personagens, itens, atributos, regras de D&D), sem imports de frameworks, banco de dados ou UI.
  * **Use Cases (Casos de Uso):** Regras de aplicação do RPG (criar personagem, rolar dado, iniciar combate, executar ação de NPC), orquestram entidades e não conhecem adaptadores externos.
  * **Interface Adapters:** Controladores, presenters, gateways e adaptadores de plug-in que traduzem entre casos de uso e formatos externos (Lua, Kekse, UI).
  * **Frameworks & Drivers:** UI (Epsilon/Delta/Omega), banco Kekse, interpretador Lua, rede — todos na camada mais externa.

### 1.2 Hexagonal Architecture (Ports and Adapters)
* **Decisão:** O núcleo adota Ports and Adapters para isolar o domínio de tecnologias externas.
* **Portas (abstratas, definidas no domínio):**
  * `ICharacterRepository` — persistência e reconstituição de personagens.
  * `IPluginLoader` — carregamento e descarregamento de plug-ins de sistemas de RPG.
  * `IDiceEngine` — rolagem de dados abstraída.
  * `IRenderer` — saída para interface do usuário.
  * `IBrain` — algoritmo de decisão para NPCs e monstros.
* **Adapters (concretos, nas bordas):**
  * `KekseCharacterRepository` — implementa `ICharacterRepository` via `extern "C"` com `graph_db_*`.
  * `LuaPluginLoader` — implementa `IPluginLoader` carregando plug-ins Lua com `schema.kekse`.
  * `EpsilonRenderer` / `DeltaRenderer` / `OmegaRenderer` — implementam `IRenderer`.
  * `StandardDiceEngine` — implementa `IDiceEngine`.
  * `NpcBrainAdapter` — implementa `IBrain` usando n-grams e cadeias de Markov.

### 1.3 Domain-Driven Design (DDD)
* **Decisão:** O domínio do RPG aplica padrões táticos de DDD.
* **Bounded Contexts:**
  * **Contexto Regras D&D:** Modelo nativo de atributos, classes, raças, magias e combate de D&D.
  * **Contexto Plug-ins de Sistemas:** Modelo genérico para sistemas de RPG externos (Ordem Paranormal, Tormenta, etc.), com anti-corruption layer para traduzir conceitos entre contextos.
  * **Contexto Persistência Kekse:** Abstração de acesso ao banco de dados grafo não-relacional.
  * **Contexto Interface:** Abstração de renderização e entrada do usuário.
* **Padrões Táticos:**
  * **Entidades:** `Personagem`, `NPC`, `Monstro`, `Item`, `Magia` — com identidade estável.
  * **Objetos de Valor:** `Atributo` (0–100), `Personalidade` (conjunto de traços 0–100), `Dano`, `Modificador`, `Rolagem`.
  * **Agregados:** `Personagem` é raiz do agregado; itens e efeitos só são acessados através dele.
  * **Repositórios:** `ICharacterRepository` abstrai persistência; implementação concreta no adapter Kekse.
  * **Eventos de Domínio:** `PersonagemCriado`, `AcaoNPCExecutada`, `DanoAplicado` — para coordenação entre agregados.

### 1.4 Núcleo Interno (`src/internal/`) e Configuração Nativa D&D
* **Decisão:** O núcleo interno reside em `src/internal/` e vem nativamente configurado para **D&D**.
* **Estrutura:**
  * `src/internal/domain/` — entidades, objetos de valor, interfaces de portas.
  * `src/internal/application/` — casos de uso.
  * `src/internal/adapters/` — adaptadores de infraestrutura (Kekse, Lua, UI).
  * `src/internal/frameworks/` — drivers externos, main, CMake.
* **Justificativa:** Fornece experiência pronta para uso imediato com D&D, sem configuração externa.

### 1.5 Cérebro de NPCs e Monstros
* **Decisão:** NPCs e monstros possuem um algoritmo interno de "cérebro" (`IBrain`) que gera comportamento autônomo.
* **Atributos de Personalidade:** Cada NPC/monstro possui traços de personalidade definidos em escala de **0 a 100** (ex: agressividade, coragem, sociabilidade, inteligência, impulsividade).
* **Geração de Características:** Características textuais e comportamentais são geradas por modelos estatísticos de **n-grams** e **cadeias de Markov**, implementados nativamente em C++ (sem Python).
* **Integração:** O cérebro é injetado como dependência no domínio via porta `IBrain`, permitindo substituição por diferentes algoritmos sem alterar entidades.

### 1.6 Sistema de Plug-ins e Compatibilidade com o Núcleo
* **Decisão:** O núcleo internal e o sistema de plug-ins devem trabalhar em conjunto com **compatibilidade extrema**.
* **Plug-ins:** Escritos em **Lua**, localizados em `src/internal/plugins/`, cada um acompanhado de um arquivo `schema.kekse`.
* **Interconexão:** O núcleo expõe portas estáveis e versionadas que os plug-ins consomem. O schema.kekse padrão vem embutido no núcleo e é configurado nativamente para D&D.
* **Estabilidade:** A interconexão é feita de forma compacta, com mínimo overhead de rede e processamento, preferindo chamadas in-process e FFI leve.

### 1.7 Integração com Banco de Dados Kekse
* **Decisão:** A persistência ocorre por trás da porta `ICharacterRepository`, sem chamadas diretas `extern "C"` espalhadas pela lógica de negócio.
* **Implementação:** Um adapter concreto em C/C++ faz a ponte com `graph_db_*` via `extern "C"` (Zero-cost FFI), mantendo o domínio livre de detalhes de armazenamento.

### 1.8 Restrições de Rede e Processamento
* **Decisão:** O projeto utiliza o mínimo possível de rede e processamento externo.
* **Justificativa:** Toda a interconexão entre núcleo e plug-ins é in-process, compacta e estável, priorizando execução local e baixo consumo de recursos.

## Consequências
* **Positivas:** Domínio isolado e testável sem banco ou UI; plug-ins intercambiáveis; cérebro de NPCs configurável; compatibilidade extrema entre núcleo e sistemas externos; baixo consumo de recursos.
* **Negativas:** Complexidade inicial maior devido a camadas e interfaces; necessidade de manter anti-corruption layers entre bounded contexts; esforço para implementar n-grams e Markov nativos em C++.

## Versão
C1.1
