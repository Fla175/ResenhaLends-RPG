# ADR 001: Visão Geral, Nomenclatura e Versionamento do Projeto

## Status
Aprovado

## Contexto
O ADR original (ADRs/adr_projeto_rpg.md) foi modularizado em documentos independentes, um para cada etapa arquitetural do projeto. Esta ADR serve como índice e define a nomenclatura de versionamento padrão adotada por todo o ecossistema ResenhaLends.

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

| ADR | Etapa | Foco |
|-----|-------|------|
| ADR 002 | 1 | Núcleo e Motor Interno (`src/internal/`) + D&D nativo + plug-ins de sistemas de RPG |
| ADR 003 | 2 | Harness de IA L.I.A.N.D.R.A. + modelo A.L.I.C.E. (70M parâmetros) |
| ADR 004 | 3 | Banco de Dados Kekse (grafo não-relacional) + `schema.kekse` + plug-ins Lua |

## Consequências
* Facilita a leitura e manutenção isolada de cada camada.
* Permite versionamento independente por temporada de atualização.
* Mantém o rastreamento de decisões por domínio (núcleo, IA, banco de dados).
