#ifndef __UNIVERSAL_SEARCH__
#define __UNIVERSAL_SEARCH__
#include <functional>
#include <gmpxx.h>

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

void universal_search(char* input_string, const std::function<bool(char*)>& test) {

}





#endif // __UNIVERSAL_SEARCH__