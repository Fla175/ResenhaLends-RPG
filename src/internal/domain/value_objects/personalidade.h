#ifndef RPG_INTERNAL_DOMAIN_VALUE_OBJECTS_PERSONALIDADE_H
#define RPG_INTERNAL_DOMAIN_VALUE_OBJECTS_PERSONALIDADE_H

#include <cstdint>
#include <map>
#include <string>
#include "atributo.h"

namespace rpg {
namespace domain {
namespace value_objects {

class Personalidade {
public:
    void set_traco(const std::string& nome, std::int8_t valor) {
        tracos_[nome] = Atributo(valor);
    }

    Atributo get_traco(const std::string& nome) const {
        auto it = tracos_.find(nome);
        if (it == tracos_.end()) {
            return Atributo{0};
        }
        return it->second;
    }

private:
    std::map<std::string, Atributo> tracos_;
};

} // namespace value_objects
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_VALUE_OBJECTS_PERSONALIDADE_H
