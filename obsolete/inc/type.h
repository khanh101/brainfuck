#ifndef __TYPE__
#define __TYPE__
#include <cstdint>
#include <vector>
#include <unordered_map>
#include <unordered_set>
#include <tuple>
#include <functional>

using uint64 = std::uint64_t;
template<typename T>
using vec = std::vector<T>;
template<typename KT, typename VT>
using dict = std::unordered_map<KT, VT>;
template<typename T>
using set = std::unordered_set<T>;
template<typename T1, typename T2>
using tup = std::tuple<T1, T2>;
template<typename T>
using func = std::function<T>;

#endif // __TYPE__