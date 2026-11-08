# ADR 004: Banco de Dados Kekse

## Status
Aprovado

## Etapa
3 — Banco de Dados

## Contexto
O projeto exige um banco de dados próprio, não-relacional, baseado em grafos, escrito em C99 puro, com suporte a múltiplos sistemas de RPG e persistência em arquivos binários.

## Decisões Tecnológicas e Arquiteturais

### 3.1 Nome e Estrutura
* **Decisão:** O banco de dados será batizado de **Kekse**.
* **Modelo:** Não-relacional baseado em grafos (nós e arestas).
* **Linguagem:** C99 puro, com controle manual de ponteiros.

### 3.2 Suporte a Sistemas de RPG
* **Decisão:** Suporte nativo ao registro de diferentes sistemas de RPG, com foco em compatibilidade de atributos, itens e novas tabelas de dados.
* **Exemplo — Ordem Paranormal:** Deve conter uma nova tabela para tipos de personagens elementais (como sangue e outros), controlada pelos plug-ins.

### 3.3 Plug-ins Lua e Schemas
* **Decisão:** Os diferentes sistemas de RPG serão controlados por plug-ins escritos em **Lua**.
* **Arquivo de Esquema:** Cada plug-in deve vir acompanhado de um arquivo chamado **schema.kekse**.
* **Schema Padrão:** O arquivo `schema.kekse` padrão do projeto já vem embutido no núcleo e é nativamente configurado para **D&D**.

### 3.4 Persistência
* **Decisão:** Persistência em arquivos binários.
* **Justificativa:** Garante estrutura ultra-enxuta e fundação pedagógica sobre persistência de dados de baixo nível.

## Consequências
* **Positivas:** Consumo mínimo de memória, portabilidade extrema, isolamento de sistemas de RPG via Lua, sandboxing possível ao desativar bibliotecas `os` e `io` do Lua.
* **Negativas:** Gerenciamento manual de memória; necessidade de implementar do zero operações de busca e gravação em grafo.

## Versão
C1.0
