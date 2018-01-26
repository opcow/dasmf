// Fixes segment headers for Atari 8-bit in dasm created binaries
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	var nextSeg, segLen int
	var ifName, ofName string

	if len(os.Args) < 2 {
		fmt.Println("Usage: dasmf infile [outfile]")
		os.Exit(1)
	}
	ifName = os.Args[1]
	if len(os.Args) > 2 {
		ofName = os.Args[2]
	} else {
		ofName = os.Args[1]
	}

	finfo, err := os.Stat(ifName)
	if err != nil {
		fmt.Println("dasmf: couldn't stat input file")
		os.Exit(1)
	}
	fSize := finfo.Size()
	f, err := os.Open(ifName)
	if err != nil {
		panic(err)
	}
	if fSize > 16384 {
		fmt.Println("dasmf: file too large.")
	}
	data := make([]byte, fSize)
	{
		defer f.Close()
		bufr := bufio.NewReader(f)
		_, err := bufr.Read(data)
		if err != nil {
			panic(err)
		}
		if fSize < 5 || fSize > 0xffc || (data[0] == 0xff && data[1] == 0xff) {
			fmt.Println("dasmf: input file is the wrong format")
			os.Exit(1)
		}
	}
	of, err := os.Create(ofName)
	if err != nil {
		panic(err)
	}
	defer of.Close()
	header := [2]byte{0xff, 0xff}
	of.Write(header[:])
	for i := 1; nextSeg <= int(fSize-2); i++ {
		fmt.Printf("\nProcessing segment %d\n", i)
		segLen = fixSegment(data[nextSeg:])
		of.Write(data[nextSeg : nextSeg+segLen+4])
		nextSeg += segLen + 4
	}
}

func fixSegment(segment []byte) int {
	var segAddr, segEnd, segLen int

	segAddr = int(segment[0]) | (int(segment[1]) << 8)
	fmt.Printf("Segment load address: 0x%04X\n", segAddr)
	segLen = (int(segment[2]) | (int(segment[3]) << 8))
	segEnd = segLen + (int(segment[0]) | (int(segment[1]) << 8)) - 1
	fmt.Printf("Segment end: 0x%04X\n", segEnd)
	fmt.Printf("Segment length: 0x%04X\n", segLen)

	segment[2] = byte(segEnd & 0xff)
	segment[3] = byte(segEnd >> 8)

	return segLen
}
