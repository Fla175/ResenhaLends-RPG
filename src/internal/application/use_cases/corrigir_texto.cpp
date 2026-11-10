#include "corrigir_texto.h"

namespace rpg {
namespace application {
namespace use_cases {

std::unique_ptr<domain::value_objects::SugestaoCorrecao> CorrigirTexto::executar(
    const std::string& texto) {
    return corretor_->corrigir(texto);
}

} // namespace use_cases
} // namespace application
} // namespace rpg
