# ui_pygame.py — mesa do mestre
import os
import sys
import pygame

import agente_estado
import agente_narrador
import dados
import gerenciador_jogadores as elenco
import gerenciador_json
import gerenciador_capitulos
import gerador_enc
import gerador_npc

LARGURA, ALTURA = 1200, 720
BG = (14, 18, 26)
PAINEL = (24, 32, 44)
PAINEL2 = (32, 42, 58)
BORDA = (70, 160, 175)
TEXTO = (228, 232, 238)
OURO = (232, 196, 96)
VERDE = (92, 196, 124)
VERMELHO = (220, 88, 96)
CINZA = (140, 150, 165)
MAGENTA = (196, 120, 200)

TAGS = {
    "[VERDE]": VERDE,
    "[VERMELHO]": VERMELHO,
    "[AMARELO]": OURO,
    "[MAGENTA]": MAGENTA,
    "[/]": TEXTO,
}


def _pasta_dialogo(salvar=True, zipfile=False):
    try:
        import tkinter as tk
        from tkinter import filedialog
        root = tk.Tk()
        root.withdraw()
        root.attributes("-topmost", True)
        if zipfile:
            caminho = filedialog.askopenfilename(
                title="Importar campanha (.zip)",
                filetypes=[("Zip", "*.zip")],
            )
        else:
            caminho = filedialog.askdirectory(title="Escolha a pasta de destino")
        root.destroy()
        return caminho or None
    except Exception:
        return None


def wrap(fonte, texto, largura):
    palavras = texto.split()
    linhas, atual = [], ""
    for p in palavras:
        teste = (atual + " " + p).strip()
        if fonte.size(teste)[0] <= largura:
            atual = teste
        else:
            if atual:
                linhas.append(atual)
            atual = p
    if atual:
        linhas.append(atual)
    return linhas or [""]


class Botao:
    def __init__(self, rect, rotulo, cor=BORDA, acao=None):
        self.rect = pygame.Rect(rect)
        self.rotulo = rotulo
        self.cor = cor
        self.acao = acao

    def desenhar(self, tela, fonte, mouse):
        hover = self.rect.collidepoint(mouse)
        pygame.draw.rect(tela, PAINEL2 if hover else PAINEL, self.rect, border_radius=8)
        pygame.draw.rect(tela, self.cor, self.rect, 2, border_radius=8)
        t = fonte.render(self.rotulo, True, TEXTO)
        tela.blit(t, t.get_rect(center=self.rect.center))

    def clicar(self, pos):
        if self.rect.collidepoint(pos) and self.acao:
            self.acao()
            return True
        return False


