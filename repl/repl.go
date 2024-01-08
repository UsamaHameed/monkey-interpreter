package repl

import (
	"bufio"
	"fmt"
	"io"
	"github.com/UsamaHameed/monkey-interpreter/lexer"
	"github.com/UsamaHameed/monkey-interpreter/token"
)

func Start(in io.Reader, out io.Writer) {

    scanner := bufio.NewScanner(in)

    for {

        fmt.Fprintf(out, ">> ")
        scanned := scanner.Scan()

        if !scanned {
            return
        }

        line := scanner.Text()
        l := lexer.New(line)

        for t := l.NextToken(); t.Type != token.EOF; t = l.NextToken() {
            fmt.Println(t.Literal, t.Type)
            fmt.Fprintf(out, "%+v\n", t)
        }
    }
}
