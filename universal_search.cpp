#include "universal_search.h"

vec<char> from_string(const char* s) {
    vec<char> v;
    for(uint64 i=0; i<std::strlen(s); i++) {
        v[i] = s[i];
    }
    return v;
}

int main() {
    vec<char> truth = from_string("2x2");

    auto [i, o] = universal_search(from_string("4"), [&](const vec<char>& out) -> bool {
        if (out.size() != truth.size()) {
            return false;
        }
        for (uint64 i=0; i<truth.size(); i++) {
            if (out[i] != truth[i]) {
                return false;
            }
        }
        return true;
    });
}