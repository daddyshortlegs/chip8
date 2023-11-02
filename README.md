# Chip 8

## Intro

This is my implementation of the Chip 8 "Emulator" in Go. Chip 8 was a VM that was written for the
Cosmac VIP computer. There's plenty of good articles that explain what it is and how to write one,
such as this one:

https://tobiasvl.github.io/blog/write-a-chip-8-emulator/

This is not an emulator as such as it does not emulate hardware and is really an interpreter that
executes bytecode. This bytecode can be found in ROMs that you can download from all over the
internet. A simple one is included in this repository.

## Status

Most of the opcodes have been implemented however there are some issues with some ROMs. 

## Running

There is a Makefile included that you can use to build and run. It has been built and run on a
Mac and has not been tested on anything else! To compile the project run:

`make`

To run with the included ROM, run:

`make run`

If you want to run a different ROM, you do:

`./chip8-app -rom "test_opcode.ch8"`


## The code

The bulk of the code is in the `chip8` directory and package. This contains the core logic. You
can run the tests from within this directory by running `make test`. The `main`program and anything 
to do with the display can be found in the main package. I use SDL to interact with the display.
