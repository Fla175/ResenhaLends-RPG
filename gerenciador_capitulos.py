# gerenciador_capitulos.py
from datetime import datetime

import gerenciador_json
from why_man.cronica import fechar_capitulo, gerar_contexto_inicial

LIMITE = {"pequena": 3, "media": 4, "medio": 4, "extensa": 6}


def _agora():
    return datetime.now().strftime("%Y-%m-%d %H:%M")


def listar(nome_campanha):
    return gerenciador_json.carregar(nome_campanha, ".capitulos.json") or []


def atual(nome_campanha):
    return gerenciador_json.carregar(nome_campanha, ".capitulo_atual.json") or {
        "numero": 1,
        "eventos": [],
        "aberto": True,
    }


def _gravar_atual(nome_campanha, cap):
    gerenciador_json.salvar(nome_campanha, ".capitulo_atual.json", cap)


def garantir_abertura(nome_campanha):
    """Gera o briefing inicial do mestre a partir da descrição e o prólogo."""
    ctx = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    caps = listar(nome_campanha)
    if not ctx.get("contexto_inicial"):
        texto = gerar_contexto_inicial(
            ctx.get("descricao", ""),
            ctx.get("objetivo_principal", ""),
            ctx.get("nome_campanha", nome_campanha),
        )
        ctx["contexto_inicial"] = texto
        gerenciador_json.salvar(nome_campanha, ".context.json", ctx)
    if not any(c.get("tipo") == "prologo" for c in caps):
        caps.insert(0, {
            "numero": 0,
            "titulo": "Prólogo — o mestre abre o mundo",
            "resumo": ctx.get("contexto_inicial", ""),
            "eventos": [ctx.get("descricao", ""), ctx.get("objetivo_principal", "")],
            "tipo": "prologo",
            "fechado_em": _agora(),
        })
        gerenciador_json.salvar(nome_campanha, ".capitulos.json", caps)
    if gerenciador_json.carregar(nome_campanha, ".capitulo_atual.json") is None:
        hist = gerenciador_json.carregar(nome_campanha, ".hist.json") or []
        _gravar_atual(nome_campanha, {
            "numero": 1,
            "eventos": list(hist),
            "aberto": True,
        })
    return ctx.get("contexto_inicial", "")


def anotar(nome_campanha, desfecho, narracao=""):
    cap = atual(nome_campanha)
    cap.setdefault("eventos", []).append(desfecho if not narracao else f"{desfecho} // {narracao}")
    _gravar_atual(nome_campanha, cap)
    ctx = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    limite = LIMITE.get(str(ctx.get("tamanho_campanha", "media")).lower(), 4)
    if len(cap.get("eventos") or []) >= limite:
        return encerrar(nome_campanha)
    return None


def encerrar(nome_campanha):
    cap = atual(nome_campanha)
    eventos = cap.get("eventos") or []
    if not eventos:
        return None
    ctx = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    fechado = fechar_capitulo(
        int(cap.get("numero") or 1),
        eventos,
        ctx.get("objetivo_principal", "a jornada"),
    )
    fechado["fechado_em"] = _agora()
    caps = listar(nome_campanha)
    caps.append(fechado)
    gerenciador_json.salvar(nome_campanha, ".capitulos.json", caps)
    _gravar_atual(nome_campanha, {
        "numero": int(fechado["numero"]) + 1,
        "eventos": [],
        "aberto": True,
    })
    return fechado
