#include "markov_correcao.h"

namespace rpg {
namespace adapters {
namespace correcao {

MarkovCorrecao::MarkovCorrecao() {
    // Inicialização de stub de cadeia de Markov.
}

std::unique_ptr<domain::value_objects::SugestaoCorrecao> MarkovCorrecao::corrigir(
    const std::string& texto) {
    auto sugestao = std::make_unique<domain::value_objects::SugestaoCorrecao>(
        domain::value_objects::Texto(texto),
        domain::value_objects::Texto(texto),
        0.15);
    return sugestao;
}

} // namespace correcao
} // namespace adapters
} // namespace rpg
