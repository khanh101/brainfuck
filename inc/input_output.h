#ifndef __INPUT_OUTPUT__
#define __INPUT_OUTPUT__

#include "type.h"
#include <cstdio>
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

struct char_output_string: char_output {
    vec<char> buffer;
    void put(char c) override {
        buffer.push_back(c);
    }
    char* to_string() {
        char* s = new char[buffer.size() + 1];
        std::memcpy(s, buffer.data(), buffer.size());
        s[buffer.size()] = '\0';
        return s;
    }
};

#endif // __INPUT_OUTPUT__