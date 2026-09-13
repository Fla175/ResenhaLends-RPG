# main.py
import sys
from colorama import init, Fore, Style
from questionary import questionary
import gerenciador_json
import gerenciador_jogadores
import gerenciador_capitulos
import dados
import agente_estado
import agente_narrador
import codex
import ia

init(autoreset=True)

def exibir_logo():
    logo = fr"""{Fore.CYAN}
    =================================================
 ____                      _             _                  _
|  _ \ ___  ___  ___ _ __ | |__   ____  | |    ___ _ __   __| |____
| |_) / _ \/ __|/ _ \ '_ \| '_ \ | '_ \ | |   / _ \ '_ \ / _` / ___|
|  _ <  __/\__ \  __/ | | | | | || |_| || |__|  __/ | | | (_| \__  \
|_| \_\___||___/\___|_| |_|_| |_||_| |_||____\___|__| |_|\__,_|____/
                                                     
        RPG SELF-HOSTED - MESTRE DE MESA PARA RPG
    =================================================
    """
    print(logo)

def menu():
    exibir_logo()
    titulo = f"{Fore.YELLOW}Escolha o seu destino:{Style.RESET_ALL}"
    opcoes = [
        f"{Fore.GREEN}1.{Style.RESET_ALL} Nova Campanha",
        f"{Fore.GREEN}2.{Style.RESET_ALL} Carregar Campanha",
        f"{Fore.GREEN}3.{Style.RESET_ALL} Gerenciar Regras",
        f"{Fore.GREEN}4.{Style.RESET_ALL} Exportar Campanha",
        f"{Fore.GREEN}5.{Style.RESET_ALL} Excluir Campanha",
        f"{Fore.GREEN}6.{Style.RESET_ALL} Sair"]
    escolha = questionary.select(titulo, choices=opcoes, pointer=f"{Fore.CYAN}$ {Style.RESET_ALL}").ask()
    return opcoes.index(escolha) if escolha else 5

def escolher_campanha():
    campanhas = gerenciador_json.listar_campanhas()
    if not campanhas:
        print(f"{Fore.RED}Nenhuma campanha encontrada.{Style.RESET_ALL}")
        return None
    for i, c in enumerate(campanhas, 1):
        print(f"[{Fore.YELLOW}{i}{Style.RESET_ALL}] {c}")
    escolha = input(f"\n{Fore.BLUE}Número ou Nome da campanha: {Style.RESET_ALL}")
    if escolha.isdigit() and 1 <= int(escolha) <= len(campanhas):
        return campanhas[int(escolha) - 1]
    if escolha in campanhas:
        return escolha
    print(f"{Fore.RED}Inválido.{Style.RESET_ALL}")
    return None

def exportar():
    nome = escolher_campanha()
    if not nome:
        return
    destino = input("Pasta de destino (Enter = exportacoes/): ").strip() or None
    try:
        arq = gerenciador_json.exportar_campanha(nome, destino)
        print(f"{Fore.GREEN}Exportado: {arq}{Style.RESET_ALL}")
    except Exception as e:
        print(f"{Fore.RED}{e}{Style.RESET_ALL}")

def exibir_hud(nome_campanha, jogadores):
    contexto = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    print(f"\n{Fore.CYAN}{'='*60}")
    print(f" CAMPANHA: {nome_campanha.upper()} | OBJETIVO: {contexto.get('objetivo_principal', '?')}")
    print(f"{'='*60}{Style.RESET_ALL}")
    print(f"{Fore.GREEN} | {Style.RESET_ALL}".join(
        f"{j['nome']} ({j.get('classe', 'Avent')}) [HP: {j.get('vida', 20)}]" for j in jogadores
    ))

def criar_jogadores(nome_campanha):
    racas = ["Humano", "Elfo", "Anão", "Halfling", "Meio-Orc", "Tiefling"]
    while True:
        nome = input("\nNome do jogador (ou Enter para finalizar): ").strip()
        if not nome:
            break
        classe = input(f"Classe de {nome}: ").strip()
        desc = input(f"Descrição de {nome}: ").strip()
        raca = questionary.select(f"Raça de {nome}:", choices=racas).ask() or "Humano"
        try:
            gerenciador_jogadores.adicionar(
                nome_campanha,
                {"nome": nome, "classe": classe, "descricao": desc, "raca": raca},
                rolar_atts=True,
            )
            print(f"{Fore.GREEN}{nome} adicionado.{Style.RESET_ALL}")
        except Exception as e:
            print(f"{Fore.RED}{e}{Style.RESET_ALL}")

def criar_campanha():
    nome = input("Nome da campanha: ").strip().replace(" ", "_").lower()
    if not nome:
        return
    descricao = input("Descrição da campanha: ").strip()
    objetivo = input("Objetivo principal: ").strip()
    dificuldade = questionary.select("Dificuldade:", choices=["Fácil", "Normal", "hardcore"]).ask() or "Normal"
    tam_map = {"Pequena": "pequena", "Média": "media", "Extensa": "extensa"}
    tam = tam_map.get(questionary.select("Tamanho:", choices=list(tam_map)).ask(), "media")
    gerenciador_json.setup_campanha(nome, dificuldade, descricao, objetivo, tam)
    gerenciador_capitulos.garantir_abertura(nome)
    print(Fore.MAGENTA + "Mestre:" + Style.RESET_ALL)
    ctx = gerenciador_json.carregar(nome, ".context.json") or {}
    print((ctx.get("contexto_inicial") or "").replace("[MAGENTA]", "").replace("[AMARELO]", "").replace("[VERDE]", "").replace("[/]", ""))
    criar_jogadores(nome)
    loop_jogo(nome)

