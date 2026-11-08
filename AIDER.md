# Configuração do Agente Aider

## Política de Branches e Merge

Toda atualização no repositório deve seguir este fluxo:

1. **Branch exclusivo**: Crie sempre um branch dedicado para a atualização (exemplo: `feat/nome-da-feature`, `fix/nome-do-bug`, `docs/nome-da-doc`).
2. **Alterações no branch**: Realize todas as modificações nesse branch exclusivo.
3. **Merge para main**: Após revisão e testes, faça o merge das alterações para a branch `main`.

Não faça commits diretamente na branch `main`. Todas as atualizações devem ser isoladas em branches exclusivos antes de serem mescladas.
