#ifndef RPG_INTERNAL_DOMAIN_VALUE_OBJECTS_ATRIBUTO_H
#define RPG_INTERNAL_DOMAIN_VALUE_OBJECTS_ATRIBUTO_H

#include <cstdint>

namespace rpg {
namespace domain {
namespace value_objects {

class Atributo {
public:
    explicit Atributo(std::int8_t value = 0)
        : value_(clamp(value)) {}

    std::int8_t value() const { return value_; }

    bool operator==(const Atributo& other) const { return value_ == other.value_; }
    bool operator!=(const Atributo& other) const { return !(*this == other); }

private:
    static std::int8_t clamp(std::int8_t v) {
        return (v < 0) ? 0 : ((v > 100) ? 100 : v);
    }

    std::int8_t value_;
};

} // namespace value_objects
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_VALUE_OBJECTS_ATRIBUTO_H
