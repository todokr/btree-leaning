package main

import (
	"bufio"
	"os"

	"b-young-plant/btree"
	"b-young-plant/cli"
)

func main() {
	tree := btree.NewBtree()
	scanner := bufio.NewScanner(os.Stdin)
	demo := cli.NewCLI(scanner, tree)
	demo.Start()
}
