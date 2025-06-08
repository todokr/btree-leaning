package cli

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"strconv"
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
	c.printPrompt()
	for {
		if c.scanner.Scan() {
			c.processInput(c.scanner.Text())
		}
	}
}

func (c *CLI) printHelp() {
	fmt.Println(`Usage:
set <key>=<val>
del <key>
get <key>
rand <n>`)
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
		c.printHelp()
	case "set":
		args := fields[1:]
		if len(args) != 1 {
			fmt.Println("Usage: set <key>=<value>")
			return
		}
		c.processSetCommand(args[0])

	case "del":
		args := fields[1:]
		if len(args) != 1 {
			fmt.Println("Usage: del <key>")
			return
		}
		c.processDeleteCommand(args[0])
	case "get":
		args := fields[1:]
		if len(args) != 1 {
			fmt.Println("Usage: get <key>")
			return
		}
		c.processGetCommand(fields[0])
	case "rand":
		args := fields[1:]
		if len(args) != 1 {
			fmt.Println("Usage: get <key>")
			return
		}
		c.processRandomCommand(args[0])
	}
	c.printPrompt()
}

func (c *CLI) processSetCommand(arg string) {
	pair := strings.Split(arg, "=")
	c.tree.Insert([]byte(pair[0]), []byte(pair[1]))
	fmt.Println(c.tree)
}

func (c *CLI) processDeleteCommand(key string) {
	res := c.tree.Delete([]byte(key))
	if !res {
		fmt.Println("Key not found.")
		return
	}
	fmt.Println(c.tree)
}

func (c *CLI) processGetCommand(key string) {
	val, err := c.tree.Find([]byte(key))

	if err != nil {
		fmt.Println("Key not found.")
		return
	}
	fmt.Println(string(val))
}

func (c *CLI) processRandomCommand(n string) {
	num, err := strconv.Atoi(n)
	if err != nil {
		fmt.Println("not a number")
		return
	}
	for range num {
		s, err := genRandomStr(10)
		if err != nil {
			fmt.Println("error generating random key")
			return
		}
		c.tree.Insert([]byte(s), []byte(s))
	}
	fmt.Println(c.tree)
}

func genRandomStr(digit uint32) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"

	b := make([]byte, digit)
	if _, err := rand.Read(b); err != nil {
		return "", errors.New("unexpected error...")
	}

	var result string
	for _, v := range b {
		// index が letters の長さに収まるように調整
		result += string(letters[int(v)%len(letters)])
	}
	return result, nil
}
