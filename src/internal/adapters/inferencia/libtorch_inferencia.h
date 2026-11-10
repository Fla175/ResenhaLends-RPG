// Adaptador externo de inferência baseado em LibTorch.
// Implementa a porta IInferenciaIA (ADR 003, seção 3.5).
// Fica na camada Frameworks & Drivers, atrás da porta abstrata.

#ifndef RPG_ADAPTERS_INFERENCIA_LIBTORCH_INFERENCIA_H
#define RPG_ADAPTERS_INFERENCIA_LIBTORCH_INFERENCIA_H

#include <string>
#include <memory>
#include "../../application/ports/iinferencia_ia.h"

namespace rpg {
namespace adapters {
namespace inferencia {

// Forward declaration do tipo do LibTorch para evitar incluir torch aqui no header.
struct TorchModel;

class LibTorchInferencia : public application::ports::IInferenciaIA {
public:
    LibTorchInferencia();
    ~LibTorchInferencia() override;

    bool carregar_modelo(const std::string& caminho_modelo) override;
    std::string gerar_resposta(const std::string& prompt) override;

private:
    std::unique_ptr<TorchModel> modelo_;
    bool carregado_ = false;
};

} // namespace inferencia
} // namespace adapters
} // namespace rpg

#endif // RPG_ADAPTERS_INFERENCIA_LIBTORCH_INFERENCIA_H
