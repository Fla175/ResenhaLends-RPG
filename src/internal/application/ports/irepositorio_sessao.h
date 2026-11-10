// Porta conduzida para persistência e reconstituição do aggregate Sessao.

#ifndef RPG_APPLICATION_PORTS_IREPOSITORIO_SESSAO_H
#define RPG_APPLICATION_PORTS_IREPOSITORIO_SESSAO_H

#include <cstdint>
#include <memory>
#include <optional>
#include "../../domain/entities/sessao.h"

namespace rpg {
namespace application {
namespace ports {

class IRepositorioSessao {
public:
    virtual ~IRepositorioSessao() = default;

    virtual std::int64_t salvar(const domain::entities::Sessao& sessao) = 0;
    virtual std::optional<domain::entities::Sessao> find_by_id(std::int64_t id) = 0;
    virtual bool remover(std::int64_t id) = 0;
};

} // namespace ports
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_PORTS_IREPOSITORIO_SESSAO_H
