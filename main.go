package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	var nextSeg = 0
	var segLen = 0

	if len(os.Args) < 3 {
		fmt.Println("Usage: dasmf infile outfile")
		os.Exit(1)
	}
	ifname := os.Args[1]
	ofname := os.Args[2]

	if _, err := os.Stat(ifname); os.IsNotExist(err) {
		fmt.Println("dasmf: input file not found")
		os.Exit(1)
	}
	f, err := os.Open(ifname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data := make([]byte, 8192)

	bufr := bufio.NewReader(f)
	br, err := bufr.Read(data)
	if br < 5 || br > 0xffc {
		fmt.Println("dasm: input file is the wrong format")
		os.Exit(1)
	}

	of, err := os.Create(ofname)
	if err != nil {
		panic(err)
	}
	defer of.Close()
	header := [2]byte{0xff, 0xff}
	of.Write(header[:])
	//	of.Write(data[:br])
	for i := 1; nextSeg <= br-2; i++ {
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
	fmt.Printf("Segment lenght: 0x%04X\n", segLen)

	segment[2] = byte(segEnd & 0xff)
	segment[3] = byte(segEnd >> 8)

	return segLen
}
