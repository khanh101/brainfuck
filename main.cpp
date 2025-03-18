#include <cstdio>
#include <vector>
#include <unordered_map>
#include <cstdint>

using uint64 = std::uint64_t;
template<typename T>
using vec = std::vector<T>;
template<typename KT, typename VT>
using dict = std::unordered_map<KT, VT>;

struct interpreter {
    uint64 data_ptr;
    uint64 code_ptr;
    vec<char> data;
    interpreter(uint64 count): data_ptr(0), code_ptr(0), data(count, 0) {}

    const char* code;
    dict<uint64, uint64> jump_table;
    void load_code(const char* code, bool debug = true) {
        vec<uint64> bracket_index_stack;
        for (uint64 i=0; i<std::strlen(code); i++) {
            switch (code[i]) {
                case '[':
                    bracket_index_stack.push_back(i);
                    break;
                case ']':
                    uint64 j = bracket_index_stack.back();
                    bracket_index_stack.pop_back();
                    jump_table[i] = j;
                    jump_table[j] = i;
                    break;
            }
        }
        if (debug) {
            std::printf("jump table:\n");
            for (auto& [k, v] : jump_table) {
                if (k < v) {
                    std::printf("\t%llu <-> %llu\n", k, v);
                }
            }
        }
    }

    void step() {
        switch (code[code_ptr]) {
            case '>':
                data_ptr++;
                break;
            case '<':
                data_ptr--;
                break;
            case '+':
                data[data_ptr]++;
                break;
            case '-':
                data[data_ptr]--;
                break;
            case '.':
                std::putchar(data[data_ptr]);
                break;
            case ',':
                data[data_ptr] = std::getchar();
                break;
            case '[':
                if (data[data_ptr] == 0) {
                    code_ptr = jump_table[code_ptr];
                }
                break;
            case ']':
                if (data[data_ptr] != 0) {
                    code_ptr = jump_table[code_ptr];
                }
                break;
        }
        code_ptr++;        
    }
};


int main() {
    interpreter i(30000);
    i.load_code("[caca][cac[ccc]]");
    return 0;
}