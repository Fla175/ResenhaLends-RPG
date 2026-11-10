// Porta conduzida para recuperação de fragmentos de contexto (RAG).

#ifndef RPG_APPLICATION_PORTS_IREPOSITORIO_CONTEXTO_H
#define RPG_APPLICATION_PORTS_IREPOSITORIO_CONTEXTO_H

#include <string>
#include <vector>
#include <memory>
#include "../../domain/value_objects/fragmento_contexto.h"

namespace rpg {
namespace application {
namespace ports {

class IRepositorioContexto {
public:
    virtual ~IRepositorioContexto() = default;

    virtual std::vector<domain::value_objects::FragmentoContexto> buscar_fragmentos(
        const std::string& consulta, std::size_t limite = 5) = 0;
};

} // namespace ports
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_PORTS_IREPOSITORIO_CONTEXTO_H
