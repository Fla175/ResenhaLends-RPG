#ifndef RPG_INTERNAL_DOMAIN_PORTS_IRENDERER_H
#define RPG_INTERNAL_DOMAIN_PORTS_IRENDERER_H

#include <string>

namespace rpg {
namespace domain {
namespace ports {

class IRenderer {
public:
    virtual ~IRenderer() = default;

    virtual void render_text(const std::string& message) = 0;
    virtual void clear_screen() = 0;
};

} // namespace ports
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_PORTS_IRENDERER_H
