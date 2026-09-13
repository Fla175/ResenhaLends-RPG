# dados.py
import random
from colorama import Fore, Style

def rolar(faces=20, mostrar=True):
    """Rola um dado de N faces e, se mostrar=True, exibe no terminal."""
    if not isinstance(faces, int) or faces < 2:
        raise ValueError("O número de faces deve ser um inteiro maior que 1.")

    resultado = random.randint(1, faces)

    if mostrar:
        if resultado == faces:
            cor = Fore.GREEN
        elif resultado == 1:
            cor = Fore.RED
        else:
            cor = Fore.YELLOW
        print(f"{cor}🎲 D{faces}: {resultado}{Style.RESET_ALL}")

    return resultado

# Teste rápido se executar o arquivo diretamente
if __name__ == "__main__":
    from colorama import init
    init(autoreset=True)
    rolar(20) # Teste normal