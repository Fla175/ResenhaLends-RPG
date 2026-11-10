// Domain Event: disparado quando uma correção estatística é aplicada a um texto.

#ifndef RPG_DOMAIN_EVENTS_CORRECAO_APLICADA_H
#define RPG_DOMAIN_EVENTS_CORRECAO_APLICADA_H

#include <cstdint>
#include <string>
#include "../../domain/value_objects/sugestao_correcao.h"

namespace rpg {
namespace domain {
namespace events {

class CorrecaoAplicada {
public:
    CorrecaoAplicada(std::int64_t sessao_id, std::string texto_original,
                     const value_objects::SugestaoCorrecao& sugestao)
        : sessao_id_(sessao_id),
          texto_original_(std::move(texto_original)),
          sugestao_(sugestao) {}

    std::int64_t sessao_id() const { return sessao_id_; }
    const std::string& texto_original() const { return texto_original_; }
    const value_objects::SugestaoCorrecao& sugestao() const { return sugestao_; }

private:
    std::int64_t sessao_id_;
    std::string texto_original_;
    value_objects::SugestaoCorrecao sugestao_;
};

} // namespace events
} // namespace domain
} // namespace rpg

#endif // RPG_DOMAIN_EVENTS_CORRECAO_APLICADA_H
