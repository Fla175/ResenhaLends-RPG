// Caso de uso: recuperação local de contexto para aumento de geração (RAG).

#ifndef RPG_APPLICATION_USE_CASES_RECUPERAR_CONTEXTO_RAG_H
#define RPG_APPLICATION_USE_CASES_RECUPERAR_CONTEXTO_RAG_H

#include <string>
#include <vector>
#include <memory>
#include "../../application/ports/irepositorio_contexto.h"
#include "../../domain/value_objects/fragmento_contexto.h"

namespace rpg {
namespace application {
namespace use_cases {

class RecuperarContextoRAG {
public:
    explicit RecuperarContextoRAG(std::shared_ptr<ports::IRepositorioContexto> repositorio)
        : repositorio_(std::move(repositorio)) {}

    std::vector<domain::value_objects::FragmentoContexto> executar(
        const std::string& consulta, std::size_t limite = 5);

private:
    std::shared_ptr<ports::IRepositorioContexto> repositorio_;
};

} // namespace use_cases
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_USE_CASES_RECUPERAR_CONTEXTO_RAG_H
