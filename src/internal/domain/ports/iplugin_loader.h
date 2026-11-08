#ifndef RPG_INTERNAL_DOMAIN_PORTS_IPLUGIN_LOADER_H
#define RPG_INTERNAL_DOMAIN_PORTS_IPLUGIN_LOADER_H

#include <string>

namespace rpg {
namespace domain {
namespace ports {

class IPluginLoader {
public:
    virtual ~IPluginLoader() = default;

    virtual bool load_plugin(const std::string& path) = 0;
    virtual void unload_plugin(const std::string& name) = 0;
};

} // namespace ports
} // namespace domain
} // namespace rpg

#endif // RPG_INTERNAL_DOMAIN_PORTS_IPLUGIN_LOADER_H
