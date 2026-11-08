#ifndef RPG_DB_GRAPH_DB_H
#define RPG_DB_GRAPH_DB_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

/* Banco de dados de grafos próprio, não-relacional, escrito em C99 puro.
 * Conforme ADR 001: controle manual de ponteiros e persistência em arquivos binários. */

typedef struct GraphDB GraphDB;

/* Cria ou abre um banco de dados em arquivo binário */
GraphDB* graph_db_open(const char* path);

/* Fecha o banco e libera recursos */
void graph_db_close(GraphDB* db);

/* Insere um nó. Retorna o ID do nó ou -1 em falha. */
int64_t graph_db_add_node(GraphDB* db, const char* label, const void* properties, size_t properties_len);

/* Insere uma aresta entre dois nós. Retorna true se sucesso. */
bool graph_db_add_edge(GraphDB* db, int64_t from_node_id, int64_t to_node_id, const char* relation);

/* Busca nós por label. Retorna array alocado; *out_count recebe a quantidade. */
int64_t* graph_db_find_nodes(GraphDB* db, const char* label, size_t* out_count);

/* Persiste alterações em disco */
bool graph_db_flush(GraphDB* db);

#ifdef __cplusplus
}
#endif

#endif /* RPG_DB_GRAPH_DB_H */
