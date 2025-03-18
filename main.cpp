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
    uint64 data_length;
    interpreter(uint64 data_length): data_ptr(0), code_ptr(0), data(data_length, 0) {}

    const char* code;
    uint64 code_length;
    dict<uint64, uint64> jump_table;
    void load_code(const char* code) {
        this->code = code;
        this->code_length = std::strlen(code);
        this->jump_table = dict<uint64, uint64>();
        vec<uint64> bracket_index_stack;
        for (uint64 i=0; i<code_length; i++) {
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
    }

    void print_jump_table() {
        std::printf("jump table:\n");
        for (auto& [k, v] : jump_table) {
            if (k < v) {
                std::printf("\t%llu <-> %llu\n", k, v);
            }
        }
    }

    void print_data() {
        for (uint64 i=0; i<data_length; i++) {
            std::printf("%c", data[i]);
        }
        std::printf("\n");
    }

    bool step() {
        if (code_ptr >= code_length) {
            return false; // halt
        }
        switch (code[code_ptr]) {
            case '>':
                data_ptr = (data_ptr + 1) % data_length;
                break;
            case '<':
                data_ptr = (data_ptr - 1 + data_length) % data_length;
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
        return true;        
    }
};


int main() {
    interpreter i(50);
    i.load_code("++++++++[>++++[>++>+++>+++>+<<<<-]>+>+>->>+[<]<-]>>.>---.+++++++..+++.>>.<-.<.+++.------.--------.>>+.>++."); // hello world
    while (i.step()) {
        // i.print_data();
    }
    return 0;
}