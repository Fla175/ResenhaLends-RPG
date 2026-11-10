#include "controlador_mestre.h"

namespace rpg {
namespace adapters {
namespace controllers {

void ControladorMestre::processar_entrada(std::int64_t sessao_id,
                                          const std::string& entrada) {
    std::string resposta = use_case_->executar(sessao_id, entrada);
    renderer_->render_text(resposta);
}

} // namespace controllers
} // namespace adapters
} // namespace rpg
