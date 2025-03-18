.PHONY: run

run:
	clang -O4 -lstdc++ -std=c++23 -march=native -fno-math-errno -Iinc -o main.out main.cpp
	./main.out 50 code/hello.bf
