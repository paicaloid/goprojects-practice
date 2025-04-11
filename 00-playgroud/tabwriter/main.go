package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// func main() {

// 	// Observe that the third line has no trailing tab,
// 	// so its final cell is not part of an aligned column.
// 	const padding = 3
// 	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', tabwriter.AlignRight|tabwriter.Debug)
// 	fmt.Fprintln(w, "a\tb\taligned\t")
// 	fmt.Fprintln(w, "aa\tbb\taligned\t")
// 	fmt.Fprintln(w, "aaa\tbbb\tunaligned\t") // no trailing tab
// 	fmt.Fprintln(w, "aaaa\tbbbb\taligned\t")
// 	w.Flush()

// }

func main() {
	w := new(tabwriter.Writer)

	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "aaaaaa\tb\tc\td\t")
	fmt.Fprintln(w, "123\t12345\t1234567\t123456789\t")
	fmt.Fprintln(w)
	w.Flush()
}
