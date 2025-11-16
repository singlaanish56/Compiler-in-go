package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/singlaanish56/Compiler-in-go/compiler"
	"github.com/singlaanish56/Compiler-in-go/lexer"
	"github.com/singlaanish56/Compiler-in-go/object"
	"github.com/singlaanish56/Compiler-in-go/parser"
	"github.com/singlaanish56/Compiler-in-go/vm"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	constants := []object.Object{}
	globalStore := make([]object.Object, vm.GlobalSize)
	symbolTable := compiler.NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	for {
		fmt.Fprintf(out, PROMPT)
		scannedLine := scanner.Scan()
		if !scannedLine {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParserProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		compiler := compiler.NewWithState(symbolTable, constants)
		err := compiler.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops, Compiler failed:\n %s\n", err)
			continue
		}

		bytecode := compiler.Bytecode()
		constants = bytecode.Constants

		vmMachine := vm.NewWithGlobalStore(bytecode, globalStore)
		err = vmMachine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops, VM failed:\n %s\n", err)
			continue
		}

		stackTop := vmMachine.LastPoppedStackElement()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, parserErrors []error) {
	io.WriteString(out, "ran into these parser errors:\n")
	for _, err := range parserErrors {
		io.WriteString(out, "\t"+err.Error()+"\n")
	}
}
