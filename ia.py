# ia.py — cliente Ollama com diagnóstico, timeout e limites para CPU
import os
from typing import Iterator, Optional

import httpx
from ollama import Client, ResponseError

MODELO = os.getenv("OLLAMA_MODEL", "llama3.2:1b")
HOST = os.getenv("OLLAMA_HOST", "http://127.0.0.1:11434")
TIMEOUT = float(os.getenv("OLLAMA_TIMEOUT", "300"))
KEEP_ALIVE = os.getenv("OLLAMA_KEEP_ALIVE", "30m")

_client: Optional[Client] = None


def cliente() -> Client:
    global _client
    if _client is None:
        _client = Client(host=HOST, timeout=TIMEOUT)
    return _client


def _opcoes(temperature: float, num_predict: int) -> dict:
    return {
        "temperature": temperature,
        "num_ctx": 2048,
        "num_predict": num_predict,
    }


def traduzir_erro(erro: Exception) -> str:
    if isinstance(erro, (httpx.ConnectError, httpx.ConnectTimeout, ConnectionError)):
        return (
            f"Não foi possível conectar ao Ollama em {HOST}. "
            "Abra o aplicativo Ollama e tente de novo."
        )
    if isinstance(erro, (httpx.ReadTimeout, TimeoutError)):
        return (
            "A IA demorou demais para responder. O modelo está rodando na CPU e fica lento. "
            "Aguarde e tente outra ação, ou feche outros programas pesados."
        )
    if isinstance(erro, ResponseError):
        texto = str(erro).lower()
        if "not found" in texto or "404" in texto:
            return f"Modelo '{MODELO}' não encontrado. No terminal, rode: ollama pull {MODELO}"
        return f"Ollama recusou o pedido: {erro}"
    return str(erro)


def verificar_conexao() -> tuple[bool, str]:
    try:
        lista = cliente().list()
        nomes = []
        for item in getattr(lista, "models", []) or []:
            nome = getattr(item, "model", None) or getattr(item, "name", None)
            if nome:
                nomes.append(nome)
        if not nomes:
            return False, "Ollama respondeu, mas nenhum modelo está instalado. Rode: ollama pull llama3.2:1b"
        if not any(MODELO in n or n.startswith(MODELO.split(":")[0]) for n in nomes):
            return False, (
                f"Modelo '{MODELO}' não está na lista ({', '.join(nomes)}). "
                f"Rode: ollama pull {MODELO}"
            )
        return True, f"Ollama ok | modelo {MODELO} | host {HOST}"
    except Exception as e:
        return False, traduzir_erro(e)


def gerar(prompt: str, *, json_mode: bool = False, temperature: float = 0.7, num_predict: int = 220) -> str:
    kwargs = {
        "model": MODELO,
        "prompt": prompt,
        "options": _opcoes(temperature, num_predict),
        "keep_alive": KEEP_ALIVE,
    }
    if json_mode:
        kwargs["format"] = "json"
    try:
        resposta = cliente().generate(**kwargs)
        return getattr(resposta, "response", None) or resposta.get("response", "")
    except Exception as e:
        raise RuntimeError(traduzir_erro(e)) from e


def gerar_stream(prompt: str, *, temperature: float = 0.9, num_predict: int = 220) -> Iterator[str]:
    try:
        fluxo = cliente().generate(
            model=MODELO,
            prompt=prompt,
            stream=True,
            options=_opcoes(temperature, num_predict),
            keep_alive=KEEP_ALIVE,
        )
        for pedaco in fluxo:
            texto = getattr(pedaco, "response", None)
            if texto is None and isinstance(pedaco, dict):
                texto = pedaco.get("response", "")
            if texto:
                yield texto
    except Exception as e:
        raise RuntimeError(traduzir_erro(e)) from e
