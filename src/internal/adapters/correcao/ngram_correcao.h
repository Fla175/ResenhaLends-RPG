// Adaptador concreto de correção ortográfica baseado em N-Gram.
// Implementa a porta ICorrecaoGramatical (ADR 003, seção 3.3).

#ifndef RPG_ADAPTERS_CORRECAO_NGRAM_CORRECAO_H
#define RPG_ADAPTERS_CORRECAO_NGRAM_CORRECAO_H

#include <string>
#include <memory>
#include <unordered_map>
#include "../../application/ports/icorrecao_gramatical.h"

namespace rpg {
namespace adapters {
namespace correcao {

class NgramCorrecao : public application::ports::ICorrecaoGramatical {
public:
    NgramCorrecao();

    std::unique_ptr<domain::value_objects::SugestaoCorrecao> corrigir(
        const std::string& texto) override;

private:
    // Modelo N-Gram simplificado: mapa de sequências para contagens.
    std::unordered_map<std::string, std::int64_t> ngramas_;
};

} // namespace correcao
} // namespace adapters
} // namespace rpg

#endif // RPG_ADAPTERS_CORRECAO_NGRAM_CORRECAO_H
