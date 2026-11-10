// Caso de uso: correção estatística de gramática de um texto arbitrário.

#ifndef RPG_APPLICATION_USE_CASES_CORRIGIR_TEXTO_H
#define RPG_APPLICATION_USE_CASES_CORRIGIR_TEXTO_H

#include <string>
#include <memory>
#include "../../application/ports/icorrecao_gramatical.h"

namespace rpg {
namespace application {
namespace use_cases {

class CorrigirTexto {
public:
    explicit CorrigirTexto(std::shared_ptr<ports::ICorrecaoGramatical> corretor)
        : corretor_(std::move(corretor)) {}

    // Retorna a sugestão de correção (pode ser nula se não houver correção).
    std::unique_ptr<domain::value_objects::SugestaoCorrecao> executar(
        const std::string& texto);

private:
    std::shared_ptr<ports::ICorrecaoGramatical> corretor_;
};

} // namespace use_cases
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_USE_CASES_CORRIGIR_TEXTO_H
