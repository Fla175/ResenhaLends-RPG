// Value Object que encapsula uma sugestão de correção estatística.

#ifndef RPG_DOMAIN_VALUE_OBJECTS_SUGESTAO_CORRECAO_H
#define RPG_DOMAIN_VALUE_OBJECTS_SUGESTAO_CORRECAO_H

#include <string>
#include <optional>

namespace rpg {
namespace domain {
namespace value_objects {

class SugestaoCorrecao {
public:
    SugestaoCorrecao(Texto original, Texto sugerida, double confianca)
        : original_(std::move(original)),
          sugerida_(std::move(sugerida)),
          confianca_(confianca) {}

    const Texto& original() const { return original_; }
    const Texto& sugerida() const { return sugerida_; }
    double confianca() const { return confianca_; }

    bool deve_aplicar(double limiar = 0.5) const {
        return confianca_ >= limiar;
    }

private:
    Texto original_;
    Texto sugerida_;
    double confianca_;
};

} // namespace value_objects
} // namespace domain
} // namespace rpg

#endif // RPG_DOMAIN_VALUE_OBJECTS_SUGESTAO_CORRECAO_H
