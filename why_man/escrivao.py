"""Why-Man escrivão: narra como mestre de mesa e anota a crônica."""
from __future__ import annotations

import random
import re

from why_man.narrador_local import narrar_local


def _aspas(texto: str) -> str:
    m = re.search(r"'([^']+)'", texto)
    return m.group(1).strip() if m else ""


def _nome(desfecho: str) -> str:
    m = re.search(
        r":\s+(.+?)\s+(?:tenta|realiza|consegue|não consegue)",
        desfecho,
        re.I,
    )
    if m:
        return m.group(1).strip()
    return "O grupo"


def _tipo(desfecho: str) -> str:
    t = desfecho.lower()
    if "falha crítica" in t or "falha critica" in t:
        return "crit_falha"
    if "sucesso crítico" in t or "sucesso critico" in t:
        return "crit_sucesso"
    if "sucesso parcial" in t:
        return "parcial"
    if desfecho.lower().startswith("sucesso"):
        return "sucesso"
    if "encontra " in t or "npc" in t:
        return "npc"
    if "evento inesperado" in t or "combate" in t or "emboscada" in t:
        return "encontro"
    if t.startswith("falha"):
        return "falha"
    return "neutro"


def _cena(tipo: str, nome: str, acao: str, mundo: str, objetivo: str, tom: str) -> str:
    acao = acao or "seguir em frente"
    mundo = mundo or "o caminho à frente"
    cenas = {
        "crit_sucesso": [
            f"[MAGENTA]{nome} mal termina de declarar '{acao}' e o dado explode em vinte.[/] "
            f"[VERDE]O mundo inteiro parece conspirar a favor: {mundo} se abre como se esperasse exatamente isso.[/] "
            f"Por um instante, {objetivo} parece perto o bastante para tocar.",
            f"[VERDE]Não foi sorte. Foi teatro.[/] {nome} executa '{acao}' com uma precisão ofensiva. "
            f"[MAGENTA]Até o cenário fica quieto, respeitoso.[/]",
        ],
        "crit_falha": [
            f"[VERMELHO]Um. Natural. {nome} tenta '{acao}' e o universo responde com um tapinha na alma.[/] "
            f"Em {mundo}, nada disso deveria ter dado tão errado — e mesmo assim deu. "
            f"[VERMELHO]Alguém vai lembrar disso na taverna. Principalmente você.[/]",
            f"[VERMELHO]{nome} escorrega no timing, na gravidade e na dignidade ao mesmo tempo.[/] "
            f"'{acao}' vira um acidente público. O objetivo ({objetivo}) assiste de camarote.",
        ],
        "sucesso": [
            f"[VERDE]{nome} vai de '{acao}' e funciona.[/] Não é épico, é competente — o tipo de competência que mantém grupo vivo. "
            f"Em {mundo}, as peças andam um palmo. [AMARELO]{objetivo} continua no horizonte, mas agora com um caminho a mais.[/]",
            f"Você descreve '{acao}'. Eu aceito. [VERDE]O teste passa.[/] "
            f"{nome} sente o chão responder: portas rangem para o lado certo, olhares desviam, o ar fica utilizável.",
        ],
        "parcial": [
            f"[AMARELO]{nome} quase leva '{acao}' no colo... quase.[/] O resultado vem pela metade, com juros. "
            f"[VERMELHO]Dói um pouco. Funciona um pouco.[/] É o clássico 'técnicamente sim'.",
            f"Olha, eu deixo. [AMARELO]Mas não de graça.[/] '{acao}' abre uma fresta em {mundo}. "
            f"{nome} passa — mancando, teimoso, vivo.",
        ],
        "falha": [
            f"[VERMELHO]{nome} tenta '{acao}' e o dado não colabora.[/] "
            f"Em {mundo}, a recusa é física: um golpe, um fechamento, um 'não' com peso. "
            f"O grupo ainda respira. O orgulho, nem tanto.",
            f"Eu até queria. [VERMELHO]O dado não queria.[/] '{acao}' falha na cara de {nome}. "
            f"{objetivo} fica um pouco mais longe, como se tivesse dado de ombros.",
        ],
        "encontro": [
            f"[MAGENTA]Espera. Eu não pedi isso, o mundo que pediu.[/] {desfecho_curto(acao, tom)} "
            f"A mesa muda de tom: mãos na ficha, olho no mapa. [VERMELHO]Isso aqui não é flavor text.[/]",
        ],
        "npc": [
            f"[AMARELO]Alguém entra em cena que não estava no roteiro de vocês — estava no meu.[/] {tom} "
            f"Conversem. Ou peguem no cabo da arma. Eu narrarei as duas com o mesmo carinho.",
        ],
        "neutro": [
            f"A cena em {mundo} se acomoda em volta de {nome}. {tom} "
            f"Nada explode. Ainda. O escrivão anota assim mesmo.",
        ],
    }
    opcoes = cenas.get(tipo, cenas["neutro"])
    return random.choice(opcoes)


