"""Why-Man: prólogo do mestre e fechamento de capítulos."""
from __future__ import annotations

import random

from why_man.narrador_local import narrar_local


def _limpo(texto: str) -> str:
    for t in ("[VERDE]", "[VERMELHO]", "[AMARELO]", "[MAGENTA]", "[/]"):
        texto = texto.replace(t, "")
    return " ".join(texto.split())


def gerar_contexto_inicial(descricao: str, objetivo: str, nome_campanha: str) -> str:
    """Gera contexto inicial usando IA ou fallback local."""
    mundo = descricao.strip() or "um mundo que ainda não confessou o nome"
    meta = objetivo.strip() or "sobreviver à própria história"
    
    # Tenta gerar contexto mais rico usando o narrador local
    try:
        tom = narrar_local(f"{mundo} {meta}")
        # Se o narrador local retornar algo além do input básico, usa IA
        if tom and len(tom) > len(f"{mundo} {meta}"):
            aberturas = [
                f"[MAGENTA]Antes de qualquer dado, escutem.[/] {mundo} "
                f"[AMARELO]O que vocês acham que vieram fazer: {meta}.[/] "
                f"{tom} Eu, o mestre, já vi o mapa. Vocês só têm as botas e a teimosia.",
                f"[MAGENTA]Sessão zero, mas com o mundo já respirando.[/] {mundo} "
                f"Se falharem, {meta} vira lenda de taverna. Se forem competentes, vira ferida. "
                f"[VERDE]Hoje a mesa abre neste tom.[/] {tom}",
                f"Eu não vou ler um livro. Vou descrever o chão. {mundo} "
                f"[AMARELO]No horizonte, sempre: {meta}.[/] "
                f"Quando o dado cair, lembrem que o universo desta campanha ({nome_campanha.replace('_', ' ')}) tem opiniões.",
            ]
            return random.choice(aberturas)
    except Exception:
        pass
    
    # Fallback para geração determinística
    aberturas = [
        f"[MAGENTA]Antes de qualquer dado, escutem.[/] {mundo} "
        f"[AMARELO]O que vocês acham que vieram fazer: {meta}.[/] "
        f"Eu, o mestre, já vi o mapa. Vocês só têm as botas e a teimosia.",
        f"[MAGENTA]Sessão zero, mas com o mundo já respirando.[/] {mundo} "
        f"Se falharem, {meta} vira lenda de taverna. Se forem competentes, vira ferida. "
        f"[VERDE]Hoje a mesa abre neste tom.[/]",
        f"Eu não vou ler um livro. Vou descrever o chão. {mundo} "
        f"[AMARELO]No horizonte, sempre: {meta}.[/] "
        f"Quando o dado cair, lembrem que o universo desta campanha ({nome_campanha.replace('_', ' ')}) tem opiniões.",
    ]
    return random.choice(aberturas)


def fechar_capitulo(numero: int, eventos: list, objetivo: str) -> dict:
    bloco = " ".join(str(e) for e in eventos).lower()
    tom = narrar_local(bloco[:400] or objetivo)

    if "falha crítica" in bloco or bloco.count("falha") >= 2:
        titulo = random.choice([
            f"Capítulo {numero}: O preço da pressa",
            f"Capítulo {numero}: Quando o dado recusou o pedido",
        ])
        clima = "sangue, silêncio e um mestre fazendo cara de paisagem"
    elif "sucesso crítico" in bloco:
        titulo = random.choice([
            f"Capítulo {numero}: Uma página que a mesa não esquece",
            f"Capítulo {numero}: O dia em que o vinte apareceu",
        ])
        clima = "glória curta, ego alto, mundo ainda perigoso"
    elif "encontro" in bloco or "combate" in bloco or "emboscada" in bloco:
        titulo = random.choice([
            f"Capítulo {numero}: A estrada morde de volta",
            f"Capítulo {numero}: Iniciativa sem aviso",
        ])
        clima = "aço, poeira e alguém perguntando quem ainda tem poção"
    elif "encontra" in bloco:
        titulo = random.choice([
            f"Capítulo {numero}: Rostos no caminho",
            f"Capítulo {numero}: Alguém que não estava no roteiro de vocês",
        ])
        clima = "conversa, suspeita e um NPC com sorriso demais"
    else:
        titulo = random.choice([
            f"Capítulo {numero}: Mais um palmo rumo a {objetivo[:40]}",
            f"Capítulo {numero}: A jornada continua, teimosa",
        ])
        clima = "passos, portas e o objetivo ainda no horizonte"

    fatos = eventos[-4:] if eventos else ["A etapa passou em silêncio."]
    resumo = (
        f"{titulo}. {tom} Nesta etapa: {clima}. "
        f"O escrivão guarda: {' | '.join(_limpo(str(f))[:90] for f in fatos)}."
    )
    return {
        "numero": numero,
        "titulo": titulo,
        "resumo": _limpo(resumo),
        "eventos": [_limpo(str(e)) for e in eventos],
        "tipo": "etapa",
    }
