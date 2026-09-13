# gerenciador_json.py
import json
import os
import shutil
import stat
from datetime import datetime

PASTA_BASE = "campanhas"
PASTA_EXPORT = "exportacoes"

def obter_caminho(nome_campanha, arquivo=""):
    if arquivo:
        return os.path.join(PASTA_BASE, nome_campanha, arquivo)
    return os.path.join(PASTA_BASE, nome_campanha)

def listar_campanhas():
    if not os.path.exists(PASTA_BASE):
        return []
    return [n for n in os.listdir(PASTA_BASE) if os.path.isdir(obter_caminho(n))]

def setup_campanha(nome_campanha, dificuldade, descricao="", objetivo_principal="", tamanho_campanha="medio"):
    caminho = obter_caminho(nome_campanha)
    os.makedirs(caminho, exist_ok=True)

    arquivos_base = {
        ".context.json": {
            "nome_campanha": nome_campanha,
            "descricao": descricao,
            "objetivo_principal": objetivo_principal,
            "tamanho_campanha": tamanho_campanha
        },
        ".regras.json": {
            "dificuldade": dificuldade,
            "magia_permitida": True
        },
        ".jogadores.json": [],
        ".hist.json": [],
        ".npc.json": {},
        ".monstros.json": {},
        ".notas_mestre.json": {
            "anotacoes": "",
            "segredos": ""
        },
        "resumo.json": {
            "historia": "",
            "dialogos_importantes": []
        },
        ".capitulos.json": [],
        ".capitulo_atual.json": {
            "numero": 1,
            "eventos": [],
            "aberto": True
        }
    }

    for arquivo, dados_iniciais in arquivos_base.items():
        arq_path = obter_caminho(nome_campanha, arquivo)
        if not os.path.exists(arq_path):
            salvar(nome_campanha, arquivo, dados_iniciais)
    migrar_bestiario(nome_campanha)

def salvar(nome_campanha, arquivo, dados):
    caminho = obter_caminho(nome_campanha, arquivo)
    os.makedirs(os.path.dirname(caminho), exist_ok=True)
    with open(caminho, "w", encoding="utf-8") as f:
        json.dump(dados, f, indent=4, ensure_ascii=False)

def carregar(nome_campanha, arquivo):
    caminho = obter_caminho(nome_campanha, arquivo)
    try:
        with open(caminho, "r", encoding="utf-8") as f:
            return json.load(f)
    except FileNotFoundError:
        return None

def tratar_erro_remocao(func, path, exc_info):
    os.chmod(path, stat.S_IWRITE)
    func(path)

def excluir_campanha(nome_campanha):
    caminho = obter_caminho(nome_campanha)
    if os.path.exists(caminho):
        shutil.rmtree(caminho, onerror=tratar_erro_remocao)
        return True
    return False

def migrar_bestiario(nome_campanha):
    antigo = carregar(nome_campanha, ".mostros.json")
    atual = carregar(nome_campanha, ".monstros.json")
    if antigo and not atual:
        salvar(nome_campanha, ".monstros.json", antigo)
    if atual is None:
        salvar(nome_campanha, ".monstros.json", antigo or {})

def exportar_campanha(nome_campanha, destino_dir=None):
    origem = obter_caminho(nome_campanha)
    if not os.path.isdir(origem):
        raise FileNotFoundError(f"Campanha '{nome_campanha}' não existe.")
    destino_dir = destino_dir or PASTA_EXPORT
    os.makedirs(destino_dir, exist_ok=True)
    stamp = datetime.now().strftime("%Y%m%d_%H%M")
    base = os.path.join(destino_dir, f"{nome_campanha}_{stamp}")
    arquivo = shutil.make_archive(base, "zip", root_dir=PASTA_BASE, base_dir=nome_campanha)
    return arquivo

def importar_campanha(caminho_zip):
    if not os.path.isfile(caminho_zip):
        raise FileNotFoundError("Arquivo zip não encontrado.")
    os.makedirs(PASTA_BASE, exist_ok=True)
    shutil.unpack_archive(caminho_zip, PASTA_BASE)
    nomes = listar_campanhas()
    if nomes:
        migrar_bestiario(nomes[-1])
    return nomes
