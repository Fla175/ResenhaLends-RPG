// Domain Event: disparado quando um fragmento de contexto é recuperado.
// Conforme ADR 003: coordenação entre bounded contexts.

#ifndef RPG_DOMAIN_EVENTS_CONTEXTO_RECUPERADO_H
#define RPG_DOMAIN_EVENTS_CONTEXTO_RECUPERADO_H

#include <cstdint>
#include <string>
#include <vector>
#include "../../domain/value_objects/fragmento_contexto.h"

namespace rpg {
namespace domain {
namespace events {

class ContextoRecuperado {
public:
    ContextoRecuperado(std::int64_t sessao_id, std::string consulta,
                       std::vector<value_objects::FragmentoContexto> fragmentos)
        : sessao_id_(sessao_id),
          consulta_(std::move(consulta)),
          fragmentos_(std::move(fragmentos)) {}

    std::int64_t sessao_id() const { return sessao_id_; }
    const std::string& consulta() const { return consulta_; }
    const std::vector<value_objects::FragmentoContexto>& fragmentos() const { return fragmentos_; }

private:
    std::int64_t sessao_id_;
    std::string consulta_;
    std::vector<value_objects::FragmentoContexto> fragmentos_;
};

} // namespace events
} // namespace domain
} // namespace rpg

#endif // RPG_DOMAIN_EVENTS_CONTEXTO_RECUPERADO_H
