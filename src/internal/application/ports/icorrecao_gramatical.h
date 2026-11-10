// Porta conduzida (driven) para correção estatística de gramática.
// Implementada por adaptadores N-Gram e Cadeias de Markov (ADR 003, seção 3.3).

#ifndef RPG_APPLICATION_PORTS_ICORRECAO_GRAMATICAL_H
#define RPG_APPLICATION_PORTS_ICORRECAO_GRAMATICAL_H

#include <string>
#include <memory>
#include "../../domain/value_objects/sugestao_correcao.h"

namespace rpg {
namespace application {
namespace ports {

class ICorrecaoGramatical {
public:
    virtual ~ICorrecaoGramatical() = default;

    // Retorna uma sugestão de correção para o texto informado.
    virtual std::unique_ptr<domain::value_objects::SugestaoCorrecao> corrigir(
        const std::string& texto) = 0;
};

} // namespace ports
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_PORTS_ICORRECAO_GRAMATICAL_H
