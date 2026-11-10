// Value Object imutável que representa um trecho de texto.
// Conforme ADR 003: identificados por atributos, sem identidade própria.

#ifndef RPG_DOMAIN_VALUE_OBJECTS_TEXTO_H
#define RPG_DOMAIN_VALUE_OBJECTS_TEXTO_H

#include <string>

namespace rpg {
namespace domain {
namespace value_objects {

class Texto {
public:
    explicit Texto(std::string conteudo)
        : conteudo_(std::move(conteudo)) {}

    const std::string& conteudo() const { return conteudo_; }
    std::size_t tamanho() const { return conteudo_.size(); }

    bool operator==(const Texto& other) const { return conteudo_ == other.conteudo_; }
    bool operator!=(const Texto& other) const { return !(*this == other); }

private:
    std::string conteudo_;
};

} // namespace value_objects
} // namespace domain
} // namespace rpg

#endif // RPG_DOMAIN_VALUE_OBJECTS_TEXTO_H
