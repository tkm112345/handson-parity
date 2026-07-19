package main

import "fmt"

func dataBits(b byte, n int) []int {
	bits := make([]int, n)
	for i := 0 ; i < n; i++ {
		bits[i] = int((b >> i) & 1) // i番目の下位bitを取り出す //uartが下位ビットから電線に送るため
	}
	return bits
}

// countOnes はビット列中の 1 の個数を返す
func countOnes(bits []int) int {
	c := 0
	for _, b := range bits {
			c += b
	}
	return c
}

// parityBit はデータから偶数(E)/奇数(O)パリティを計算する（N=なしは -1）
func parityBit(data []int, mode string) int {
	ones := countOnes(data)
	switch mode {
	case "E":
			return ones % 2   // 偶数: 1 の数を偶数に揃える
	case "O":
			return 1 - ones%2 // 奇数: 1 の数を奇数に揃える
	default:
			return -1         // なし
	}
}

// buildFrame は 1 バイトを UART フレーム(start+data+parity+stop)に組み立てる
func buildFrame(b byte, n int, mode string) []int {
	frame := []int{0}
	data := dataBits(b,n)
	frame = append(frame,data...)
	if mode !="N" {
		frame = append(frame,parityBit(data, mode)) // パリティ
	}
	frame = append(frame, 1) // ストップビット
	return frame
}

// showFrame はフレームをビット列として表示する
func showFrame(frame []int) {
	for i, b := range frame {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(b)
	}
	fmt.Printf("   (%dbit)\n", len(frame))
}

// receive はフレームを分解し、復元バイトとパリティ検査結果(true=正常)を返す
func receive(frame []int, n int, mode string)(byte, bool) {
	data := frame[1:1+n] // スタートの次からn個がデータ
	ok := true
	if mode != "N" {
		p := frame[1+n]
		ok = parityBit(data, mode) == p // 再計算と一致するか
	}
	var b byte
	for i := 0; i< n; i++ {
		b |= byte(data[i] << i) // LSB first でバイトに戻す
	}
	return b, ok
}

