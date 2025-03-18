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
    vec<char> code;
    dict<uint64, uint64> jump_table;
    prog(const vec<char>& code) : code(code) {
        uint64 i = 0;
        while (i < code.size()) {
            if (code[i] == '[') {
                uint64 j = i;
                uint64 depth = 1;
                while (depth > 0) {
                    j++;
                    if (code[j] == '[') {
                        depth++;
                    } else if (code[j] == ']') {
                        depth--;
                    }
                }
                jump_table[i] = j;
                jump_table[j] = i;
            }
            i++;
        }
    }
    void print_jump_table() {
        for (auto& [k, v] : jump_table) {
            printf("%lu -> %lu\n", k, v);
        }
    }
};


int main() {
    char* code_str = "[caca][cac[ccc]]";
    vec<char> code(code_str, code_str + std::strlen(code_str));;
    prog p(code);
    p.print_jump_table();
    return 0;
}