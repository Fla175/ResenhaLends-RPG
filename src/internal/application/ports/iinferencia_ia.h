// Porta conduzida para inferência do modelo A.L.I.C.E. (LibTorch ou llama.cpp).

#ifndef RPG_APPLICATION_PORTS_IINFERENCIA_IA_H
#define RPG_APPLICATION_PORTS_IINFERENCIA_IA_H

#include <string>
#include <memory>

namespace rpg {
namespace application {
namespace ports {

class IInferenciaIA {
public:
    virtual ~IInferenciaIA() = default;

    // Gera uma resposta de texto a partir de um prompt.
    virtual std::string gerar_resposta(const std::string& prompt) = 0;

    // Carrega o modelo A.L.I.C.E. de um caminho no disco.
    virtual bool carregar_modelo(const std::string& caminho_modelo) = 0;
};

} // namespace ports
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_PORTS_IINFERENCIA_IA_H
