# gerador_enc.py
import random
from colorama import Fore, Style
import gerenciador_json

def gerar_encontro(nome_campanha):
    regras = gerenciador_json.carregar(nome_campanha, ".regras.json") or {}
    dificuldade = regras.get("dificuldade", "normal").lower()
    
    tipos = ["Combate", "Emboscada", "Armadilha", "Descoberta", "Clima Extremo", "Encontro Social", "Enigma", "NPC Perdido", "Tesouro Oculto", "Obstáculo Natural", "Sinal Misterioso"]
    
    ameacas = {
        "fácil": [
            "Bando de Goblins",
            "Lobos Selvagens",
            "Bandidos Novatos",
            "Ratos Gigantes",
            "Slime Ácido",
            "Cultistas Aprendizes",
            "Kobolds Ladinos",
        ],
        "normal": [
            "Orcs Furiosos",
            "Esqueletos Armados",
            "Urso-Coruja",
            "Aranhas Gigantes",
            "Gárgula de Pedra",
            "Elemental da Terra",
            "Manticora",
            "Feiticeiro Renegado",
        ],
        "hardcore": [
            "Troll de Pedra",
            "Vampiro Ancião",
            "Dragão Jovem",
            "Lich Enfraquecido",
            "Devorador de Mentes",
            "Demônio da Fornalha",
            "Beholder (Contemplador)",
            "Hydra das Cavernas",
        ],
    }
    
    locais = [
        # Florestas e Natureza
        "em uma clareira densa",
        "atrás de folhagens espessas",
        "sob as raízes de uma árvore gigante",
        "em um pântano lamacento e fétido",
        "no coração de um bosque sussurrante",
        
        # Estradas e Passagens
        "no meio da estrada",
        "em uma ponte de pedra antiga",
        "em uma bifurcação deserta da trilha",
        "no cruzamento de estradas abandonadas",
        "em um túnel estreito escavado na rocha",
        
        # Ruínas e Estruturas
        "perto de ruínas desmoronadas",
        "no interior de uma cabana em ruínas",
        "nos escombros de uma torre de vigia",
        "no pátio de um castelo esquecido",
        "na entrada de um templo sagrado profanado",
        
        # Relevo e Cavernas
        "à beira de um penhasco íngreme",
        "dentro de uma caverna úmida",
        "no topo de uma colina enevoada",
        "no fundo de um desfiladeiro rochoso",
        "na entrada de uma mina abandonada",
        
        # Água e Litoral
        "próximo a um rio turbulento",
        "à margem de um lago de águas escuras",
        "ao lado de uma cachoeira barulhenta",
        "em uma praia de areia negra",
        
        # Sombrios e Misteriosos
        "em um cemitério abandonado",
        "no centro de um círculo de pedras rúnicas",
        "onde a névoa impede de enxergar um palmo à frente"
    ]

    tipo = random.choice(tipos)
    ameaca = random.choice(ameacas.get(dificuldade, ameacas["normal"]))
    local = random.choice(locais)
    
    if tipo == "Descoberta":
        evento = f"O grupo encontra um baú trancado ou santuário antigo {local}."
    elif tipo == "Armadilha":
        evento = f"Alguém pisa em uma armadilha escondida {local}."
    elif tipo == "Clima Extremo":
        evento = f"Uma tempestade violenta ou nevasca repentina atinge o grupo {local}."
    else:
        evento = f"{tipo}: Vocês são surpreendidos por um {ameaca} {local}."
        
    print(f"\n{Fore.RED}⚠️ ENCONTRO GERADO: {evento}{Style.RESET_ALL}")
    
    # Injeta o encontro no histórico para o Ollama ler no próximo turno
    historico = gerenciador_json.carregar(nome_campanha, ".hist.json") or []
    historico.append(f"Evento Aleatório: {evento}")
    gerenciador_json.salvar(nome_campanha, ".hist.json", historico)
    
    return evento