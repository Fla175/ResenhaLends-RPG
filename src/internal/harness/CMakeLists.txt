# Build do Harness de IA L.I.A.N.D.R.A. (ADR 003)
# Sem dependência de Python em tempo de execução.

set(HARNESS_SRCS
    domain/entities/sessao.cpp
    domain/value_objects/texto.cpp
    domain/value_objects/sugestao_correcao.cpp
    domain/value_objects/fragmento_contexto.cpp
    domain/events/contexto_recuperado.cpp
    domain/events/correcao_aplicada.cpp
    application/use_cases/gerar_resposta_mestre.cpp
    application/use_cases/corrigir_texto.cpp
    application/use_cases/recuperar_contexto_rag.cpp
    adapters/correcao/ngram_correcao.cpp
    adapters/correcao/markov_correcao.cpp
    adapters/inferencia/libtorch_inferencia.cpp
    adapters/memoria/repositorio_sessao_memoria.cpp
    adapters/controllers/controlador_mestre.cpp
)

add_library(harness_ia STATIC ${HARNESS_SRCS})

target_include_directories(harness_ia PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/../..
)

target_link_libraries(harness_ia PUBLIC
    rpg_domain
    rpg_ports
)

# LibTorch é opcional e fica no Frameworks & Drivers.
# Descomentar quando o backend LibTorch estiver integrado:
# find_package(Torch REQUIRED)
# target_link_libraries(harness_ia PUBLIC ${TORCH_LIBRARIES})
# set_property(TARGET harness_ia PROPERTY CXX_STANDARD 17)