def desfecho_curto(acao: str, tom: str) -> str:
    if tom:
        return tom
    return acao


def comentario_de_mesa(tipo: str, nome: str, vida: int | None, historico: list, acao: str) -> str | None:
    """Falas de mestre fora da narração — só em ocasiões certas."""
    recentes = " ".join(historico[-4:]).lower() if historico else ""
    falhas_seguidas = recentes.count("falha") >= 2
    acao_l = (acao or "").lower()

    if tipo == "crit_falha":
        return random.choice([
            f"Mesa: {nome}... eu juro que não empurrei o dado. O universo que tem inveja.",
            "Mesa: Anota no caderno. Isso vira história. Infelizmente.",
            "Mesa: Podem rir. Vocês vão precisar.",
        ])
    if tipo == "crit_sucesso":
        return random.choice([
            "Mesa: Pode comemorar. Até eu fiquei impressionado, e eu que narro isso.",
            f"Mesa: {nome} acaba de comprar um capítulo inteiro na crônica. Não desperdicem.",
            "Mesa: Vinte. Eu fecho a boca e só descrevo o estrago.",
        ])
    if vida is not None and vida <= 4:
        return random.choice([
            f"Mesa: {nome} está com {vida} de vida. Alguém tem poção ou a gente mede o caixão?",
            "Mesa: Estou narrando em voz baixa de propósito. Não morre agora, fica feio na ata.",
        ])
    if falhas_seguidas and tipo == "falha":
        return random.choice([
            "Mesa: Terceiro tombo na sequência. Eu começo a achar que vocês gostam de apanhar.",
            "Mesa: Olha, a campanha continua. O orgulho coletivo, não sei.",
        ])
    if any(p in acao_l for p in ("pai", "seduz", "explode", "implod", "deus", "namor")):
        return random.choice([
            "Mesa: Vocês se metem em cada uma. Eu só anoto. Com gosto.",
            "Mesa: Esse é o tipo de ação que o escrivão sublinha duas vezes.",
        ])
    if tipo == "parcial" and random.random() < 0.7:
        return "Mesa: Sucesso com asterisco. Eu amo um asterisco."
    if tipo in ("sucesso", "npc", "encontro") and random.random() < 0.35:
        return random.choice([
            "Mesa: E então? O que o grupo faz com isso? Não me olhem assim, a cena é de vocês.",
            "Mesa: Vou beber água. Vocês pensem. Quinze segundos de tensão dramática grátis.",
            f"Mesa: {nome} puxou o foco. Resto da mesa, não durmam no dado.",
        ])
    return None


def narrar_como_mestre(
    desfecho: str,
    *,
    cenario: str = "",
    objetivo: str = "",
    vida: int | None = None,
    historico: list | None = None,
) -> tuple[str, str | None]:
    tipo = _tipo(desfecho)
    nome = _nome(desfecho)
    acao = _aspas(desfecho) or desfecho
    tom = narrar_local(desfecho)
    cena = _cena(tipo, nome, acao, cenario, objetivo, tom)
    comentario = comentario_de_mesa(tipo, nome, vida, historico or [], acao)
    return cena, comentario
