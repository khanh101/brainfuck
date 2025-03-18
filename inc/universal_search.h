#ifndef __UNIVERSAL_SEARCH__
#define __UNIVERSAL_SEARCH__
#include <gmpxx.h>
#include "type.h"
#include "interpreter.h"

using Z = mpz_class;

Z pow(Z a, Z n) {
    if (n < 0) {
        return pow(a, -n);
    }
    if (n == 0) {
        return 1;
    }

    Z half = pow(a, n / 2);
    if (n % 2 == 0) {
        return half * half;
    } else {
        return a * half * half;
    }
}

vec<char> get_code_from_z(Z i) {
    dict<uint64, char> ztoc = {
        {0, '_'},
        {1, '['},
        {2, ']'},
        {3, '<'},
        {4, '>'},
        {5, '+'},
        {6, '-'},
        {7, '.'},
        {8, ','},
        {9, 'a'},
        {10, 's'},
        {11, 'm'},
        {12, 'd'},
        {13, 'r'},
        {14, 'z'},
        {15, 'w'}
    };

    vec<char> source_code;
    while (i > 0) {
        Z mod = i % 16;
        auto it = ztoc.find(mod.get_ui());
        char c = it->second;
        source_code.push_back(c);
        i = i / 16;
    }
    return source_code;
}

tup<Z, vec<char>> universal_search(const vec<char>& input_string, const func<bool(const vec<char>&)>& test) {
    vec<tup<Z, interpreter*>> space;

    Z i = 0;

    while (true) {
        // make new program
        interpreter* ip = new interpreter(300, get_code_from_z(i), new char_input_string(input_string), new char_output_string());
        space.push_back({i, ip});
        std::printf("running %zu programs ...\n", space.size());
        // exec
        for (uint64 k=0; k < space.size(); k++) {
            auto [j, jp] = space[k];
            
            Z num_steps = pow(2, i - j); // run 2^{i - j} steps
            while (num_steps > 0) {
                bool keep_running = jp->step();
                if (not keep_running) { // halt
                    if (test(((char_output_string*)(jp->output))->buffer)) {
                        vec<char> buffer(((char_output_string*)(jp->output))->buffer); // copy
                        for (uint64 kk=0; kk <space.size(); kk++) {
                            delete std::get<1>(space[kk]);
                        }
                        return {j, buffer};
                    }
                    // otherwise, remove from running list
                    if (space.size() <= 1) {
                        space = vec<tup<Z, interpreter*>>();
                    } else {
                        delete jp;
                        space[k] = space[space.size()-1];
                        space.pop_back();
                    }
                    break;
                }
                num_steps -= 1;
            }
        }
        i++;
    }
}





#endif // __UNIVERSAL_SEARCH__