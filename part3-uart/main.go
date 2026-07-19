package main

import "fmt"

func main() {
	// 'A' = 0x41 = 0b0100_0001
	// fmt.Println("'A' の 8 データビット (LSB first):", dataBits('A', 8))
	
	// fmt.Print("8E1  'A': ")
	// showFrame(buildFrame('A', 8, "E"))
	// fmt.Print("8O1  'A': ")
	// showFrame(buildFrame('A', 8, "O"))
	// fmt.Print("8N1  'A': ")
	// showFrame(buildFrame('A', 8, "N"))
	
	frame := buildFrame('A', 8, "E")
	b, ok := receive(frame, 8, "E")
	fmt.Print("受信フレーム: ")
	showFrame(frame)
	fmt.Printf("復元バイト = %q, パリティ検査 = %t\n", b, ok)


}