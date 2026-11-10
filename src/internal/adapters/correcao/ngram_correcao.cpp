#include "ngram_correcao.h"
#include <sstream>
#include <cmath>

namespace rpg {
namespace adapters {
namespace correcao {

NgramCorrecao::NgramCorrecao() {
    // Inicialização estática de um pequeno modelo N-Gram de exemplo.
    // Em produção, isso seria carregado de um arquivo binário ou treinado.
    ngramas_["a casa"] = 100;
    ngramas_["a casa e"] = 80;
    ngramas_["a casa está"] = 60;
    ngramas_["o jogador"] = 90;
    ngramas_["o jogador rola"] = 70;
}

std::unique_ptr<domain::value_objects::SugestaoCorrecao> NgramCorrecao::corrigir(
    const std::string& texto) {
    // Implementação de stub: retorna o mesmo texto com confiança baixa.
    // A lógica estatística completa de N-Gram seria implementada aqui.
    auto sugestao = std::make_unique<domain::value_objects::SugestaoCorrecao>(
        domain::value_objects::Texto(texto),
        domain::value_objects::Texto(texto),
        0.1);
    return sugestao;
}

} // namespace correcao
} // namespace adapters
} // namespace rpg
