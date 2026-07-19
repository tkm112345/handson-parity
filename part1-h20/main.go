package main

import "fmt"

func countOnes(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1) // 最下位ビットが1なら足す
		b >>= 1 // 1ビット右にずらして次を見る
	}
	return count
}

// parity は偶数パリティビット ( 0 or 1 )を返す
func parity(code byte) byte {
	return byte(countOnes(code) % 2) // 1の個数が奇数なら1, 偶数なら0
}

// encode は7bitコードを左に一つずらし，最下位にパリティを入れて8bitにする
func encode(code byte) byte {
	return (code << 1) | parity(code)
}

// check は8bitデータの1の個数が偶数ならtrue
func check(data byte) bool {
	return countOnes(data)%2 ==0
}

// flip はpos 番目(最下位が0)のビットを1つ反転させる( = 通信路の誤りを模倣)
func flip(data byte, pos int) byte {
	return data ^ (1 << pos) // XORで1ビットだけ反転　// XORの性質として，反転させたいbitだけ相手を1にしたらいい
}

func main() {
	// var code byte = 0b1011001 // 7bitコード

	// fmt.Printf("code     = %08b\n", code)
	// fmt.Printf("1 の個数 = %d\n", countOnes(code))

	// codes := []byte{0b1011001, 0b1011000, 0b0000000, 0b1111111}

	// fmt.Println("code(7bit) 1の数 parity encode(8bit)")
	// for _, c := range codes {
	// 		fmt.Printf("%07b  %d      %d     %08b\n", c, countOnes(c), parity(c), encode(c))
	// }

	// var data byte = 0b10110010

	// fmt.Printf("data     = %08b\n", data)
	// fmt.Printf("1 の個数 = %d\n", countOnes(data))
	// fmt.Printf("check    = %t   (偶数なら true = 正常)\n", check(data))

	var sent byte = 0b10110010 // 正しく送ったデータ（1 の個数 = 偶数）
	recv1 := flip(sent,3) // 途中で 3 番目のビット// 途中で 3 番目のビットが化けたが化けた
	recv2 := flip(flip(sent,3), 5) // 2bit誤り(3と5番め)

	// fmt.Printf("sent = %08b  check=%t\n", sent, check(sent))
	// fmt.Printf("recv = %08b  check=%t\n", recv, check(recv))

	fmt.Printf("sent       = %08b  check=%t\n", sent, check(sent))
	fmt.Printf("1bit error = %08b  check=%t  → 検出できた\n", recv1, check(recv1))
	fmt.Printf("2bit error = %08b  check=%t  → 見逃した！\n", recv2, check(recv2))


  }
