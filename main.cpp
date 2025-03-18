#include <cstdio>
#include <vector>
#include <unordered_map>
#include <cstdint>

using uint64 = std::uint64_t;
template<typename T>
using vec = std::vector<T>;
template<typename KT, typename VT>
using dict = std::unordered_map<KT, VT>;

struct prog {
    const char* code;
    dict<uint64, uint64> jump_table;
    prog(const char* code) : code(code) {
        vec<uint64> bracket_index_stack;
        for (uint64 i=0; i<std::strlen(code); i++) {
            switch (code[i]) {
                case '[':
                    bracket_index_stack.push_back(i);
                    continue;
                case ']':
                    uint64 j = bracket_index_stack.back();
                    bracket_index_stack.pop_back();
                    jump_table[i] = j;
                    jump_table[j] = i;
                    continue;
            }
        }
    }
    void print_jump_table() {
        for (auto& [k, v] : jump_table) {
            std::printf("%llu -> %llu\n", k, v);
        }
    }
};


int main() {
    prog p("[caca][cac[ccc]]");
    p.print_jump_table();
    return 0;
}