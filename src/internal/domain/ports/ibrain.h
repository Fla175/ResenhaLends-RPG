#ifndef RPG_INTERNAL_DOMAIN_PORTS_IBRAIN_H
#define RPG_INTERNAL_DOMAIN_PORTS_IBRAIN_H

#include <cstdint>
#include <string>

namespace rpg {
namespace domain {
namespace ports {

class Personagem;

class IBrain {
public:
    virtual ~IBrain() = default;

    virtual std::string decide_action(const Personagem& npc) = 0;
    virtual void set_personality_trait(const std::string& trait, std::int8_t value) = 0;
    virtual std::int8_t get_personality_trait(const std::string& trait) const = 0;
};

} // namespace ports
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_PORTS_IBRAIN_H
