package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Symbol struct {
	Name               string
	Type               string
	Line, StartPos     int
	EndPos             int
}

func (s Symbol) String() string {
	return fmt.Sprintf("%s %s %d:%d-%d",
		s.Name, s.Type, s.Line, s.StartPos, s.EndPos)
}

func makeSymbol(line string, lineNum int, name, typ string) Symbol {
	start := strings.Index(line, name) + 1
	end   := start + len(name)            
	return Symbol{Name: name, Type: typ, Line: lineNum, StartPos: start, EndPos: end}
}

var patterns = []struct {
	re  *regexp.Regexp
	typ string
}{
	{regexp.MustCompile(`^\s*import\s+"([^"]+)"`), "import"},
	{regexp.MustCompile(`^\s*service\s+(\w+)`),     "service"},
	{regexp.MustCompile(`^\s*rpc\s+(\w+)`),         "method"},
	{regexp.MustCompile(`^\s*enum\s+(\w+)`),        "enum"},
	{regexp.MustCompile(`^\s*message\s+(\w+)`),     "message"},
}

func parseProtoFile(filename string) ([]Symbol, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var symbols []Symbol
	sc := bufio.NewScanner(f)
	lineNum := 0

	for sc.Scan() {
		lineNum++
		line := sc.Text()

		if isInsideBlock(line) {
			continue
		}

		for _, p := range patterns {
			if m := p.re.FindStringSubmatch(line); m != nil {
				symbols = append(symbols, makeSymbol(line, lineNum, m[1], p.typ))
				break 
			}
		}
	}
	return symbols, sc.Err()
}

func isInsideBlock(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	leading := len(line) - len(strings.TrimLeft(line, " \t"))
	if leading > 0 &&
		(strings.HasPrefix(trimmed, "message ") ||
			strings.HasPrefix(trimmed, "enum ") ||
			strings.HasPrefix(trimmed, "oneof ")) {
		return true
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <proto-file>\n", os.Args[0])
		os.Exit(1)
	}
	syms, err := parseProtoFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	for _, s := range syms {
		fmt.Println(s)
	}
}
