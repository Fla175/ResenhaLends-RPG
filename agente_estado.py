# agente_estado.py
"""Motor de regras local — zero LLM, feito para notebook de 4 GB."""
import re
import gerenciador_json
import codex

_DIFICULDADE = {
    "fácil": 8,
    "facil": 8,
    "normal": 12,
    "hardcore": 15,
}


def _extrair_nome(acao, jogadores):
    m = re.search(r"jogador\s+(.+?)\s+tenta", acao, re.I)
    if m:
        candidato = m.group(1).strip()
        for j in jogadores:
            if j.get("nome") and j["nome"].lower() == candidato.lower():
                return j["nome"]
    return jogadores[0]["nome"] if jogadores else "Aventureiro"


def _aplicar_vida(jogadores, nome, delta):
    atualizados = []
    for j in jogadores:
        copia = dict(j)
        if copia.get("nome") == nome:
            vida = int(copia.get("vida", 20)) + delta
            copia["vida"] = max(0, min(99, vida))
        atualizados.append(copia)
    return atualizados


def processar_acao(nome_campanha, acao, dado):
    jogadores = gerenciador_json.carregar(nome_campanha, ".jogadores.json") or []
    historico = gerenciador_json.carregar(nome_campanha, ".hist.json") or []
    regras = gerenciador_json.carregar(nome_campanha, ".regras.json") or {}
    contexto = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    monstros = codex.carregar_bestiario(nome_campanha)

    nome_j = _extrair_nome(acao, jogadores)
    texto_acao = acao.split("tenta:", 1)[-1].strip() if "tenta:" in acao.lower() or "tenta:" in acao else acao
    cd = _DIFICULDADE.get(str(regras.get("dificuldade", "Normal")).lower(), 12)
    ameaca = next(iter(monstros.keys()), None)
    objetivo = contexto.get("objetivo_principal", "seguir a jornada")

    if dado == 1:
        novos = _aplicar_vida(jogadores, nome_j, -4)
        desfecho = f"Falha crítica: {nome_j} tenta '{texto_acao}' e se machuca gravemente (-4 HP)."
    elif dado == 20:
        novos = jogadores
        extra = f" A ameaça '{ameaca}' recua." if ameaca else ""
        desfecho = f"Sucesso crítico: {nome_j} realiza '{texto_acao}' de forma lendária e avança rumo a {objetivo}.{extra}"
    elif dado >= cd:
        novos = jogadores
        desfecho = f"Sucesso: {nome_j} consegue '{texto_acao}' (dado {dado} vs CD {cd})."
    elif dado >= cd - 3:
        novos = _aplicar_vida(jogadores, nome_j, -1)
        desfecho = f"Sucesso parcial: {nome_j} tenta '{texto_acao}', paga um preço (-1 HP) e só avança pela metade."
    else:
        dano = 3 if str(regras.get("dificuldade", "")).lower() == "hardcore" else 2
        novos = _aplicar_vida(jogadores, nome_j, -dano)
        desfecho = f"Falha: {nome_j} não consegue '{texto_acao}' (dado {dado} vs CD {cd}) e sofre {dano} de dano."

    gerenciador_json.salvar(nome_campanha, ".jogadores.json", novos)
    historico.append(desfecho)
    gerenciador_json.salvar(nome_campanha, ".hist.json", historico)
    return desfecho
