package main 

import "fmt"

var data = [4][4]int {
	{1, 0, 1, 1},
	{0, 1, 0, 1},
	{1, 1, 0, 0},
	{0, 0, 1, 1},
}

// encodeは4x4データに行・列パリティをつけて，5x5にする
func encode(data [4][4]int) [5][5]int {
	var g[5][5]int
	//元データをコピー
	for i :=0 ; i< 4; i++ {
		for j :=0; j<4; j++ {
			g[i][j] = data[i][j]
		}
	}

	// 行パリティ：各行(0..3)の1の個数の遇機を右端(列 4)に入れる
	for i := 0; i < 4; i++ {
		sum :=0
		for j:= 0; j<4;j++{
			sum += data[i][j]
		}
		g[i][4] = sum % 2
	}

	// 列パリティ：各列(0..4) の1の個数の遇機を下端(行 4)に入れる
	// 列4も含めるので，右下g[4][4]は全体のパリティになる
	for j := 0; j< 5; j++ {
		sum := 0
		for i:=0; i<4; i++ {
			sum += g[i][j]
		}
		g[4][j] = sum % 2
	}
	return g
}

// checkRows は各行の1の個数が奇数(=異常)ならtrueを返す
func checkRows(g [5][5]int) [5]bool {
	var bad [5]bool
	for i := 0; i< 5; i++ {
		sum := 0
		for j := 0; j< 5; j++ {
			sum += g[i][j]
		}
		bad[i] = sum%2 != 0
	}
	return bad
}

// checkCols は各列の1の個数が奇数(=異常)ならtrueを返す
func checkCols(g [5][5]int) [5]bool {
	var bad [5]bool
	for j := 0; j< 5; j++ {
		sum := 0
		for i := 0; i < 5; i++ {
			sum += g[i][j]
		}
		bad[j] = sum%2 != 0
	}
	return bad
}

// pritnGrid5は5x5のbit格子を表示
func printGrid5(g [5][5]int) {
	for i := 0; i< 5; i++ {
		for j := 0; j < 5; j ++ {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(g[i][j])
		}
		fmt.Println()
	}
}


// pritnGridは4x4のbit格子を表示
func printGrid(g [4][4]int) {
	for i := 0; i< 4; i++ {
		for j := 0; j < 4; j ++ {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(g[i][j])
		}
		fmt.Println()
	}
}

func main() {
	g := encode(data)
	printGrid5(g)
	fmt.Println("行ごとの異常(true=奇数):", checkRows(g))
	fmt.Println("列ごとの異常(true=奇数):", checkCols(g))

	fmt.Println("=== 誤りを1ビット注入: (行1,列2) を反転 ===")
	g[1][2] ^= 1 // XOR で 1 ビット反転（第1部の flip と同じ考え方）
	printGrid5(g)

	rows := checkRows(g)
	cols := checkCols(g)
	fmt.Println("異常な行:", rows)
	fmt.Println("異常な列:", cols)

	// 異常な行・列を1つずつ探す
	r,c := -1,-1
	for i := 0; i < 5; i++ {
		if rows[i] {
			r = i
		}
	}
	for j := 0; j<5;j++ {
		if cols[j] {
			c = j
		}
	}
	fmt.Printf("→ 交点 (%d,%d) が誤り。ここを反転して訂正\n", r, c)
	g[r][c] ^= 1 // 訂正

	fmt.Println("訂正後 異常な行:", checkRows(g))
	fmt.Println("訂正後 異常な列:", checkCols(g))


	fmt.Println("=== 誤りを2ビット注入: (1,2) と (2,3) を反転 ===")
	g[1][2] ^= 1
	g[2][3] ^= 1
	printGrid5(g)
	fmt.Println("異常な行:", checkRows(g))
	fmt.Println("異常な列:", checkCols(g))

}
