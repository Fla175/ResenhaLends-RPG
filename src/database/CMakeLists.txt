# Banco de dados de grafos próprio em C99 puro (ADR 001, item 2)
add_library(rpg_db STATIC
    src/graph_db.c
)

target_include_directories(rpg_db PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}/include
)

target_compile_features(rpg_db PUBLIC c_std_99)

set_target_properties(rpg_db PROPERTIES
    LANGUAGE C
    ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/lib
)
