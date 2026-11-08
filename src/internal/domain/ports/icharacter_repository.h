#ifndef RPG_INTERNAL_DOMAIN_PORTS_ICHARACTER_REPOSITORY_H
#define RPG_INTERNAL_DOMAIN_PORTS_ICHARACTER_REPOSITORY_H

#include <cstdint>
#include <memory>
#include <vector>

namespace rpg {
namespace domain {
namespace ports {

class Personagem;

class ICharacterRepository {
public:
    virtual ~ICharacterRepository() = default;

    virtual std::int64_t save(const Personagem& personagem) = 0;
    virtual std::unique_ptr<Personagem> find_by_id(std::int64_t id) = 0;
    virtual std::vector<std::unique_ptr<Personagem>> find_all() = 0;
    virtual bool remove(std::int64_t id) = 0;
};

} // namespace ports
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_PORTS_ICHARACTER_REPOSITORY_H
