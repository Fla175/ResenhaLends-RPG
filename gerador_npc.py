# gerador_npc.py
import random
from colorama import Fore, Style
import gerenciador_json

def gerar_npc(nome_campanha):
    nomes = ["Thalor", "Elara", "Bram", "Kael", "Lyra", "Gorn", "Silas", "Maeve", "Doran", "Vex", "Eldrin", "Mira", "Tavish", "Kira", "Borin", "Zalriel", "Faelar", "Yvaine", "Garrick", "Sariel", "Fenric", "Liora", "Dax", "Valen", "Rowan", "Talia", "Kaelen", "Morwen", "Orin", "Reyna", "Zephyr", "Aveline", "Caelum", "Elowen", "Gideon", "Hadrian", "Isolde", "Jarek", "Malakor", "Theron", "Ignis", "Storm", "Thurn", "Zorvath", "Kallista", "Brottor", "Fira", "Keth", "Milo", "Zanna", "Xalthor", "Ruk", "Garnet", "Marrow", "Echo", "Rusty", "Torinn", "Nisha", "Aurelius", "Bryn"]

    sexo = ["Masculino", "Feminino"]

    racas = ["Humano", "Elfo", "Anão", "Halfling", "Meio-Orc", "Draconato", "Tiefling", "Gnomo", "Goliath", "Aasimar", "Firbolg", "Tabaxi", "Kenqu", "Lizardfolk", "Changeling", "Warforged", "Centauro", "Minotauro", "Goblinoide", "Orc", "Drow", "Duérgar", "Githyanki", "Githzerai", "Satir", "Tritão", "Genasi do Fogo", "Genasi da Água", "Genasi da Terra", "Genasi do Ar", "Shifter", "Leonino", "Loxodonte", "Vedalken", "Simic Híbrido", "Kalashtar", "Aerie", "Gnomo das Profundezas", "Elfo da Floresta", "Anão da Montanha"]

    ocupacoes = ["Ferreiro", "Mercador Suspeito", "Guarda da Cidade", "Mago Andarilho", "Ladino", "Alquimista Excêntrico", "Coveiro Ranzinza", "Taberneiro Falador", "Caçador de Recompensas", "Sacerdote Corrupto", "Bardo Fracassado", "Ervanária da Vila", "Mensageiro Real", "Navegador Aposentado", "Estudante de Arqueologia", "Guarda-costas Silencioso", "Cozinheiro de Taverna", "Contrabandista de Carga", "Alfaiate da Alta Sociedade", "Eremita da Floresta"]

    motivos = ["Busca vingança", "Quer ficar rico rápido", "Foge do passado", "Procura um artefato perdido", "Protege um segredo antigo", "Busca o perdão por um crime", "Quer provar seu valor a uma guilda", "Deseja curar uma doença incurável", "Pesquisa uma lenda esquecida", "Busca fama e glória no combate", "Quer reencontrar um amor perdido", "Tenta saldar uma dívida de jogo", "Protege sua família ameaçada", "Quer dominar uma magia proibida", "Foge de um casamento arranjado", "Deseja recuperar a honra da família", "Busca entender uma visão profética", "Quer derrubar um tirano local", "Explora o mundo por puro tédio", "Coleta espécimes raros para estudo"]

    segredos = ["É um espião disfarçado", "Tem uma maldição oculta", "Sabe onde está o chefe final", "Deve dinheiro à guilda das sombras", "Na verdade, é um lobisomem", "Roubou a identidade de um nobre falecido", "Serve secretamente a uma divindade maligna", "Possui um mapa para o covil do dragão", "Foi o responsável por incendiar a vila vizinha", "É o herdeiro legítimo do trono que foi deposto", "Sua riqueza vem de pactos com demônios", "Está envenenando lentamente o líder local", "Vendeu informações sigilosas para os inimigos", "É uma criatura mágica disfarçada de humano", "Está fingindo ter amnésia para esconder um crime", "Possui a última cópia de um grimório proibido", "Foi enviado para sabotar a missão do grupo", "É o verdadeiro assassino do antigo rei", "Sua linhagem sanguínea ativa portais mágicos", "Guarda a chave que liberta o grande vilão"]

    npc = {
        "nome": random.choice(nomes),
        "sexo": random.choice(sexo),
        "raca": random.choice(racas),
        "ocupacao": random.choice(ocupacoes),
        "motivacao": random.choice(motivos),
        "segredo": random.choice(segredos)
    }

    # Salva diretamente no banco de NPCs da campanha
    npcs_salvos = gerenciador_json.carregar(nome_campanha, ".npc.json") or {}
    npcs_salvos[npc["nome"]] = npc
    gerenciador_json.salvar(nome_campanha, ".npc.json", npcs_salvos)

    print(f"\n{Fore.MAGENTA}--- NOVO NPC CRIADO ---{Style.RESET_ALL}")
    for chave, valor in npc.items():
        print(f"{Fore.CYAN}{chave.capitalize()}:{Style.RESET_ALL} {valor}")
    print(f"{Fore.MAGENTA}-----------------------{Style.RESET_ALL}")
    
    return npc