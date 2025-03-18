#include <cstdio>
#include <vector>
#include <unordered_map>
#include <cstdint>

using uint64 = std::uint64_t;
template<typename T>
using vec = std::vector<T>;
template<typename KT, typename VT>
using dict = std::unordered_map<KT, VT>;

// abstract class for input and output
struct char_input {
    virtual char get() = 0;
};
struct char_output {
    virtual void put(char c) = 0;
};

// concrete classes for input and output implementing stdin and stdout
struct char_input_stdin : char_input {
    char get() override {
        return std::getchar();
    }
};

struct char_output_stdout : char_output {
    void put(char c) override {
        std::putchar(c);
    }
};

// concrete class for input implementing a string
struct char_input_string : char_input {
    uint64 ptr;
    const char* string;
    uint64 length;
    char boundary;
    char_input_string(const char* string, char boundary = '\0'): 
        ptr(0),
        string(string),
        length(std::strlen(string)),
        boundary(boundary)
    {}
    char get() override {
        if (ptr >= length) {
            return boundary;
        }
        return string[ptr++];
    }
};

struct interpreter {
    uint64 data_ptr;
    uint64 code_ptr;
    vec<char> data;
    const vec<char>& code;

    char_input* input;
    char_output* output;

    dict<uint64, uint64> jump_table;
    interpreter(uint64 data_length, const vec<char>& code, char_input* input = new char_input_stdin(), char_output* output = new char_output_stdout()):
        data_ptr(0),
        code_ptr(0),
        data(data_length, 0),
        code(code),
        input(input),
        output(output),
        jump_table()
    {
        vec<uint64> bracket_index_stack;
        for (uint64 i=0; i<code.size(); i++) {
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
        for (uint64 i=0; i<data.size(); i++) {
            std::printf("%c", data[i]);
        }
        std::printf("\n");
    }

    bool step() {
        if (code_ptr >= code.size()) {
            return false; // halt
        }
        switch (code[code_ptr]) {
            case '>':
                data_ptr = (data_ptr + 1) % data.size();
                break;
            case '<':
                data_ptr = (data_ptr + data.size() - 1) % data.size();
                break;
            case '+':
                data[data_ptr]++;
                break;
            case '-':
                data[data_ptr]--;
                break;
            case '.':
                output->put(data[data_ptr]);
                break;
            case ',':
                data[data_ptr] = input->get();
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

vec<char> read_code_from_file(const char* filename) {
    FILE* code_file = std::fopen(filename, "r");
    if (code_file == nullptr) {
        std::printf("Error: cannot open file %s\n", filename);
        return vec<char>();
    }
    std::fseek(code_file, 0, SEEK_END);
    uint64 code_length = std::ftell(code_file);
    std::rewind(code_file);
    vec<char> code(code_length);
    std::fread(code.data(), 1, code_length, code_file);
    std::fclose(code_file);
    return code;
}

std::tuple<uint64, vec<char>> read_args(int argc, char** argv) {
    if (argc < 2) {
        std::printf("Usage: %s <data_length> <code_filename>\n", argv[0]);
        return {0, vec<char>()};
    }
    uint64 data_length = std::stoull(argv[1]);
    vec<char> code = read_code_from_file(argv[2]);
    return {data_length, code};
}

int main(int argc, char** argv) {
    if (argc < 2) {
        std::printf("Usage: %s <data_length> <code_filename>\n", argv[0]);
        return 1;
    }
    auto [data_length, code] = read_args(argc, argv);

    interpreter i(data_length, code);
    while (i.step()) {
        // i.print_data();
    }
    return 0;
}