// gsize tells you how large standard input will be after gzip compression
package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/tabwriter"
)

type countWriter int64

func (size *countWriter) Write(p []byte) (n int, err error) {
	n, err = ioutil.Discard.Write(p)
	*size += countWriter(n)
	return
}

func humanize(size int64) string {
	const (
		_        = iota
		kilobyte = 1 << (10 * iota)
		megabyte
		gigabyte
		terabyte
	)

	format := "%.f"
	value := float32(size)

	switch {
	case size >= terabyte:
		format = "%3.1f TB"
		value = value / terabyte
	case size >= gigabyte:
		format = "%3.1f GB"
		value = value / gigabyte
	case size >= megabyte:
		format = "%3.1f MB"
		value = value / megabyte
	case size >= kilobyte:
		format = "%3.1f KB"
		value = value / kilobyte
	}
	return fmt.Sprintf(format, value)
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(2)
	}
}

func main() {
	var compressedSize countWriter
	gw := gzip.NewWriter(&compressedSize)
	originalSize, err := io.Copy(gw, os.Stdin)
	die(err)
	die(gw.Flush())
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "  Original\t%s\t%d\t\n",
		humanize(originalSize), originalSize)
	fmt.Fprintf(tw, "Compressed\t%s\t%d\t\n",
		humanize(int64(compressedSize)), compressedSize)
	fmt.Fprintf(tw, "     Ratio\t%01.2f\t\n",
		float64(compressedSize)/float64(originalSize))
	tw.Flush()
}
