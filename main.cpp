
#include "interpreter.h"



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