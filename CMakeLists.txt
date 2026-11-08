cmake_minimum_required(VERSION 3.18 FATAL_ERROR)
project(ResenhaLends VERSION 1.0.0 LANGUAGES C CXX)

# Padrões de linguagem conforme ADR 001
set(CMAKE_C_STANDARD 99)
set(CMAKE_C_STANDARD_REQUIRED ON)
set(CMAKE_C_EXTENSIONS OFF)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if(NOT CMAKE_BUILD_TYPE)
    set(CMAKE_BUILD_TYPE Release CACHE STRING "Tipo de build" FORCE)
endif()

# --- Opções de Interface (Views) ---
# Conforme ADR: 4 variantes independentes desacopladas do núcleo
option(RESSENHALENDS_OMEGA  "GUI: raylib + Dear ImGui" OFF)
option(RESSENHALENDS_DELTA  "TUI estilizada: FTXUI"    OFF)
option(RESSENHALENDS_IOTA   "TUI leve: ncurses"        OFF)
option(RESSENHALENDS_EPSILON "TUI ultra-leve: cout + ANSI" OFF)

# Caso nenhuma interface seja escolhida, ativa Epsilon por padrão
if(NOT RESSENHALENDS_OMEGA AND NOT RESSENHALENDS_DELTA AND NOT RESSENHALENDS_IOTA AND NOT RESSENHALENDS_EPSILON)
    set(RESSENHALENDS_EPSILON ON CACHE BOOL "Interface padrão ultra-leve" FORCE)
endif()

# --- Dependências externas ---
find_package(CURL REQUIRED)

# Sol3 (Sol2) é header-only; disponibilizado via FetchContent para não exigir instalação prévia
include(FetchContent)
FetchContent_Declare(
    sol3
    GIT_REPOSITORY https://github.com/ThePhD/sol2.git
    GIT_TAG        v3.3.0
)
FetchContent_MakeAvailable(sol3)

# --- Módulos do núcleo ---
add_subdirectory(src/database)
add_subdirectory(src/core)
add_subdirectory(src/ai)
add_subdirectory(src/ui)

# --- Executável principal ---
add_executable(resenhalends src/main.cpp)

target_link_libraries(resenhalends PRIVATE
    rpg_db
    rpg_core
    ai_harness
)

target_include_directories(resenhalends PRIVATE
    src/core/include
    src/database/include
    src/ai/include
)

target_compile_features(resenhalends PRIVATE cxx_std_17)

# Vincula a interface selecionada
if(RESSENHALENDS_OMEGA)
    target_link_libraries(resenhalends PRIVATE ui_omega)
    target_compile_definitions(resenhalends PRIVATE RESSENHALENDS_UI_OMEGA)
elseif(RESSENHALENDS_DELTA)
    target_link_libraries(resenhalends PRIVATE ui_delta)
    target_compile_definitions(resenhalends PRIVATE RESSENHALENDS_UI_DELTA)
elseif(RESSENHALENDS_IOTA)
    target_link_libraries(resenhalends PRIVATE ui_iota)
    target_compile_definitions(resenhalends PRIVATE RESSENHALENDS_UI_IOTA)
elseif(RESSENHALENDS_EPSILON)
    target_link_libraries(resenhalends PRIVATE ui_epsilon)
    target_compile_definitions(resenhalends PRIVATE RESSENHALENDS_UI_EPSILON)
endif()

# Saída de binários
set_target_properties(resenhalends PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY ${CMAKE_BINARY_DIR}/bin
)
