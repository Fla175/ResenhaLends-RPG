// Aggregate root do Contexto de Mestragem (Bounded Context).
// Conforme ADR 003: apenas o root é acessível de fora do aggregate.

#ifndef RPG_DOMAIN_ENTITIES_SESSAO_H
#define RPG_DOMAIN_ENTITIES_SESSAO_H

#include <cstdint>
#include <string>
#include <vector>
#include <memory>
#include <optional>

#include "../../domain/entities/personagem.h"

namespace rpg {
namespace domain {
namespace entities {

class Sessao {
public:
    using Id = std::int64_t;

    Sessao(Id id, std::string nome_mestre)
        : id_(id), nome_mestre_(std::move(nome_mestre)) {}

    Id id() const { return id_; }
    const std::string& nome_mestre() const { return nome_mestre_; }

    void adicionar_personagem(std::shared_ptr<Personagem> personagem) {
        personagens_.push_back(std::move(personagem));
    }

    const std::vector<std::shared_ptr<Personagem>>& personagens() const {
        return personagens_;
    }

    void adicionar_evento_narrativo(std::string evento) {
        eventos_narrativos_.push_back(std::move(evento));
    }

    const std::vector<std::string>& eventos_narrativos() const {
        return eventos_narrativos_;
    }

    void definir_cena_atual(std::string cena) {
        cena_atual_ = std::move(cena);
    }

    const std::optional<std::string>& cena_atual() const {
        return cena_atual_;
    }

private:
    Id id_;
    std::string nome_mestre_;
    std::vector<std::shared_ptr<Personagem>> personagens_;
    std::vector<std::string> eventos_narrativos_;
    std::optional<std::string> cena_atual_;
};

} // namespace entities
} // namespace domain
} // namespace rpg

#endif // RPG_DOMAIN_ENTITIES_SESSAO_H
