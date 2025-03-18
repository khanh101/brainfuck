#ifndef __INTERPRETER__
#define __INTERPRETER__
#include "type.h"
#include "io.h"

struct interpreter {
    uint64 data_ptr;
    uint64 code_ptr;
    vec<char> data;

    char_input* input;
    char_output* output;

    vec<char> code;
    dict<uint64, uint64> jump_table;
    interpreter(uint64 data_length, const vec<char>& source_code, char_input* input = new char_input_stdin(), char_output* output = new char_output_stdout()):
        data_ptr(0),
        code_ptr(0),
        data(data_length, 0),
        code(),
        input(input),
        output(output),
        jump_table()
    {
        // shorten code
        set<char> allowed_char_set = {'[', ']', '<', '>', '+', '-', '.', ','};
        for (uint64 i = 0; i < source_code.size(); i++) {
            if (allowed_char_set.find(source_code[i]) != allowed_char_set.end()) {
                code.push_back(source_code[i]);
            }
        }
        // calculate jump table
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

#endif // __INTERPRETER__