def excluir_campanha():
    nome = escolher_campanha()
    if nome and input(f"{Fore.RED}Excluir '{nome}'? (s/n): {Style.RESET_ALL}").lower() == "s":
        gerenciador_json.excluir_campanha(nome)

def gerenciar_regras():
    nome = escolher_campanha()
    if not nome:
        return
    regras = gerenciador_json.carregar(nome, ".regras.json") or {}
    print(f"Regras atuais: {regras}")
    chave = input("Qual regra alterar/criar? ")
    if chave:
        regras[chave] = input(f"Novo valor para '{chave}': ")
        gerenciador_json.salvar(nome, ".regras.json", regras)

def loop_jogo(nome_campanha):
    gerenciador_capitulos.garantir_abertura(nome_campanha)
    ctx = gerenciador_json.carregar(nome_campanha, ".context.json") or {}
    if ctx.get("contexto_inicial"):
        print(f"\n{Fore.MAGENTA}Mestre — contexto inicial:{Style.RESET_ALL}\n{ctx['contexto_inicial']}")
    jogadores = gerenciador_jogadores.listar(nome_campanha)
    if not jogadores:
        criar_jogadores(nome_campanha)
        jogadores = gerenciador_jogadores.listar(nome_campanha)
    turno_idx = 0
    while True:
        exibir_hud(nome_campanha, jogadores)
        jogador_atual = jogadores[turno_idx]
        nome_j = jogador_atual["nome"]
        macros = jogador_atual.get("macros", {})
        acao = input(f"\n{Fore.YELLOW}[Turno de {nome_j}] ação / capitulos / encerrar / exportar / sair: {Style.RESET_ALL}").strip()
        if acao.lower() in ("sair", "menu"):
            break
        if acao.lower() in ("capitulos", "capítulos"):
            for cap in gerenciador_capitulos.listar(nome_campanha):
                print(f"\n{Fore.YELLOW}{cap.get('titulo')}{Style.RESET_ALL} ({cap.get('fechado_em', '')})")
                print(cap.get("resumo", ""))
            atual = gerenciador_capitulos.atual(nome_campanha)
            print(f"\nEm andamento: capítulo {atual.get('numero')} — {len(atual.get('eventos') or [])} cenas")
            continue
        if acao.lower() == "encerrar":
            fechado = gerenciador_capitulos.encerrar(nome_campanha)
            if fechado:
                print(f"{Fore.YELLOW}{fechado.get('titulo')}{Style.RESET_ALL}\n{fechado.get('resumo')}")
            else:
                print("Nada para fechar ainda.")
            continue
        if acao.lower() in ("sair", "menu"):
            break
        if acao.lower() == "exportar":
            print(gerenciador_json.exportar_campanha(nome_campanha))
            continue
        if acao.lower() == "!encontro":
            agente_narrador.disparar_encontro_procedural(nome_campanha)
            continue
        if acao.lower() == "!npc":
            agente_narrador.disparar_npc_procedural(nome_campanha)
            continue
        if acao.lower() == "!bestiario":
            codex.menu_bestiario(nome_campanha)
            continue
        if acao in macros and macros[acao]:
            acao = macros[acao]
        dado = dados.rolar()
        desfecho = agente_estado.processar_acao(nome_campanha, f"O jogador {nome_j} tenta: {acao}", dado)
        agente_narrador.narrar(nome_campanha, desfecho)
        jogadores = gerenciador_jogadores.listar(nome_campanha)
        turno_idx = (turno_idx + 1) % len(jogadores)

def main_terminal():
    ok, status = ia.verificar_conexao()
    if ok:
        print(f"{Fore.GREEN}{status}{Style.RESET_ALL}")
    print(f"{Fore.YELLOW}Regras locais. Why-Man. Mesa gráfica: python main.py{Style.RESET_ALL}")
    while True:
        op = menu()
        if op == 0:
            criar_campanha()
        elif op == 1:
            c = escolher_campanha()
            if c:
                loop_jogo(c)
        elif op == 2:
            gerenciar_regras()
        elif op == 3:
            exportar()
        elif op == 4:
            excluir_campanha()
        elif op == 5:
            print(f"{Fore.GREEN}Até a próxima!{Style.RESET_ALL}")
            break

def main():
    if "--terminal" in sys.argv:
        main_terminal()
        return
    try:
        from ui_pygame import executar
        executar()
    except Exception as e:
        print(f"{Fore.RED}Pygame indisponível ({e}). pip install pygame-ce{Style.RESET_ALL}")
        main_terminal()

if __name__ == "__main__":
    main()
