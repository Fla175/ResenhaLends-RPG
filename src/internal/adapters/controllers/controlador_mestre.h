// Controlador da Interface Adapters: traduz entre entrada do usuário e o caso de uso.
// Não conhece frameworks de UI específicos; depende apenas de IRenderer (porta).

#ifndef RPG_ADAPTERS_CONTROLE_CONTROLADOR_MESTRE_H
#define RPG_ADAPTERS_CONTROLE_CONTROLADOR_MESTRE_H

#include <string>
#include <memory>
#include "../../application/use_cases/gerar_resposta_mestre.h"
#include "../../domain/ports/irenderer.h"

namespace rpg {
namespace adapters {
namespace controllers {

class ControladorMestre {
public:
    ControladorMestre(
        std::shared_ptr<application::use_cases::GerarRespostaMestre> use_case,
        std::shared_ptr<domain::ports::IRenderer> renderer)
        : use_case_(std::move(use_case)),
          renderer_(std::move(renderer)) {}

    void processar_entrada(std::int64_t sessao_id, const std::string& entrada);

private:
    std::shared_ptr<application::use_cases::GerarRespostaMestre> use_case_;
    std::shared_ptr<domain::ports::IRenderer> renderer_;
};

} // namespace controllers
} // namespace adapters
} // namespace rpg

#endif // RPG_ADAPTERS_CONTROLE_CONTROLADOR_MESTRE_H
