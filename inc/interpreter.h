#ifndef __INTERPRETER__
#define __INTERPRETER__
#include "type.h"
#include "input_output.h"

struct token {
    char operation;  // the operation character ('+', '-', etc.)
    int count;      // how many times to perform it (written in base 10) (default 1)
};

vec<token> parse_code(const vec<char>& source_code) {
    set<char> command_set = {'[', ']', '<', '>', '+', '-', '.', ','};
    set<char> multi_command_set = {'<', '>', '+', '-', '.', ','};
    vec<token> token_list;
    for (uint64 i = 0; i < source_code.size(); ++i) {
        char c = source_code[i];
        if (c == '#') { // skip comment
            uint64 j = i+1;
            while (j < source_code.size()) {
                if (source_code[j] == '\n') {
                    break;
                }
                ++j;
            }
            i = j;
        }
        if (command_set.count(c) == 0) {
            continue; // ignore other characters 
        }

        token t;

        if (multi_command_set.count(c)) {
            int count = 0;
            // Look ahead for digits.
            uint64 j = i + 1;
            while (j < source_code.size() and std::isdigit(source_code[j])) {
                count = count * 10 + (source_code[j] - '0');
                ++j;
            }
            if (count == 0) { // If no digit is found, count remains 0. Use default 1.
                count = 1;
            }
            t = token{c, count};
            i = j - 1;  // Skip processed digits.
        }
        else if (command_set.count(c)) {
            // For the other commands, count is just 1.
            t = token{c, 1};
        }

        if (token_list.size() > 0 and token_list[token_list.size()-1].operation == t.operation) {
            token_list[token_list.size()-1].count += t.count; // update prev token
        } else {
            token_list.push_back(t);
        }
    }
    return token_list;
}

dict<uint64, uint64> calc_jump_table(const vec<token>& code) {
    dict<uint64, uint64> jump_table;
    vec<uint64> bracket_index_stack;
    for (uint64 i=0; i<code.size(); i++) {
        switch (code[i].operation) {
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
    return jump_table;
}

struct interpreter {
    uint64 data_ptr;
    uint64 code_ptr;
    vec<char> data;

    char_input* input;
    char_output* output;

    vec<token> code;
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
        this->code = parse_code(source_code);
        this->jump_table = calc_jump_table(this->code);
    }

    void print_code() {
        std::printf("code: ");
        for (uint64 i=0; i<code.size(); i++) {
            token t = code[i];
            if (t.count == 1) {
                std::printf("%c", t.operation);
            } else {
                std::printf("%c%d", t.operation, t.count);
            }
        }
        std::printf("\n");
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
            std::printf("|%d", data[i]);
        }
        std::printf("|\n");
    }

    bool step() {
        if (code_ptr >= code.size()) {
            return false; // halt
        }
        token t = code[code_ptr];
        switch (t.operation) {
            case '>':
                data_ptr = (data_ptr + t.count) % data.size();
                break;
            case '<':
                data_ptr = (data_ptr + data.size() - t.count) % data.size();
                break;
            case '+':
                data[data_ptr] = data[data_ptr] + t.count;
                break;
            case '-':
                data[data_ptr] = data[data_ptr] - t.count;
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