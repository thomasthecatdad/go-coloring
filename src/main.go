package main

import (
	g "./graphs"
	t "./testHarness"
)

func main() {
	gr := t.ParseFile("res/Sample01.txt")
	//gr := t.ParseFile("res/Error01.txt")
	//gr := t.ParseFile("res/Error02.txt")
	g.PrintGraph(&gr)
}
