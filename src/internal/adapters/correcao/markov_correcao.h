// Adaptador concreto de correção ortográfica baseado em Cadeias de Markov.
// Implementa a porta ICorrecaoGramatical (ADR 003, seção 3.3).

#ifndef RPG_ADAPTERS_CORRECAO_MARKOV_CORRECAO_H
#define RPG_ADAPTERS_CORRECAO_MARKOV_CORRECAO_H

#include <string>
#include <memory>
#include <unordered_map>
#include "../../application/ports/icorrecao_gramatical.h"

namespace rpg {
namespace adapters {
namespace correcao {

class MarkovCorrecao : public application::ports::ICorrecaoGramatical {
public:
    MarkovCorrecao();

    std::unique_ptr<domain::value_objects::SugestaoCorrecao> corrigir(
        const std::string& texto) override;

private:
    std::unordered_map<std::string, std::unordered_map<std::string, double>> transicoes_;
};

} // namespace correcao
} // namespace adapters
} // namespace rpg

#endif // RPG_ADAPTERS_CORRECAO_MARKOV_CORRECAO_H
