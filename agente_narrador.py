# agente_narrador.py — sub-agente escrivão (Why-Man)
from colorama import Fore, Style
import gerenciador_json
import gerador_enc
import gerador_npc
from why_man.escrivao import narrar_como_mestre
import gerenciador_capitulos


def _colorir(texto: str) -> str:
    return (
        texto.replace("[VERDE]", Fore.GREEN)
        .replace("[VERMELHO]", Fore.RED)
        .replace("[AMARELO]", Fore.YELLOW)
        .replace("[MAGENTA]", Fore.MAGENTA)
        .replace("[/]", Style.RESET_ALL)
    )


def _vida_do_desfecho(nome_campanha, desfecho: str):
    jogadores = gerenciador_json.carregar(nome_campanha, ".jogadores.json") or []
    for j in jogadores:
        nome = j.get("nome") or ""
        if nome and nome in desfecho:
            return int(j.get("vida", 20))
    if jogadores:
        return int(jogadores[0].get("vida", 20))
    return None


def _registrar_cronica(nome_campanha, narracao: str, comentario: str | None):
    resumo = gerenciador_json.carregar(nome_campanha, "resumo.json") or {
        "historia": "",
        "dialogos_importantes": [],
    }
    frase = " ".join(narracao.replace("[VERDE]", "").replace("[VERMELHO]", "")
                     .replace("[AMARELO]", "").replace("[MAGENTA]", "")
                     .replace("[/]", "").split())
    historia = (resumo.get("historia") or "").strip()
    if frase:
        historia = (historia + " " + frase).strip()
        if len(historia) > 2500:
            historia = historia[-2500:]
    resumo["historia"] = historia
    if comentario:
        falas = resumo.get("dialogos_importantes") or []
        falas.append(comentario)
        resumo["dialogos_importantes"] = falas[-12:]
    gerenciador_json.salvar(nome_campanha, "resumo.json", resumo)


def gerar_narracao(nome_campanha, desfecho_tecnico):
    resumo = gerenciador_json.carregar(nome_campanha, "resumo.json") or {}
    contexto = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    historico = gerenciador_json.carregar(nome_campanha, ".hist.json") or []
    cenario = contexto.get("contexto_inicial") or resumo.get("historia") or contexto.get("descricao") or "uma estrada que ainda não tem nome"
    if len(cenario) > 280:
        cenario = cenario[:280]
    objetivo = contexto.get("objetivo_principal", "a jornada")
    vida = _vida_do_desfecho(nome_campanha, desfecho_tecnico)

    cena, comentario = narrar_como_mestre(
        desfecho_tecnico,
        cenario=cenario,
        objetivo=objetivo,
        vida=vida,
        historico=historico,
    )
    _registrar_cronica(nome_campanha, cena, comentario)
    fechado = gerenciador_capitulos.anotar(nome_campanha, desfecho_tecnico, cena)
    return cena, comentario, fechado


def narrar(nome_campanha, desfecho_tecnico, imprimir=True):
    cena, comentario, fechado = gerar_narracao(nome_campanha, desfecho_tecnico)
    if imprimir:
        print(f"\n{Fore.MAGENTA}Mestre (Why-Man, escrivão):{Style.RESET_ALL}\n")
        print(_colorir(cena) + Style.RESET_ALL)
        if comentario:
            print(f"\n{Fore.CYAN}{Style.DIM}{comentario}{Style.RESET_ALL}")
        if fechado:
            print(f"\n{Fore.YELLOW}Capítulo encerrado: {fechado.get('titulo')}{Style.RESET_ALL}")
            print(_colorir(fechado.get("resumo", "")))
    return cena, comentario, fechado


def disparar_encontro_procedural(nome_campanha):
    evento = gerador_enc.gerar_encontro(nome_campanha)
    narrar(nome_campanha, f"Um evento inesperado acontece: {evento}")


def disparar_npc_procedural(nome_campanha):
    npc = gerador_npc.gerar_npc(nome_campanha)
    desfecho = (
        f"O grupo encontra {npc['nome']}, um(a) {npc['raca']} {npc['ocupacao']}. "
        f"Motivação: {npc['motivacao']}."
    )
    narrar(nome_campanha, desfecho)
