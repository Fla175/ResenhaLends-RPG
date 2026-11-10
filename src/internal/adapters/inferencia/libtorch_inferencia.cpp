#include "libtorch_inferencia.h"
#include <sstream>

// Nota: Em integração real, incluir <torch/torch.h> e usar LibTorch para
// carregar o modelo A.L.I.C.E. (~70M parâmetros) e executar inferência.
// Aqui mantemos um stub para preservar a independência do núcleo de domínio.

namespace rpg {
namespace adapters {
namespace inferencia {

struct TorchModel {
    std::string caminho;
};

LibTorchInferencia::LibTorchInferencia() = default;

LibTorchInferencia::~LibTorchInferencia() = default;

bool LibTorchInferencia::carregar_modelo(const std::string& caminho_modelo) {
    modelo_ = std::make_unique<TorchModel>();
    modelo_->caminho = caminho_modelo;
    carregado_ = true;
    return true;
}

std::string LibTorchInferencia::gerar_resposta(const std::string& prompt) {
    if (!carregado_) {
        return "Modelo A.L.I.C.E. não carregado.";
    }
    std::ostringstream resposta;
    resposta << "[A.L.I.C.E. via LibTorch] Resposta simulada para: " << prompt;
    return resposta.str();
}

} // namespace inferencia
} // namespace adapters
} // namespace rpg
