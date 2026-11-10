#include "gerar_resposta_mestre.h"
#include <sstream>

namespace rpg {
namespace application {
namespace use_cases {

std::string GerarRespostaMestre::executar(std::int64_t sessao_id,
                                          const std::string& entrada_jogador) {
    // Recupera o aggregate root da sessão.
    auto sessao_opt = repositorio_sessao_->find_by_id(sessao_id);
    if (!sessao_opt.has_value()) {
        return "Sessão não encontrada.";
    }

    const auto& sessao = sessao_opt.value();

    // Recupera contexto relevante via RAG (porta abstrata).
    auto fragmentos = repositorio_contexto_->buscar_fragmentos(entrada_jogador, 3);

    // Monta o prompt para o modelo A.L.I.C.E.
    std::ostringstream prompt;
    prompt << "Você é o mestre de mesa da sessão \"" << sessao.nome_mestre() << "\".\n";
    prompt << "Cena atual: ";
    if (sessao.cena_atual().has_value()) {
        prompt << sessao.cena_atual().value();
    } else {
        prompt << "indefinida";
    }
    prompt << "\n\nContexto recuperado:\n";
    for (const auto& frag : fragmentos) {
        prompt << "- [" << frag.fonte() << "] " << frag.texto() << "\n";
    }
    prompt << "\nJogador diz: " << entrada_jogador << "\n";
    prompt << "Mestre:";

    // Gera resposta via porta de inferência (adaptador externo: LibTorch/llama.cpp).
    return inferencia_->gerar_resposta(prompt.str());
}

} // namespace use_cases
} // namespace application
} // namespace rpg
