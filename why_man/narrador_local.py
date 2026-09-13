"""Why-Man local: classificador linear treinado no corpus (poucos mil parâmetros)."""
from __future__ import annotations

import json
import math
import os
import random
import re

from why_man.corpus import EXEMPLOS

PASTA = os.path.dirname(os.path.abspath(__file__))
PESOS = os.path.join(PASTA, "why_man_pesos.json")
DIM = 96
LR = 0.08
EPOCAS = 40


def _tokens(texto: str) -> list[str]:
    return re.findall(r"[a-zà-ÿ0-9]+", texto.lower())


def _hash(token: str) -> int:
    h = 2166136261
    for ch in token:
        h ^= ord(ch)
        h = (h * 16777619) & 0xFFFFFFFF
    return h % DIM


def vetor(texto: str) -> list[float]:
    v = [0.0] * DIM
    toks = _tokens(texto)
    if not toks:
        return v
    for t in toks:
        v[_hash(t)] += 1.0
    n = math.sqrt(sum(x * x for x in v)) or 1.0
    return [x / n for x in v]


def _dot(a: list[float], b: list[float]) -> float:
    return sum(x * y for x, y in zip(a, b))


def treinar_e_salvar(caminho: str = PESOS) -> dict:
    n = len(EXEMPLOS)
    pesos = [[random.uniform(-0.05, 0.05) for _ in range(DIM)] for _ in range(n)]
    vies = [0.0] * n
    xs = [vetor(fato) for fato, _ in EXEMPLOS]
    for _ in range(EPOCAS):
        ordem = list(range(n))
        random.shuffle(ordem)
        for i in ordem:
            logits = [_dot(pesos[k], xs[i]) + vies[k] for k in range(n)]
            m = max(logits)
            exps = [math.exp(z - m) for z in logits]
            s = sum(exps) or 1.0
            probs = [e / s for e in exps]
            for k in range(n):
                g = probs[k] - (1.0 if k == i else 0.0)
                for d in range(DIM):
                    pesos[k][d] -= LR * g * xs[i][d]
                vies[k] -= LR * g
    modelo = {
        "nome": "why-man",
        "parametros": n * DIM + n,
        "pesos": pesos,
        "vies": vies,
        "respostas": [narracao for _, narracao in EXEMPLOS],
    }
    with open(caminho, "w", encoding="utf-8") as f:
        json.dump(modelo, f, ensure_ascii=False)
    return modelo


def _carregar() -> dict:
    if not os.path.exists(PESOS):
        return treinar_e_salvar()
    with open(PESOS, "r", encoding="utf-8") as f:
        return json.load(f)


def narrar_local(fato: str) -> str:
    modelo = _carregar()
    x = vetor(fato)
    logits = [_dot(w, x) + b for w, b in zip(modelo["pesos"], modelo["vies"])]
    idx = max(range(len(logits)), key=lambda i: logits[i])
    base = modelo["respostas"][idx]
    if "crític" in fato.lower() or "critic" in fato.lower():
        if "Falha crítica" in fato or "Falha critica" in fato:
            return "[VERMELHO]" + base.replace("[VERDE]", "").replace("[AMARELO]", "")
    return base
