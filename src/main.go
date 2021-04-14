package main

import (
	t "./testHarness"
)

func main() {
	g := t.ParseFile("res/SampleGraph01.txt")
	t.PrintGraph(&g)
}
