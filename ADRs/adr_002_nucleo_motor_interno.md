# ADR 002: Núcleo e Motor Interno

## Status
Aprovado

## Etapa
1 — Núcleo e motor interno

## Contexto
O projeto requer um núcleo de jogo de RPG portátil, escrito em C++ Moderno, que venha nativamente configurado para o sistema D&D e suporte a importação de diferentes sistemas de RPG por meio de plug-ins.

## Decisões Tecnológicas e Arquiteturais

### 1.1 Modelo Arquitetural
* **Decisão:** Monólito Modular Conectado, totalmente in-process.
* **Justificativa:** Elimina overhead de rede e latência, ideal para hardwares limitados.

### 1.2 Linguagem e Padrão
* **Decisão:** C++ Moderno (C++17/C++20).
* **Justificativa:** RAII, smart pointers (`std::unique_ptr`, `std::shared_ptr`) e containers seguros. Proibido uso de C++98 com `new`/`delete` manuais.

### 1.3 Configuração Nativa
* **Decisão:** O núcleo vem nativamente configurado para o sistema **D&D** (Dungeons & Dragons).
* **Justificativa:** Fornece uma experiência pronta para uso imediato, sem necessidade de configuração externa para o sistema padrão.

### 1.4 Plug-ins de Sistemas de RPG
* **Decisão:** Suporte a plug-ins para diferentes sistemas de RPG, localizados em `src/internal/`.
* **Justificativa:** Permite expansão para sistemas como Ordem Paranormal, Tormenta, etc., sem modificar o núcleo principal.

### 1.5 Integração com Banco de Dados
* **Decisão:** Integração com o banco de dados Kekse via diretivas `extern "C"` (Zero-cost FFI).
* **Justificativa:** Controle manual de ponteiros e persistência em arquivos binários, com custo zero de abstração entre C++ e C.

## Consequências
* **Positivas:** Portabilidade máxima, desempenho nativo, fácil expansão via plug-ins.
* **Negativas:** Gerenciamento manual de memória requer disciplina rigorosa; build multiplataforma via CMake demanda atenção.

## Versão
C1.0
