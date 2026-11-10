// Caso de uso: gerar resposta do mestre de mesa para uma sessão de RPG.
// Orquestra entidades do domínio e portas abstratas, sem conhecer adaptadores concretos.

#ifndef RPG_APPLICATION_USE_CASES_GERAR_RESPOSTA_MESTRE_H
#define RPG_APPLICATION_USE_CASES_GERAR_RESPOSTA_MESTRE_H

#include <cstdint>
#include <string>
#include <memory>
#include "../../domain/entities/sessao.h"
#include "../../application/ports/irepositorio_sessao.h"
#include "../../application/ports/iinferencia_ia.h"
#include "../../application/ports/irepositorio_contexto.h"

namespace rpg {
namespace application {
namespace use_cases {

class GerarRespostaMestre {
public:
    GerarRespostaMestre(
        std::shared_ptr<ports::IRepositorioSessao> repositorio_sessao,
        std::shared_ptr<ports::IInferenciaIA> inferencia,
        std::shared_ptr<ports::IRepositorioContexto> repositorio_contexto)
        : repositorio_sessao_(std::move(repositorio_sessao)),
          inferencia_(std::move(inferencia)),
          repositorio_contexto_(std::move(repositorio_contexto)) {}

    // Executa o caso de uso. Retorna a resposta gerada pelo mestre.
    std::string executar(std::int64_t sessao_id, const std::string& entrada_jogador);

private:
    std::shared_ptr<ports::IRepositorioSessao> repositorio_sessao_;
    std::shared_ptr<ports::IInferenciaIA> inferencia_;
    std::shared_ptr<ports::IRepositorioContexto> repositorio_contexto_;
};

} // namespace use_cases
} // namespace application
} // namespace rpg

#endif // RPG_APPLICATION_USE_CASES_GERAR_RESPOSTA_MESTRE_H
