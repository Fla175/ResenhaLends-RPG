#ifndef RPG_INTERNAL_DOMAIN_ENTITIES_PERSONAGEM_H
#define RPG_INTERNAL_DOMAIN_ENTITIES_PERSONAGEM_H

#include <cstdint>
#include <string>
#include <utility>
#include "../value_objects/atributo.h"
#include "../value_objects/personalidade.h"

namespace rpg {
namespace domain {
namespace entities {

class Personagem {
public:
    Personagem(std::int64_t id, std::string nome)
        : id_(id), nome_(std::move(nome)) {}

    std::int64_t id() const { return id_; }
    const std::string& nome() const { return nome_; }

    void set_forca(const rpg::domain::value_objects::Atributo& forca) { forca_ = forca; }
    rpg::domain::value_objects::Atributo forca() const { return forca_; }

    void set_personalidade(const rpg::domain::value_objects::Personalidade& personalidade) { personalidade_ = personalidade; }
    const rpg::domain::value_objects::Personalidade& personalidade() const { return personalidade_; }

private:
    std::int64_t id_;
    std::string nome_;
    rpg::domain::value_objects::Atributo forca_{0};
    rpg::domain::value_objects::Personalidade personalidade_;
};

} // namespace entities
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_ENTITIES_PERSONAGEM_H
