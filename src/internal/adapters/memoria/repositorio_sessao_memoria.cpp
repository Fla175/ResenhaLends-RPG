#include "repositorio_sessao_memoria.h"

namespace rpg {
namespace adapters {
namespace memoria {

std::int64_t RepositorioSessaoMemoria::salvar(const domain::entities::Sessao& sessao) {
    std::lock_guard<std::mutex> lock(mutex_);
    std::int64_t id = sessao.id();
    if (id == 0) {
        id = proximo_id_++;
    }
    sessoes_[id] = sessao;
    return id;
}

std::optional<domain::entities::Sessao> RepositorioSessaoMemoria::find_by_id(std::int64_t id) {
    std::lock_guard<std::mutex> lock(mutex_);
    auto it = sessoes_.find(id);
    if (it == sessoes_.end()) {
        return std::nullopt;
    }
    return it->second;
}

bool RepositorioSessaoMemoria::remover(std::int64_t id) {
    std::lock_guard<std::mutex> lock(mutex_);
    return sessoes_.erase(id) > 0;
}

} // namespace memoria
} // namespace adapters
} // namespace rpg
