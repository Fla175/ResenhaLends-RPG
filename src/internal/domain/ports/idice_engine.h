#ifndef RPG_INTERNAL_DOMAIN_PORTS_IDICE_ENGINE_H
#define RPG_INTERNAL_DOMAIN_PORTS_IDICE_ENGINE_H

#include <cstdint>
#include <string>

namespace rpg {
namespace domain {
namespace ports {

class IDiceEngine {
public:
    virtual ~IDiceEngine() = default;

    virtual std::int64_t roll_d6() = 0;
    virtual std::int64_t roll_d20() = 0;
    virtual std::int64_t roll(const std::string& notation) = 0;
};

} // namespace ports
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_PORTS_IDICE_ENGINE_H
