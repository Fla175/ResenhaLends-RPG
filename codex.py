# codex.py
from colorama import Fore, Style
import gerenciador_json
from pick import pick

def carregar_bestiario(nome_campanha):
    gerenciador_json.migrar_bestiario(nome_campanha)
    return gerenciador_json.carregar(nome_campanha, ".monstros.json") or {}

def salvar_monstro(nome_campanha, id_monstro, dados_monstro):
    bestiario = carregar_bestiario(nome_campanha)
    bestiario[id_monstro] = dados_monstro
    gerenciador_json.salvar(nome_campanha, ".monstros.json", bestiario)

def buscar_monstro(nome_campanha, termo):
    bestiario = carregar_bestiario(nome_campanha)
    termo = termo.lower()
    return {k: v for k, v in bestiario.items() if termo in v.get("nome", "").lower() or termo in k.lower()}

def formatar_ficha(monstro):
    linha = f"{Fore.CYAN}{'-'*40}{Style.RESET_ALL}"
    atts = monstro.get("atributos", {})
    
    ficha = [
        linha,
        f"{Fore.YELLOW}   === {monstro.get('nome', 'Monstro').upper()} ==={Style.RESET_ALL}",
        f"Tipo: {monstro.get('tipo', 'Desconhecido')} | ND: {monstro.get('nd', 'N/A')}",
        linha,
        f"PV: {monstro.get('pv', 'N/A')} | CA: {monstro.get('ca', 'N/A')}",
        f"FOR: {atts.get('for', 10)} | DES: {atts.get('des', 10)} | CON: {atts.get('con', 10)}",
        f"INT: {atts.get('int', 10)} | SAB: {atts.get('sab', 10)} | CAR: {atts.get('car', 10)}",
        linha
    ]
    return "\n".join(ficha)

def menu_bestiario(nome_campanha):
    bestiario = carregar_bestiario(nome_campanha)
    print(f"\n{Fore.MAGENTA}=== BESTIÁRIO DA CAMPANHA ==={Style.RESET_ALL}")
    if not bestiario:
        print(f"{Fore.YELLOW}Nenhum monstro registrado no Bestiário ainda.{Style.RESET_ALL}")
        return

    for i, (m_id, m_dados) in enumerate(bestiario.items(), 1):
        print(f"[{i}] {m_dados.get('nome', m_id)} (ND: {m_dados.get('nd', 'N/A')})")

    chaves = list(bestiario.keys())
    opcoes = chaves + ["voltar"]
    titulo = "Escolha um monstro para ver a ficha:"

    opt_selecionada, index = pick(
        options=opcoes,
        title=titulo,
        indicator="$ ",
        default_index=0,
    )

    if opt_selecionada != "Voltar":
        print(formatar_ficha(bestiario[opt_selecionada]))