class Campo:
    def __init__(self, rect, texto="", placeholder=""):
        self.rect = pygame.Rect(rect)
        self.texto = texto
        self.placeholder = placeholder
        self.ativo = False

    def evento(self, ev):
        if ev.type == pygame.MOUSEBUTTONDOWN:
            self.ativo = self.rect.collidepoint(ev.pos)
        if not self.ativo or ev.type != pygame.KEYDOWN:
            return
        if ev.key == pygame.K_BACKSPACE:
            self.texto = self.texto[:-1]
        elif ev.key == pygame.K_RETURN:
            self.ativo = False
        else:
            if ev.unicode and ev.unicode.isprintable() and len(self.texto) < 180:
                self.texto += ev.unicode

    def desenhar(self, tela, fonte):
        pygame.draw.rect(tela, PAINEL2, self.rect, border_radius=6)
        pygame.draw.rect(tela, OURO if self.ativo else BORDA, self.rect, 2, border_radius=6)
        mostra = self.texto if self.texto else self.placeholder
        cor = TEXTO if self.texto else CINZA
        clip = fonte.render(mostra, True, cor)
        tela.set_clip(self.rect.inflate(-8, -4))
        tela.blit(clip, (self.rect.x + 8, self.rect.y + (self.rect.h - clip.get_height()) // 2))
        tela.set_clip(None)


class Mesa:
    def __init__(self):
        pygame.init()
        pygame.display.set_caption("ResenhaLends — Mesa do Mestre")
        self.tela = pygame.display.set_mode((LARGURA, ALTURA))
        self.clock = pygame.time.Clock()
        self.fonte = pygame.font.SysFont("segoeui", 18) or pygame.font.Font(None, 22)
        self.fonte_p = pygame.font.SysFont("segoeui", 15) or pygame.font.Font(None, 18)
        self.fonte_t = pygame.font.SysFont("segoeui", 26, bold=True) or pygame.font.Font(None, 32)
        self.tela_id = "inicio"
        self.campanha = None
        self.aba = "sessao"
        self.turno = 0
        self.log = ["Why-Man à mesa. Regras locais. Exporte a campanha quando quiser."]
        self.status = ""
        self.jogador_sel = None
        self.cap_sel = 0
        self.campos = {}
        self._montar_inicio()

    def _msg(self, texto):
        self.status = texto
        self.log.append(texto)
        self.log = self.log[-80:]

    def _montar_inicio(self):
        self.campos = {
            "busca": Campo((40, 640, 500, 36), placeholder="nome da nova campanha..."),
        }

    def _montar_nova(self):
        self.campos = {
            "nome": Campo((40, 140, 500, 36), placeholder="slug da campanha"),
            "desc": Campo((40, 200, 900, 36), placeholder="descrição"),
            "objetivo": Campo((40, 260, 900, 36), placeholder="objetivo principal"),
        }
        self._dif = "Normal"
        self._tam = "media"

    def _montar_elenco(self):
        j = self._jogador()
        att = (j or {}).get("atributos") or {}
        macros = (j or {}).get("macros") or {}
        inv = ", ".join((j or {}).get("inventario") or [])
        self.campos = {
            "nome": Campo((300, 90, 320, 32), (j or {}).get("nome", ""), "nome"),
            "classe": Campo((640, 90, 220, 32), (j or {}).get("classe", ""), "classe"),
            "raca": Campo((880, 90, 220, 32), (j or {}).get("raca", ""), "raça"),
            "desc": Campo((300, 134, 800, 32), (j or {}).get("descricao", ""), "descrição"),
            "vida": Campo((300, 178, 80, 32), str((j or {}).get("vida", 20)), "HP"),
            "inv": Campo((400, 178, 700, 32), inv, "inventário (vírgulas)"),
            "m1": Campo((300, 222, 260, 32), macros.get("1", ""), "macro 1"),
            "m2": Campo((580, 222, 260, 32), macros.get("2", ""), "macro 2"),
            "m3": Campo((860, 222, 240, 32), macros.get("3", ""), "macro 3"),
            "for": Campo((300, 266, 70, 32), str(att.get("for", 10)), "FOR"),
            "hab": Campo((380, 266, 70, 32), str(att.get("hab", 10)), "HAB"),
            "res": Campo((460, 266, 70, 32), str(att.get("res", 10)), "RES"),
            "arm": Campo((540, 266, 70, 32), str(att.get("arm", 10)), "ARM"),
            "pdf": Campo((620, 266, 70, 32), str(att.get("pdf", 10)), "PDF"),
        }

    def _montar_notas(self):
        notas = gerenciador_json.carregar(self.campanha, ".notas_mestre.json") or {}
        self.campos = {
            "anotacoes": Campo((40, 100, 1100, 220), notas.get("anotacoes", ""), "anotações da sessão"),
            "segredos": Campo((40, 360, 1100, 180), notas.get("segredos", ""), "segredos do mestre"),
        }

    def _montar_sessao(self):
        self.campos = {
            "acao": Campo((300, 640, 520, 40), placeholder="ação do jogador da vez..."),
        }

    def _jogadores(self):
        if not self.campanha:
            return []
        return elenco.listar(self.campanha)

    def _jogador(self):
        js = self._jogadores()
        if self.jogador_sel:
            for j in js:
                if j.get("nome") == self.jogador_sel:
                    return j
        return js[0] if js else None

    def _entrar(self, nome):
        gerenciador_json.migrar_bestiario(nome)
        notas = gerenciador_json.carregar(nome, ".notas_mestre.json")
        if notas is None:
            gerenciador_json.salvar(nome, ".notas_mestre.json", {"anotacoes": "", "segredos": ""})
        self.campanha = nome
        self.tela_id = "mesa"
        self.aba = "sessao"
        self.turno = 0
        js = self._jogadores()
        self.jogador_sel = js[0]["nome"] if js else None
        self.cap_sel = 0
        self._montar_sessao()
        ctx = gerenciador_json.carregar(nome, ".context.json") or {}
        prologo = gerenciador_capitulos.garantir_abertura(nome)
        limpo = prologo
        for t in TAGS:
            limpo = limpo.replace(t, "")
        self.log = [limpo, f"Sessão: {nome} — {ctx.get('objetivo_principal', '')}"]
        self._msg("O mestre descreveu o mundo. Capítulos na aba Capítulos.")

    def criar_campanha(self):
        nome = (self.campos["nome"].texto or "").strip().replace(" ", "_").lower()
        if not nome:
            self._msg("Dê um nome à campanha.")
            return
        gerenciador_json.setup_campanha(
            nome,
            self._dif,
            self.campos["desc"].texto.strip(),
            self.campos["objetivo"].texto.strip(),
            self._tam,
        )
        gerenciador_capitulos.garantir_abertura(nome)
        self._entrar(nome)

    def exportar(self):
        if not self.campanha:
            return
        destino = _pasta_dialogo() or gerenciador_json.PASTA_EXPORT
        try:
            arq = gerenciador_json.exportar_campanha(self.campanha, destino)
            self._msg(f"Exportado: {arq}")
        except Exception as e:
            self._msg(f"Falha ao exportar: {e}")

    def importar(self):
        zip_path = _pasta_dialogo(zipfile=True)
        if not zip_path:
            return
        try:
            gerenciador_json.importar_campanha(zip_path)
            self._msg(f"Importado de {os.path.basename(zip_path)}")
        except Exception as e:
            self._msg(f"Falha ao importar: {e}")

    def resolver_acao(self, acao=None):
        js = self._jogadores()
        if not js:
            self._msg("Cadastre um jogador no Elenco.")
            return
        j = js[self.turno % len(js)]
        acao = (acao or self.campos.get("acao").texto or "").strip()
        if not acao:
            self._msg("Escreva a ação.")
            return
        macros = j.get("macros") or {}
        if acao in macros and macros[acao]:
            acao = macros[acao]
        dado = dados.rolar(mostrar=False)
        self.ultimo_dado = dado
        desfecho = agente_estado.processar_acao(
            self.campanha, f"O jogador {j['nome']} tenta: {acao}", dado
        )
        cena, comentario, fechado = agente_narrador.narrar(self.campanha, desfecho, imprimir=False)
        self._msg(f"D20: {dado} — {desfecho}")
        self._msg(cena.replace("[VERDE]", "").replace("[VERMELHO]", "").replace("[AMARELO]", "").replace("[MAGENTA]", "").replace("[/]", ""))
        if comentario:
            self._msg(comentario)
        if fechado:
            self._msg(f"Capítulo fechado — {fechado.get('titulo')}")
        if "acao" in self.campos:
            self.campos["acao"].texto = ""
        self.turno = (self.turno + 1) % max(1, len(self._jogadores()))
        nxt = self._jogadores()
        if nxt:
            self.jogador_sel = nxt[self.turno % len(nxt)]["nome"]

    def encontro(self):
        evento = gerador_enc.gerar_encontro(self.campanha)
        cena, com, fechado = agente_narrador.narrar(
            self.campanha, f"Um evento inesperado acontece: {evento}", imprimir=False
        )
        self._msg(evento)
        self._msg(cena.replace("[VERDE]", "").replace("[VERMELHO]", "").replace("[AMARELO]", "").replace("[MAGENTA]", "").replace("[/]", ""))
        if com:
            self._msg(com)
        if fechado:
            self._msg(f"Capítulo fechado — {fechado.get('titulo')}")

    def npc(self):
        n = gerador_npc.gerar_npc(self.campanha)
        des = f"O grupo encontra {n['nome']}, um(a) {n['raca']} {n['ocupacao']}. Motivação: {n['motivacao']}."
        cena, com, fechado = agente_narrador.narrar(self.campanha, des, imprimir=False)
        self._msg(des)
        self._msg(cena.replace("[VERDE]", "").replace("[VERMELHO]", "").replace("[AMARELO]", "").replace("[MAGENTA]", "").replace("[/]", ""))
        if com:
            self._msg(com)
        if fechado:
            self._msg(f"Capítulo fechado — {fechado.get('titulo')}")

    def salvar_elenco(self, novo=False):
        c = self.campos
        try:
            vida = int(c["vida"].texto or 20)
        except ValueError:
            vida = 20
        atts = {}
        for k in elenco.ATRIBUTOS:
            try:
                atts[k] = int(c[k].texto or 10)
            except ValueError:
                atts[k] = 10
        patch = {
            "nome": c["nome"].texto.strip(),
            "classe": c["classe"].texto.strip(),
            "raca": c["raca"].texto.strip(),
            "descricao": c["desc"].texto.strip(),
            "vida": vida,
            "inventario": c["inv"].texto,
            "macros": {"1": c["m1"].texto, "2": c["m2"].texto, "3": c["m3"].texto},
            "atributos": atts,
        }
        try:
            if novo:
                elenco.adicionar(self.campanha, patch)
                self.jogador_sel = patch["nome"]
                self._msg(f"{patch['nome']} entrou na mesa.")
            else:
                alvo = self.jogador_sel or patch["nome"]
                elenco.atualizar(self.campanha, alvo, patch)
                self.jogador_sel = patch["nome"]
                self._msg(f"{patch['nome']} atualizado.")
            self._montar_elenco()
        except Exception as e:
            self._msg(str(e))

    def rolar_atts(self):
        if not self.jogador_sel:
            return
        novos = {k: dados.rolar(mostrar=False) for k in elenco.ATRIBUTOS}
        elenco.atualizar(self.campanha, self.jogador_sel, {"atributos": novos})
        self._montar_elenco()
        self._msg(f"Atributos rolados: {novos}")

    def salvar_notas(self):
        gerenciador_json.salvar(self.campanha, ".notas_mestre.json", {
            "anotacoes": self.campos["anotacoes"].texto,
            "segredos": self.campos["segredos"].texto,
        })
        self._msg("Notas do mestre salvas.")

    def set_dificuldade(self, valor):
        regras = gerenciador_json.carregar(self.campanha, ".regras.json") or {}
        regras["dificuldade"] = valor
        gerenciador_json.salvar(self.campanha, ".regras.json", regras)
        self._msg(f"Dificuldade: {valor}")

    def botoes_inicio(self):
        ys = 160
        btns = []
        for i, nome in enumerate(gerenciador_json.listar_campanhas()):
            r = (40, ys + i * 48, 520, 40)
            btns.append(Botao(r, nome, acao=lambda n=nome: self._entrar(n)))
        btns += [
            Botao((600, 160, 240, 44), "Nova campanha", OURO, lambda: self._ir("nova")),
            Botao((600, 216, 240, 44), "Importar zip", acao=self.importar),
            Botao((600, 272, 240, 44), "Sair", VERMELHO, lambda: sys.exit(0)),
        ]
        return btns

    def _ir(self, tela):
        self.tela_id = tela
        if tela == "nova":
            self._montar_nova()
        elif tela == "inicio":
            self.campanha = None
            self._montar_inicio()

    def botoes_mesa(self):
        abas = [
            ("sessao", "Sessão"),
            ("elenco", "Elenco"),
            ("capitulos", "Capítulos"),
            ("notas", "Notas"),
            ("regras", "Regras"),
        ]
        btns = []
        for i, (k, lab) in enumerate(abas):
            btns.append(Botao((20 + i * 112, 58, 106, 34), lab, OURO if self.aba == k else BORDA,
                              lambda kk=k: self._trocar_aba(kk)))
        btns += [
            Botao((1040, 16, 140, 36), "Exportar", OURO, self.exportar),
            Botao((890, 16, 140, 36), "Menu", acao=lambda: self._ir("inicio")),
        ]
        if self.aba == "sessao":
            btns += [
                Botao((830, 640, 110, 40), "Resolver", VERDE, lambda: self.resolver_acao()),
                Botao((950, 640, 110, 40), "Encontro", VERMELHO, self.encontro),
                Botao((1070, 640, 110, 40), "NPC", MAGENTA, self.npc),
            ]
            j = self._jogador()
            if j:
                macros = j.get("macros") or {}
                x = 300
                for k, v in macros.items():
                    if v:
                        btns.append(Botao((x, 590, 160, 36), f"{k}:{v[:12]}", acao=lambda vv=v: self.resolver_acao(vv)))
                        x += 170
        if self.aba == "elenco":
            btns += [
                Botao((300, 320, 160, 40), "Salvar ficha", VERDE, lambda: self.salvar_elenco(False)),
                Botao((480, 320, 160, 40), "Novo jogador", OURO, lambda: self.salvar_elenco(True)),
                Botao((660, 320, 160, 40), "Rolar atts", acao=self.rolar_atts),
                Botao((840, 320, 80, 40), "HP+1", VERDE, lambda: self._hp(1)),
                Botao((930, 320, 80, 40), "HP-1", VERMELHO, lambda: self._hp(-1)),
                Botao((1020, 320, 80, 40), "Apagar", VERMELHO, self._apagar),
            ]
        if self.aba == "notas":
            btns.append(Botao((40, 560, 200, 40), "Salvar notas", VERDE, self.salvar_notas))
        if self.aba == "capitulos":
            btns.append(Botao((20, 640, 220, 40), "Encerrar capítulo", OURO, self.encerrar_capitulo))
        if self.aba == "regras":
            for i, d in enumerate(("Fácil", "Normal", "hardcore")):
                btns.append(Botao((40 + i * 180, 140, 160, 44), d, acao=lambda dd=d: self.set_dificuldade(dd)))
        return btns

    def encerrar_capitulo(self):
        fechado = gerenciador_capitulos.encerrar(self.campanha)
        if not fechado:
            self._msg("Nada para fechar — joguem um pouco.")
            return
        self._msg(f"Capítulo fechado — {fechado.get('titulo')}")
        self.cap_sel = len(gerenciador_capitulos.listar(self.campanha)) - 1

    def _trocar_aba(self, aba):
        self.aba = aba
        if aba == "sessao":
            self._montar_sessao()
        elif aba == "elenco":
            self._montar_elenco()
        elif aba == "notas":
            self._montar_notas()
        else:
            self.campos = {}

    def _hp(self, delta):
        if not self.jogador_sel:
            return
        try:
            elenco.ajustar_vida(self.campanha, self.jogador_sel, delta)
            self._montar_elenco()
        except Exception as e:
            self._msg(str(e))

    def _apagar(self):
        if not self.jogador_sel:
            return
        elenco.remover(self.campanha, self.jogador_sel)
        js = self._jogadores()
        self.jogador_sel = js[0]["nome"] if js else None
        self._montar_elenco()
        self._msg("Jogador removido.")

    def desenhar_texto(self, s, xy, fonte=None, cor=TEXTO):
        fonte = fonte or self.fonte
        self.tela.blit(fonte.render(s, True, cor), xy)

    def loop(self):
        while True:
            mouse = pygame.mouse.get_pos()
            if self.tela_id == "inicio":
                botoes = self.botoes_inicio()
            elif self.tela_id == "nova":
                botoes = [
                    Botao((40, 340, 140, 40), "Fácil", acao=lambda: setattr(self, "_dif", "Fácil")),
                    Botao((200, 340, 140, 40), "Normal", acao=lambda: setattr(self, "_dif", "Normal")),
                    Botao((360, 340, 140, 40), "Hardcore", acao=lambda: setattr(self, "_dif", "hardcore")),
                    Botao((40, 400, 140, 40), "Pequena", acao=lambda: setattr(self, "_tam", "pequena")),
                    Botao((200, 400, 140, 40), "Média", acao=lambda: setattr(self, "_tam", "media")),
                    Botao((360, 400, 140, 40), "Extensa", acao=lambda: setattr(self, "_tam", "extensa")),
                    Botao((40, 480, 200, 48), "Criar e abrir", VERDE, self.criar_campanha),
                    Botao((260, 480, 160, 48), "Voltar", acao=lambda: self._ir("inicio")),
                ]
            else:
                botoes = self.botoes_mesa()

            for ev in pygame.event.get():
                if ev.type == pygame.QUIT:
                    pygame.quit()
                    return
                for c in self.campos.values():
                    c.evento(ev)
                if ev.type == pygame.MOUSEBUTTONDOWN and ev.button == 1:
                    if self.tela_id == "mesa" and self.aba == "elenco":
                        for i, j in enumerate(self._jogadores()):
                            r = pygame.Rect(20, 90 + i * 44, 260, 40)
                            if r.collidepoint(ev.pos):
                                self.jogador_sel = j.get("nome")
                                self._montar_elenco()
                    if self.tela_id == "mesa" and self.aba == "sessao":
                        for i, j in enumerate(self._jogadores()):
                            r = pygame.Rect(20, 110 + i * 70, 260, 64)
                            if r.collidepoint(ev.pos):
                                self.jogador_sel = j.get("nome")
                                self.turno = i
                    if self.tela_id == "mesa" and self.aba == "capitulos":
                        for i, _c in enumerate(gerenciador_capitulos.listar(self.campanha)):
                            r = pygame.Rect(20, 100 + i * 44, 320, 40)
                            if r.collidepoint(ev.pos):
                                self.cap_sel = i
                    for b in botoes:
                        if b.clicar(ev.pos):
                            break
                if ev.type == pygame.KEYDOWN and ev.key == pygame.K_RETURN and self.tela_id == "mesa" and self.aba == "sessao":
                    self.resolver_acao()

            self.tela.fill(BG)
            if self.tela_id == "inicio":
                self.desenhar_texto("RESENHALENDS", (40, 30), self.fonte_t, OURO)
                self.desenhar_texto("Mesa do mestre — campanhas", (40, 110), cor=CINZA)
                self.desenhar_texto(self.status, (40, 680), self.fonte_p, CINZA)
            elif self.tela_id == "nova":
                self.desenhar_texto("Nova campanha", (40, 40), self.fonte_t, OURO)
                self.desenhar_texto(f"Dificuldade: {self._dif}   Tamanho: {self._tam}", (40, 90), cor=CINZA)
            else:
                ctx = gerenciador_json.carregar(self.campanha, ".context.json") or {}
                self.desenhar_texto(self.campanha.upper(), (20, 12), self.fonte_t, OURO)
                self.desenhar_texto(ctx.get("objetivo_principal", ""), (20, 40), self.fonte_p, CINZA)
                if self.aba == "sessao":
                    pygame.draw.rect(self.tela, PAINEL, (20, 100, 260, 470), border_radius=10)
                    for i, j in enumerate(self._jogadores()):
                        r = pygame.Rect(28, 110 + i * 70, 244, 64)
                        ativo = j.get("nome") == self.jogador_sel
                        pygame.draw.rect(self.tela, PAINEL2 if ativo else BG, r, border_radius=8)
                        pygame.draw.rect(self.tela, VERDE if ativo else BORDA, r, 2, border_radius=8)
                        hp = int(j.get("vida", 20))
                        self.desenhar_texto(j.get("nome", "?"), (36, 118 + i * 70), self.fonte_p, OURO)
                        self.desenhar_texto(f"{j.get('classe') or 'avent'}  HP {hp}", (36, 140 + i * 70), self.fonte_p, VERDE if hp > 6 else VERMELHO)
                    pygame.draw.rect(self.tela, PAINEL, (300, 100, 880, 470), border_radius=10)
                    y = 112
                    for linha in self.log[-14:]:
                        for w in wrap(self.fonte_p, linha, 850)[:3]:
                            self.desenhar_texto(w, (316, y), self.fonte_p)
                            y += 20
                    if self.ultimo_dado is not None:
                        self.desenhar_texto(f"Último dado: {self.ultimo_dado}", (300, 548), cor=OURO)
                elif self.aba == "elenco":
                    for i, j in enumerate(self._jogadores()):
                        r = pygame.Rect(20, 90 + i * 44, 260, 40)
                        pygame.draw.rect(self.tela, PAINEL2 if j.get("nome") == self.jogador_sel else PAINEL, r, border_radius=8)
                        self.desenhar_texto(f"{j.get('nome')}  HP {j.get('vida')}", (32, 98 + i * 44), self.fonte_p)
                    self.desenhar_texto("Ficha do herói (mestre edita tudo)", (300, 60), self.fonte_p, CINZA)
                elif self.aba == "regras":
                    regras = gerenciador_json.carregar(self.campanha, ".regras.json") or {}
                    self.desenhar_texto(f"Atual: {regras}", (40, 100), cor=CINZA)
                    self.desenhar_texto("CD: Fácil 8 · Normal 12 · hardcore 15", (40, 210), cor=CINZA)
                    hist = gerenciador_json.carregar(self.campanha, ".hist.json") or []
                    self.desenhar_texto("Últimos desfechos:", (40, 260), cor=OURO)
                    y = 290
                    for h in hist[-8:]:
                        for w in wrap(self.fonte_p, str(h), 1000)[:2]:
                            self.desenhar_texto(w, (40, y), self.fonte_p)
                            y += 20
                elif self.aba == "capitulos":
                    caps = gerenciador_capitulos.listar(self.campanha)
                    atual = gerenciador_capitulos.atual(self.campanha)
                    self.desenhar_texto(
                        f"Em andamento: cap. {atual.get('numero')} ({len(atual.get('eventos') or [])} cenas)",
                        (600, 62), self.fonte_p, CINZA)
                    for i, cap in enumerate(caps):
                        r = pygame.Rect(20, 100 + i * 44, 320, 40)
                        pygame.draw.rect(self.tela, PAINEL2 if i == self.cap_sel else PAINEL, r, border_radius=8)
                        self.desenhar_texto((cap.get("titulo") or "?")[:42], (32, 108 + i * 44), self.fonte_p, OURO)
                    pygame.draw.rect(self.tela, PAINEL, (360, 100, 820, 520), border_radius=10)
                    if caps:
                        i = min(self.cap_sel, len(caps) - 1)
                        cap = caps[i]
                        y = 116
                        self.desenhar_texto(cap.get("titulo", ""), (376, y), self.fonte, OURO)
                        y += 28
                        self.desenhar_texto(cap.get("fechado_em") or cap.get("tipo", ""), (376, y), self.fonte_p, CINZA)
                        y += 24
                        for w in wrap(self.fonte_p, cap.get("resumo") or "", 780)[:8]:
                            self.desenhar_texto(w, (376, y), self.fonte_p)
                            y += 20
                        y += 8
                        self.desenhar_texto("Histórico da etapa:", (376, y), self.fonte_p, OURO)
                        y += 22
                        for evn in (cap.get("eventos") or [])[-8:]:
                            for w in wrap(self.fonte_p, str(evn), 780)[:2]:
                                self.desenhar_texto(w, (376, y), self.fonte_p, CINZA)
                                y += 18
                elif self.aba == "notas":
                    self.desenhar_texto("Só o mestre vê. Salve antes de exportar.", (40, 70), self.fonte_p, CINZA)

            for c in self.campos.values():
                c.desenhar(self.tela, self.fonte)
            for b in botoes:
                b.desenhar(self.tela, self.fonte_p, mouse)
            if self.tela_id == "mesa":
                self.desenhar_texto(self.status, (20, 690), self.fonte_p, CINZA)

            pygame.display.flip()
            self.clock.tick(30)


def executar():
    Mesa().loop()
