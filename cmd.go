package hobbit

import (
	"bufio"
	"os"
	"strings"
)

type Cmd struct {
	in  *bufio.Reader
	out *bufio.Writer
	//fileout *bufio.Writer
}

func NewCmd() *Cmd {
	cmd := new(Cmd)
	cmd.in = bufio.NewReader(os.Stdin)
	cmd.out = bufio.NewWriter(os.Stdout)
	return cmd
}

func (cmd *Cmd) ReadCommand() (string, error) {
	cmd.Prompt(">")
	text, err := cmd.in.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(text), nil
}

// Prompt a string to
func (cmd *Cmd) Prompt(s string) {
	cmd.out.WriteString(s)
	cmd.out.Flush()
}

func (cmd *Cmd) Promptln(s string) {
	cmd.out.WriteString(s)
	cmd.out.WriteString("\n")
	cmd.out.Flush()
}

func (cmd *Cmd) Banner(s string, filler byte) {
	s = FixToLength(s, 80, filler)
	cmd.out.WriteString(s)
	cmd.out.WriteString("\n")
	cmd.out.Flush()
}
