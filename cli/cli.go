package cli

import (
	"bufio"
	"fmt"
	"strings"

	"b-young-plant/btree"
)

type CLI struct {
	scanner *bufio.Scanner
	tree    *btree.Btree
}

func NewCLI(s *bufio.Scanner, b *btree.Btree) *CLI {
	return &CLI{s, b}
}

func (c *CLI) Start() {
	c.printHelp()
	c.printPrompt()
	for {
		if c.scanner.Scan() {
			c.processInput(c.scanner.Text())
		}
	}
}

func (c *CLI) printHelp() {
	fmt.Println(`
set <key>=<val>
del <key>
get <key>`)
}

func (c *CLI) printPrompt() {
	fmt.Print("> ")
}

func (c *CLI) processInput(line string) {
	fields := strings.Fields(line)

	if len(fields) < 1 {
		return
	}
	command := strings.ToLower(fields[0])

	switch command {
	default:
		fmt.Printf("Unknown command \"%s\"\n", command)
	case "set":
		c.processSetCommand(fields[1:])
	case "del":
		c.processDeleteCommand(fields[1:])
	case "get":
		c.processGetCommand(fields[1:])
	case "auto":
		c.processSetCommand([]string{"a=a"})
		c.processSetCommand([]string{"aa=aa"})
		c.processSetCommand([]string{"aaa=aaa"})
		c.processSetCommand([]string{"b=b"})
	}
	c.printPrompt()
}

func (c *CLI) processSetCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: set <key>=<value>")
		return
	}
	pair := strings.Split(args[0], "=")
	c.tree.Insert([]byte(pair[0]), []byte(pair[1]))
	fmt.Println(c.tree)
}

func (c *CLI) processDeleteCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: del <key>")
		return
	}
	res := c.tree.Delete([]byte(args[0]))

	if !res {
		fmt.Println("Key not found.")
		return
	}
	fmt.Println(c.tree)
}

func (c *CLI) processGetCommand(args []string) {
	if len(args) != 1 {
		fmt.Println("Usage: get <key>")
		return
	}
	val, err := c.tree.Find([]byte(args[0]))

	if err != nil {
		fmt.Println("Key not found.")
		return
	}
	fmt.Println(string(val))
}
