// Value Object que representa um fragmento de contexto recuperado via RAG.

#ifndef RPG_DOMAIN_VALUE_OBJECTS_FRAGMENTO_CONTEXTO_H
#define RPG_DOMAIN_VALUE_OBJECTS_FRAGMENTO_CONTEXTO_H

#include <string>
#include <optional>

namespace rpg {
namespace domain {
namespace value_objects {

class FragmentoContexto {
public:
    FragmentoContexto(std::string texto, std::string fonte, double relevancia)
        : texto_(std::move(texto)),
          fonte_(std::move(fonte)),
          relevancia_(relevancia) {}

    const std::string& texto() const { return texto_; }
    const std::string& fonte() const { return fonte_; }
    double relevancia() const { return relevancia_; }

private:
    std::string texto_;
    std::string fonte_;
    double relevancia_;
};

} // namespace value_objects
} // namespace domain
} // namespace rpg

#endif // RPG_DOMAIN_VALUE_OBJECTS_FRAGMENTO_CONTEXTO_H
