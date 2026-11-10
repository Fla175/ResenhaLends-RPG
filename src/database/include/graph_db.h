#ifndef RPG_DB_GRAPH_DB_H
#define RPG_DB_GRAPH_DB_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stddef.h>
#include <stdint.h>
#include <stdbool.h>

/* Banco de dados de grafos próprio, não-relacional, escrito em C99 puro.
 * Conforme ADR 004 (C2.0): segurança, robustez, desempenho e persistência binária. */

typedef struct GraphDB GraphDB;

/* Códigos de erro explícitos. Nenhuma operação chama abort/exit. */
typedef enum {
    GRAPH_DB_OK = 0,
    GRAPH_DB_ERR_NOT_FOUND = -1,
    GRAPH_DB_ERR_INVALID_ARG = -2,
    GRAPH_DB_ERR_IO = -3,
    GRAPH_DB_ERR_CORRUPT = -4,
    GRAPH_DB_ERR_NOMEM = -5,
    GRAPH_DB_ERR_LOCKED = -6,
    GRAPH_DB_ERR_TXN_ACTIVE = -7,
    GRAPH_DB_ERR_NO_TXN = -8,
    GRAPH_DB_ERR_SECURITY = -9,
    GRAPH_DB_ERR_UNSUPPORTED = -10
} graph_db_error_t;

/* Configuração de abertura. Todos os campos têm valores padrão seguros se NULL. */
typedef struct {
    size_t page_size;           /* bytes, múltiplo de 4096 recomendado (padrão 4096) */
    size_t cache_size;          /* número de páginas no cache LRU (padrão 64) */
    bool wal_enabled;           /* write-ahead log (padrão true) */
    bool read_only;             /* abre apenas para leitura (padrão false) */
    const char* encryption_key; /* opcional, NULL desativa criptografia de páginas */
} graph_db_config_t;

/* Métricas expostas para tuning contínuo (ADR 4.10). */
typedef struct {
    size_t cache_hits;
    size_t cache_misses;
    size_t pages_read;
    size_t pages_written;
    double avg_op_latency_us;
} graph_db_metrics_t;

/* Inicializa config com valores padrão. */
void graph_db_config_init(graph_db_config_t* cfg);

/* Cria ou abre um banco de dados em arquivo binário.
 * Retorna NULL e preenche out_err em caso de falha. */
GraphDB* graph_db_open(const char* path,
                       const graph_db_config_t* cfg,
                       graph_db_error_t* out_err);

/* Fecha o banco e libera recursos. Seguro chamar com NULL. */
void graph_db_close(GraphDB* db);

/* Insere um nó. Retorna true se sucesso e preenche out_node_id.
 * Valida label e propriedades antes de gravar (ADR 4.8). */
bool graph_db_add_node(GraphDB* db,
                       const char* label,
                       const void* properties,
                       size_t properties_len,
                       int64_t* out_node_id,
                       graph_db_error_t* out_err);

/* Insere uma aresta entre dois nós. Retorna true se sucesso. */
bool graph_db_add_edge(GraphDB* db,
                       int64_t from_node_id,
                       int64_t to_node_id,
                       const char* relation,
                       graph_db_error_t* out_err);

/* Inserção em lote de nós (ADR 4.7/4.10). Reduz chamadas de fsync. */
bool graph_db_add_nodes_batch(GraphDB* db,
                              const char** labels,
                              const void** properties,
                              const size_t* properties_lens,
                              size_t count,
                              int64_t* out_node_ids,
                              graph_db_error_t* out_err);

/* Busca nós por label. Retorna array alocado com IDs; *out_count recebe qtd.
 * O chamador deve liberar com graph_db_free_ids. */
int64_t* graph_db_find_nodes(GraphDB* db,
                             const char* label,
                             size_t* out_count,
                             graph_db_error_t* out_err);

/* Libera array alocado por graph_db_find_nodes. */
void graph_db_free_ids(int64_t* ids);

/* Transações: begin / commit / rollback (ADR 4.9). */
bool graph_db_begin(GraphDB* db, graph_db_error_t* out_err);
bool graph_db_commit(GraphDB* db, graph_db_error_t* out_err);
bool graph_db_rollback(GraphDB* db, graph_db_error_t* out_err);

/* Persiste alterações em disco (fsync de páginas e WAL). */
bool graph_db_flush(GraphDB* db, graph_db_error_t* out_err);

/* Segurança (ADR 4.8): valida e sanitiza identificadores. */
bool graph_db_validate_identifier(const char* ident, graph_db_error_t* out_err);
size_t graph_db_sanitize_identifier(char* dst,
                                    size_t dst_size,
                                    const char* src,
                                    graph_db_error_t* out_err);

/* Robustez (ADR 4.9): recuperação e verificação de integridade. */
bool graph_db_recover(const char* path, graph_db_error_t* out_err);
bool graph_db_verify_checksum(GraphDB* db, graph_db_error_t* out_err);

/* Métricas (ADR 4.10). */
void graph_db_get_metrics(GraphDB* db, graph_db_metrics_t* out);

#ifdef __cplusplus
}
#endif

#endif /* RPG_DB_GRAPH_DB_H */
