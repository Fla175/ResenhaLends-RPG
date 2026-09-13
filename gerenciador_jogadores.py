# gerenciador_jogadores.py
import gerenciador_json
import dados

ATRIBUTOS = ("for", "hab", "res", "arm", "pdf")


def ficha_padrao(nome, **extra):
    ficha = {
        "nome": nome.strip(),
        "classe": extra.get("classe", ""),
        "descricao": extra.get("descricao", ""),
        "raca": extra.get("raca", "Humano"),
        "vida": int(extra.get("vida", 20)),
        "inventario": list(extra.get("inventario") or []),
        "macros": {"1": "", "2": "", "3": "", **(extra.get("macros") or {})},
        "atributos": {"for": 10, "hab": 10, "res": 10, "arm": 10, "pdf": 10, **(extra.get("atributos") or {})},
    }
    return ficha


def mesclar_ficha(original, patch):
    jogador = dict(original or {})
    if not jogador:
        jogador = ficha_padrao(patch.get("nome", "Sem nome"))
    for chave in ("nome", "classe", "descricao", "raca"):
        if chave in patch and patch[chave] is not None:
            jogador[chave] = patch[chave]
    if "vida" in patch and patch["vida"] is not None:
        jogador["vida"] = max(0, min(99, int(patch["vida"])))
    if "inventario" in patch and patch["inventario"] is not None:
        if isinstance(patch["inventario"], str):
            jogador["inventario"] = [i.strip() for i in patch["inventario"].split(",") if i.strip()]
        else:
            jogador["inventario"] = list(patch["inventario"])
    if "macros" in patch and patch["macros"]:
        macros = dict(jogador.get("macros") or {"1": "", "2": "", "3": ""})
        macros.update(patch["macros"])
        jogador["macros"] = macros
    if "atributos" in patch and patch["atributos"]:
        atts = dict(jogador.get("atributos") or {})
        atts.update(patch["atributos"])
        jogador["atributos"] = atts
    return jogador


def listar(nome_campanha):
    return gerenciador_json.carregar(nome_campanha, ".jogadores.json") or []


def _gravar(nome_campanha, jogadores):
    gerenciador_json.salvar(nome_campanha, ".jogadores.json", jogadores)


def adicionar(nome_campanha, ficha, rolar_atts=False):
    jogadores = listar(nome_campanha)
    dados_ficha = dict(ficha or {})
    nome = (dados_ficha.get("nome") or "").strip()
    if not nome:
        raise ValueError("Informe o nome do jogador.")
    if any(j.get("nome", "").lower() == nome.lower() for j in jogadores):
        raise ValueError(f"Já existe um jogador chamado {nome}.")
    if isinstance(dados_ficha.get("inventario"), str):
        dados_ficha["inventario"] = [i.strip() for i in dados_ficha["inventario"].split(",") if i.strip()]
    nova = ficha_padrao(nome, **{k: v for k, v in dados_ficha.items() if k != "nome"})
    if rolar_atts:
        nova["atributos"] = {k: dados.rolar(mostrar=False) for k in ATRIBUTOS}
    jogadores.append(nova)
    _gravar(nome_campanha, jogadores)
    return nova


def atualizar(nome_campanha, nome, patch):
    jogadores = listar(nome_campanha)
    achou = False
    for i, j in enumerate(jogadores):
        if j.get("nome") == nome:
            jogadores[i] = mesclar_ficha(j, patch)
            achou = True
            break
    if not achou:
        raise ValueError(f"Jogador '{nome}' não encontrado.")
    _gravar(nome_campanha, jogadores)
    return jogadores


def remover(nome_campanha, nome):
    jogadores = [j for j in listar(nome_campanha) if j.get("nome") != nome]
    _gravar(nome_campanha, jogadores)
    return jogadores


def ajustar_vida(nome_campanha, nome, delta):
    jogadores = listar(nome_campanha)
    for j in jogadores:
        if j.get("nome") == nome:
            j["vida"] = max(0, min(99, int(j.get("vida", 20)) + int(delta)))
            _gravar(nome_campanha, jogadores)
            return j["vida"]
    raise ValueError(f"Jogador '{nome}' não encontrado.")
