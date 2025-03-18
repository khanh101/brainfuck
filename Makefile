.PHONY: interpreter

interpreter:
	clang -O4 -lstdc++ -std=c++23 -march=native -fno-math-errno -Iinc -o interpreter.out interpreter.cpp
	./interpreter.out 50 code/hello.bf
