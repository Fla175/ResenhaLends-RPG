#include "recuperar_contexto_rag.h"

namespace rpg {
namespace application {
namespace use_cases {

std::vector<domain::value_objects::FragmentoContexto> RecuperarContextoRAG::executar(
    const std::string& consulta, std::size_t limite) {
    return repositorio_->buscar_fragmentos(consulta, limite);
}

} // namespace use_cases
} // namespace application
} // namespace rpg
