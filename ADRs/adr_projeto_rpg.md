# ADR 001: Definição da Arquitetura do Núcleo, Banco de Dados e Interfaces para o Projeto de RPG

## Status
Aprovado

## Contexto
O projeto consiste em um jogo de RPG projetado sob a filosofia de portabilidade extrema, com o objetivo de rodar desde computadores modestos até dispositivos móveis antigos (celulares legados), sem impor complexidade inviável ao desenvolvimento. O sistema exige um banco de dados de grafos próprio, um harness de Inteligência Artificial embutido localmente (modelo de mestre de mesa de ~70M de parâmetros), correção gramatical estatística em tempo real (N-Grams/Markov) e suporte a módulos dinâmicos para a importação de regras de diferentes sistemas de RPG (BYOK e pacotes externos).

## Decisões Tecnológicas e Arquiteturais

### 1. Modelo Arquitetural Geral
* **Decisão:** Adotar um modelo de **Monólito Modular Conectado**, totalmente embutido em processo (**In-Process**).
* **Justificativa:** Elimina o overhead de latência, o consumo de memória e a complexidade de rede que APIs internas (como HTTP ou gRPC) causariam em hardwares severamente limitados.

### 2. Camada de Banco de Dados
* **Decisão:** Desenvolvimento de um banco de dados próprio, não-relacional baseado em grafos (nós e arestas) escrito em **C99 puro**.
* **Justificativa:** Garantia de uma estrutura ultra-enxuta, controle manual e preciso de ponteiros, além de atuar como uma fundação pedagógica profunda sobre persistência de dados em arquivos binários.

### 3. Camada do Núcleo (Core Engine) e Harness de IA
* **Decisão:** Escrita do núcleo em **C++ Moderno (C++17/C++20)**, integrando nativamente o banco de dados via diretivas `extern "C"` (Zero-cost FFI). O Harness de IA local utilizará a arquitetura C++ do **LibTorch** ou **llama.cpp**, rejeitando o uso de ambientes Python em tempo de execução no produto final. A correção gramatical será feita via implementações nativas em C++ de modelos N-Gram e Cadeias de Markov.
* **Justificativa:** O C++ oferece portabilidade máxima (via Android NDK, GCC, Clang e MSVC), controle absoluto sobre o ciclo de vida da memória e a capacidade de executar o modelo de 70M (~40MB a ~140MB) diretamente na CPU/GPU do dispositivo sem a pesada stack do ecossistema Python. O suporte a chaves externas (BYOK) será exposto de forma nativa.

### 4. Camadas de Interface (Views)
* **Decisão:** Total desacoplamento do motor visual em relação à lógica de negócio do jogo. O ecossistema será compilado em 4 variantes independentes:
  * **ResenhaLends Omega (GUI):** Renderização de mapas procedimentais usando as bibliotecas leves `raylib` e `Dear ImGui`.
  * **ResenhaLends Delta (TUI Estilizada):** Interface de terminal moderna e declarativa baseada na biblioteca `FTXUI`.
  * **ResenhaLends Iota (TUI Lite):** Interface de terminal limpa com remoção de recursos pesados via `ncurses`.
  * **ResenhaLends Epsilon (TUI Ultra-Lite):** Otimização máxima voltada para o Termux e terminais antigos usando buffers diretos em `std::cout` e sequências ANSI.
* **Justificativa:** Permite que o mesmo núcleo lógico atenda a perfis de hardware e preferências de usuário drasticamente diferentes sem comprometer a estabilidade do motor principal.

### 5. Sistema de Expansibilidade e Pacotes Remotos
* **Decisão:** Suporte à importação de sistemas de RPG (além do padrão D&D) através de scripts escritos na linguagem **Lua**, hospedados em repositórios do GitHub. O download será gerenciado no Core por meio da biblioteca **libcurl**. A integração e execução dos scripts Lua dentro do C++ serão feitas pela biblioteca header-only **Sol3 (Sol2)**.
* **Justificativa:** Scripts Lua são incrivelmente leves e portáveis. O isolamento do ambiente de execução (*sandboxing*) será garantido desativando o carregamento das bibliotecas de sistema (`os`) e entrada/saída (`io`) do Lua no Core, mitigando riscos de segurança ao executar códigos de terceiros.

## Consequências
* **Positivas:** Performance bruta máxima por ciclo de clock; consumo mínimo e previsível de memória RAM; total isolamento entre interface, IA e armazenamento; ecossistema seguro de plugins; altíssimo ganho de aprendizado prático sobre computação de baixo nível.
* **Negativas:** Gerenciamento manual da complexidade do build multiplataforma (CMake); necessidade de atenção rigorosa à gerência de memória no C++ moderno para evitar vazamentos (*memory leaks*); esforço inicial acentuado por programar componentes estruturais (banco de dados e corretores) totalmente à mão.

---

## 📖 Apêndice: Roteiro de Estudos e Fontes de C++ Moderno

Para o desenvolvimento correto deste projeto, é terminantemente proibido o uso de C++ antigo (padrão C++98 baseado em `new`/`delete` manuais). O aprendizado deve ser focado nos pilares modernos de RAII, smart pointers (`std::unique_ptr`, `std::shared_ptr`) e containers seguros.

### Fontes Textuais e de Referência
1. **[LearnCpp](https://www.learncpp.com/):** Guia linear completo e gratuito para aprender C++ moderno do absoluto zero ao avançado.
2. **[Cppreference](https://en.cppreference.com/w/):** Enciclopédia técnica e oficial da biblioteca padrão (STL), sintaxe e recursos da linguagem.

### Fontes de Aprendizado Visual (Vídeos)
1. **[The Cherno - C++ Series](https://www.youtube.com/playlist?list=PLlrATfHPBB18eXmRGEMeGfOWB-4c7xG4k):** Foco em arquitetura de motores de jogos, performance de baixo nível e depuração de memória no compilador.
2. **[C++ Weekly com Jason Turner](https://www.youtube.com/@mjtrux):** Pílulas de conhecimento focadas em otimização, boas práticas e os novos recursos do C++20/C++23.
