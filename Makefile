.PHONY: run

run:
	clang -Og -lstdc++ -std=c++23 -march=native -fno-math-errno -o main.out main.cpp
	./main.out 50 code/hello.bf
