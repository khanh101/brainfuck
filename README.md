# BRAINFUCK INTERPRETER

## EXTENSION

- if a number in decimal `n` follows `<>+-`, the operation will be repeat `n` times. for example, `+48` means add `48` into the current cell

- new commands: `asmdr` add/subtract/multiply/divide/remainder current cell and next cell then put the output to the current cell

- new commands: `z` set the current cell to zero

- new commands: `w` swap the current cell and the next cell

- new commands: `_` noop

## PROGRAM

- `hello.bf`: just hello word

- `add.bf`: single digit adding `23 -> (2 + 3) -> 5`, the time complexity is of `O(2^n)`

# UNIVERSAL SEARCH

use `brainfuck` to search for algorithms

- `./bin/universal_search_factorize 323` example for factorizing `221 = 13 x 17`, (`17 x 19` takes forever)