# ADR 003: Harness de IA L.I.A.N.D.R.A. e Modelo A.L.I.C.E.

## Status
Aprovado

## Etapa
2 — Harness de Inteligência Artificial

## Contexto
O projeto necessita de um harness de IA embutido localmente, sem dependência de Python em tempo de execução, capaz de atuar como mestre de mesa, corrigir gramática estatisticamente e operar em dispositivos com recursos limitados. Além disso, o código deve ser mantível, testável e escalável, com regras de dependência claras entre camadas, evitando que lógica de negócio vaze para infraestrutura ou frameworks externos.

## Decisões Arquiteturais

### 3.1 Padrão Arquitetural: Clean Architecture e Hexagonal (Ports and Adapters)
* **Decisão:** O harness de IA seguirá **Clean Architecture** com dependências fluindo exclusivamente para dentro, implementado através de **Hexagonal Architecture (Ports and Adapters)**.
* **Camadas (do interno para o externo):**
  1. **Domain Core (Entidades e Regras de Negócio):** Modelo de domínio puro, sem imports de frameworks, bancos de dados ou bibliotecas de inferência. Contém entidades de mestragem, value objects de correção gramatical e regras de negócio centrais.
  2. **Use Cases (Casos de Uso):** Regras de aplicação que orquestram entidades do domínio (ex: gerar resposta de mestre de mesa, corrigir texto, recuperar contexto via RAG). Não conhecem adaptadores concretos.
  3. **Ports (Interfaces):** Abstrações que definem como o núcleo interage com o mundo externo. Portas dirigidas (driving) expõem casos de uso; portas conduzidas (driven) definem contratos para repositórios, inferência de IA e persistência.
  4. **Interface Adapters:** Controladores, presenters e gateways que traduzem entre use cases e formatos externos (ex: entrada de usuário, saída de texto, chamadas ao modelo A.L.I.C.E.).
  5. **Frameworks & Drivers:** LibTorch, llama.cpp, persistência em arquivos binários, Android NDK, UI — todos na camada mais externa, atrás de portas abstratas.
* **Regra de dependência:** Dependências apontam apenas para dentro. Camadas internas não sabem nada sobre camadas externas. Cada fronteira de camada é cruzada via interface abstrata.

### 3.2 Domain-Driven Design (DDD)
* **Decisão:** Aplicar padrões estratégicos e táticos de DDD dentro do Domain Core.
* **Bounded Contexts:**
  * **Contexto de Mestragem:** Modelo coerente para sessões de RPG, personagens, cenas e eventos de narrativa.
  * **Contexto de Correção Linguística:** Modelo para correção estatística de gramática, incluindo N-Gram e Cadeias de Markov.
  * **Contexto de Recuperação de Contexto (RAG):** Modelo para recuperação local e geração aumentada de contexto.
* **Ubiquitous Language:** Termos como "mestre de mesa", "sessão", "correção estatística" e "fragmento de contexto" são usados consistentemente no código, documentação e comunicação com especialistas de domínio.
* **Padrões Táticos:**
  * **Aggregates:** Consistência delimitada; apenas o root é acessível de fora (ex: `Sessao` como aggregate root).
  * **Value Objects:** Objetos imutáveis identificados por atributos (ex: `Texto`, `SugestaoCorrecao`).
  * **Repositórios:** Portas abstratas para persistir e reconstituir aggregates, desacopladas do mecanismo de armazenamento.
  * **Domain Events:** Capturam ocorrências dentro do domínio (ex: `ContextoRecuperado`, `CorrecaoAplicada`) para coordenação entre bounded contexts.

### 3.3 Componentes L.I.A.N.D.R.A. (reposicionados na arquitetura)
* **RAG Híbrido:** Implementado como Use Case que orquestra portas de recuperação local e geração aumentada para contexto consistente de mestragem.
* **Organização de Dados M.E.U.:** Estrutura interna de metadados, entidades e relações mapeada para entidades e value objects do Domain Core, acessada via portas de repositório.
* **Correção Ortográfica:** Implementação nativa em C++ de modelos **N-Gram** e **Cadeias de Markov** como adaptadores concretos da porta de correção gramatical (`ICorrecaoGramatical`).

### 3.4 Modelo A.L.I.C.E.
* **Decisão:** Desenvolvimento de um modelo de IA do zero com aproximadamente **70 milhões de parâmetros**, batizado de **A.L.I.C.E.**
* **Justificativa:** Tamanho adequado (~40MB a ~140MB) para execução direta em CPU/GPU de dispositivos modestos, sem a pesada stack do ecossistema Python.

### 3.5 Backend de Execução
* **Decisão:** Utilização de **LibTorch** ou **llama.cpp** para inferência local.
* **Justificativa:** Portabilidade via Android NDK, GCC, Clang e MSVC; execução eficiente em hardware limitado. Implementado como adaptador externo (driver) atrás da porta de inferência.

### 3.6 Restrição de Linguagem
* **Decisão:** Rejeição total de ambientes Python em tempo de execução no produto final.
* **Justificativa:** Reduz consumo de memória, tamanho do binário e superfície de ataque.

## Limites de Teste
* **Decisão:** Cada Use Case deve ser exequível em teste de unidade puro, sem banco de dados real, Docker ou rede, utilizando adaptadores em memória (in-memory adapters) que implementem as portas do domínio.
* **Justificativa:** Garante que a lógica de negócio seja testável independentemente de UI, banco de dados ou serviços externos, conforme preconizado pela Clean Architecture.

## Consequências
* **Positivas:** Execução local, baixa latência, controle total sobre o modelo, correção gramatical estatística leve, lógica de negócio independente de frameworks, testes de unidade sem dependências de infraestrutura, possibilidade de trocar LibTorch por llama.cpp (ou vice-versa) sem alterar o núcleo de domínio.
* **Negativas:** Esforço inicial para definir camadas, portas e interfaces; indireção adicional entre camadas; manutenção de implementações nativas de N-Gram e Markov; esforço para treinar o modelo do zero.

## Versão
C1.1
