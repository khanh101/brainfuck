.PHONY: interpreter

COMMON_FLAGS = -O4 -std=c++23 -march=native -fno-math-errno -I inc -I /Users/khanh/miniforge3/envs/brainfuck/include -L /Users/khanh/miniforge3/envs/brainfuck/lib

interpreter:
	clang++ $(COMMON_FLAGS) -o interpreter.out interpreter.cpp
	DYLD_LIBRARY_PATH=/Users/khanh/miniforge3/envs/brainfuck/lib ./interpreter.out 50 code/hello.bf

universal_search:
	clang++ $(COMMON_FLAGS)  -lgmp -lgmpxx -o universal_search.out universal_search.cpp
	DYLD_LIBRARY_PATH=/Users/khanh/miniforge3/envs/brainfuck/lib ./universal_search.out