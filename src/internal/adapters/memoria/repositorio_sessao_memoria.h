// Adaptador em memória do repositório de sessões.
// Usado em testes de unidade puro, sem banco de dados real (ADR 003, Limites de Teste).

#ifndef RPG_ADAPTERS_MEMORIA_REPOSITORIO_SESSAO_MEMORIA_H
#define RPG_ADAPTERS_MEMORIA_REPOSITORIO_SESSAO_MEMORIA_H

#include <cstdint>
#include <memory>
#include <optional>
#include <unordered_map>
#include <mutex>
#include "../../domain/entities/sessao.h"
#include "../../application/ports/irepositorio_sessao.h"

namespace rpg {
namespace adapters {
namespace memoria {

class RepositorioSessaoMemoria : public application::ports::IRepositorioSessao {
public:
    std::int64_t salvar(const domain::entities::Sessao& sessao) override;
    std::optional<domain::entities::Sessao> find_by_id(std::int64_t id) override;
    bool remover(std::int64_t id) override;

private:
    std::unordered_map<std::int64_t, domain::entities::Sessao> sessoes_;
    std::int64_t proximo_id_ = 1;
    mutable std::mutex mutex_;
};

} // namespace memoria
} // namespace adapters
} // namespace rpg

#endif // RPG_ADAPTERS_MEMORIA_REPOSITORIO_SESSAO_MEMORIA_H
