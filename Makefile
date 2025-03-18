.PHONY: run

run:
	clang -lstdc++ -std=c++23 -O4 -march=native -fno-math-errno -o main.out main.cpp
	./main.out
