# ADR 003: Harness de IA L.I.A.N.D.R.A. e Modelo A.L.I.C.E.

## Status
Aprovado

## Etapa
2 — Harness de Inteligência Artificial

## Contexto
O projeto necessita de um harness de IA embutido localmente, sem dependência de Python em tempo de execução, capaz de atuar como mestre de mesa, corrigir gramática estatisticamente e operar em dispositivos com recursos limitados.

## Decisões Tecnológicas e Arquiteturais

### 2.1 Arquitetura L.I.A.N.D.R.A.
* **Decisão:** O harness de IA seguirá a arquitetura **L.I.A.N.D.R.A.**
* **Componentes:**
  * **RAG Híbrido:** Combinação de recuperação local e geração aumentada para contexto consistente de mestragem.
  * **Organização de Dados M.E.U.:** Estrutura interna de metadados, entidades e relações otimizada para consulta rápida pelo núcleo.
  * **Correção Ortográfica:** Implementação nativa em C++ de modelos **N-Gram** e **Cadeias de Markov** para correção gramatical em tempo real.

### 2.2 Modelo A.L.I.C.E.
* **Decisão:** Desenvolvimento de um modelo de IA do zero com aproximadamente **70 milhões de parâmetros**, batizado de **A.L.I.C.E.**
* **Justificativa:** Tamanho adequado (~40MB a ~140MB) para execução direta em CPU/GPU de dispositivos modestos, sem a pesada stack do ecossistema Python.

### 2.3 Backend de Execução
* **Decisão:** Utilização de **LibTorch** ou **llama.cpp** para inferência local.
* **Justificativa:** Portabilidade via Android NDK, GCC, Clang e MSVC; execução eficiente em hardware limitado.

### 2.4 Restrição de Linguagem
* **Decisão:** Rejeição total de ambientes Python em tempo de execução no produto final.
* **Justificativa:** Reduz consumo de memória, tamanho do binário e superfície de ataque.

## Consequências
* **Positivas:** Execução local, baixa latência, controle total sobre o modelo, correção gramatical estatística leve.
* **Negativas:** Esforço inicial para treinar o modelo do zero; manutenção de implementações nativas de N-Gram e Markov.

## Versão
C1.0
