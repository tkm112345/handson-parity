# ハンズオン: パリティビットを Go で実験する

応用情報技術者試験の **パリティビット問題** を「腑に落とす」ための写経ハンズオン。
前半（Module 0〜5）は **H20 秋 問6**（1 次元パリティ・誤り「検出」）、
後半（Module 6〜10）は **R5 秋 問4**（二次元パリティ・誤り「訂正」）を扱う。

> 7 ビットのコードと 1 ビットのパリティビットからなる 8 ビットのデータで発生した
> 誤りに関する記述として、適切なものはどれか
>
> **答え: 1 ビットの誤りは検出できるが、その位置は分からない（訂正できない）。**
> **2 ビットの誤りは検出できない。**

この「なぜ？」を、小さな Go プログラムで 1 つずつ目で確認する。
XOR ハンズオン（`../xor`）の続き。パリティは「1 の個数の偶奇」であり、
その計算は XOR の交換法則・結合法則の応用そのもの。

## 進め方

- 各モジュールのコードを**写経（打ち込み）**する
- `go run .`（対象ディレクトリ）で実行し、「期待する出力」と一致するか確認する
- コードより **「なぜそうなるか」** に注目する（下の解説を読む）

### ディレクトリ構成

Go は「1 ディレクトリ ＝ 1 パッケージ、`func main()` は 1 つだけ」なので、
第1部と第2部を **別ディレクトリ** に分けて両方とも動かせる状態で残す。

```
learn/parity/
├── part1-h20/   # 第1部 H20秋問6（1次元・検出）Module 0〜5 の最終コード
│   └── main.go
├── part2-r5/    # 第2部 R5秋問4（二次元・訂正）Module 6〜10 を写経（各モジュールで上書き）
│   └── main.go
└── part3-uart/  # 第3部 シリアル通信(UART)のパリティを擬似検証 Module 11〜16
    ├── uart.go  # 共通部品（dataBits, parityBit, buildFrame, receive ...）を追記していく
    └── main.go  # func main() だけ。Module ごとに書き換える
```

実行は対象ディレクトリを指定する: `go run ./part1-h20` / `go run ./part2-r5` / `go run ./part3-uart`。
（第1部の Module 0〜5 の解説は当時 `main.go` を書き換えながら進めた記録。現在は上記の構成に整理済み）

パリティ方式は本ハンズオンでは **偶数パリティ**（1 の個数を偶数に揃える）で統一する。
第1部のパリティビットは **最下位ビット（一番右）**、第2部は **行の右端・列の下端** に付ける。

---

## Module 0: 環境準備

作業場所を作り、Go が動くことだけ確認する。

```bash
mkdir -p ~/learn/parity && cd ~/learn/parity
go mod init handson-parity   # go.mod がまだ無ければ
go version
```

`main.go` を作り、`fmt.Println` で 1 行出して `go run main.go` が動けば OK。

---

## Module 1: 1 の個数を数える

**ゴール**: 8 ビット中の「1 の個数」を数える関数を作る。パリティの土台。

```go
package main

import "fmt"

// countOnes は byte 内の 1 のビット数を数える
func countOnes(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1) // 最下位ビットが 1 なら足す
		b >>= 1             // 1 ビット右にずらして次を見る
	}
	return count
}

func main() {
	var code byte = 0b1011001 // 7 ビットのコード

	fmt.Printf("code     = %08b\n", code)
	fmt.Printf("1 の個数 = %d\n", countOnes(code))
}
```

**期待する出力**（`0b1011001` = 89 を `%08b` で 8 桁ゼロ埋め表示）:

```
code     = 01011001
1 の個数 = 4
```

**ポイント**:
- `b & 1` … 最下位ビットだけ取り出す（1 か 0）。
- `b >>= 1` … 右シフトで次のビットを最下位に持ってくる。
- これを 8 回まわせば全ビットの 1 を数えられる。**パリティ = この個数の偶奇**。

---

## Module 2: パリティビットを付ける（送信側）

**ゴール**: 7 ビットコードに偶数パリティを付けて 8 ビットにする。

```go
package main

import "fmt"

func countOnes(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

// parity は偶数パリティビット（0 or 1）を返す
func parity(code byte) byte {
	return byte(countOnes(code) % 2) // 1 の個数が奇数なら 1、偶数なら 0
}

// encode は 7 ビットコードを左に 1 つずらし、最下位にパリティを入れて 8 ビットにする
func encode(code byte) byte {
	return (code << 1) | parity(code)
}

func main() {
	codes := []byte{0b1011001, 0b1011000, 0b0000000, 0b1111111}

	fmt.Println("code(7bit) 1の数 parity encode(8bit)")
	for _, c := range codes {
		fmt.Printf("%07b  %d      %d     %08b\n", c, countOnes(c), parity(c), encode(c))
	}
}
```

**期待する出力**:

```
code(7bit) 1の数 parity encode(8bit)
1011001  4      0     10110010
1011000  3      1     10110001
0000000  0      0     00000000
1111111  7      1     11111111
```

**ポイント**:
- `parity` … 1 の個数が奇数なら 1 を足して偶数にする。偶数ならそのまま 0。
- `encode` … `code << 1` で下位 1 ビットを空け、`| parity` でそこに入れる。
- 結果 `encode` の 8 ビットは、**必ず 1 の個数が偶数**になっている（ここが肝）。

---

## Module 3: 受信側で検査する

**ゴール**: 受け取った 8 ビットの「1 の個数が偶数か」だけで正常/異常を判定する。

```go
package main

import "fmt"

func countOnes(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

// check は 8 ビットデータの 1 の個数が偶数なら true（正常とみなす）
func check(data byte) bool {
	return countOnes(data)%2 == 0
}

func main() {
	var data byte = 0b10110010 // Module 2 で作った正しい 8 ビット

	fmt.Printf("data    = %08b\n", data)
	fmt.Printf("1 の個数 = %d\n", countOnes(data))
	fmt.Printf("check   = %t   (偶数なら true = 正常)\n", check(data))
}
```

**期待する出力**:

```
data    = 10110010
1 の個数 = 4
check   = true   (偶数なら true = 正常)
```

**ポイント**:
- 受信側は元の 7 ビットが何だったかを**知らない**。分かるのは「1 の個数の偶奇」だけ。
- 偶数パリティで送ったので、**誤りが無ければ必ず偶数** → `check == true`。
- 逆に言うと、受信側の判定材料は「偶奇」の 1 ビット分の情報しかない。
  この情報量の少なさが、次の Module 4・5 の限界を生む。

---

## Module 4: 1 ビットの誤り → 検出できる

**ゴール**: 1 ビットだけ反転させると 1 の個数が奇数になり、`check` が false になることを見る。

```go
package main

import "fmt"

func countOnes(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func check(data byte) bool {
	return countOnes(data)%2 == 0
}

// flip は pos 番目（最下位が 0）のビットを 1 つ反転させる（＝通信路の誤りを模擬）
func flip(data byte, pos int) byte {
	return data ^ (1 << pos) // XOR で 1 ビットだけ反転（xor ハンズオンの応用）
}

func main() {
	var sent byte = 0b10110010 // 正しく送ったデータ（1 の個数 = 偶数）
	recv := flip(sent, 3)      // 途中で 3 番目のビットが化けた

	fmt.Printf("sent = %08b  check=%t\n", sent, check(sent))
	fmt.Printf("recv = %08b  check=%t\n", recv, check(recv))
}
```

**期待する出力**:

```
sent = 10110010  check=true
recv = 10111010  check=false
```

**ポイント**:
- 1 ビット反転 → 1 の個数が 1 つ増減 → 偶奇が必ず反転 → `check=false` で**検出成功**。
- ただし `check` が教えてくれるのは「どこかが変」だけ。
  **どのビットが化けたかは分からない** → 訂正はできない。これが設問の答えの前半。
- `pos` を 0〜7 のどれに変えても、必ず `check=false` になる（試してみる）。

---

## Module 5: 2 ビットの誤り → 検出できない（設問の核心）

**ゴール**: 2 ビット反転させると 1 の個数が偶数に戻り、`check` が誤って true になる（見逃す）。

```go
package main

import "fmt"

func countOnes(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func check(data byte) bool {
	return countOnes(data)%2 == 0
}

func flip(data byte, pos int) byte {
	return data ^ (1 << pos)
}

func main() {
	var sent byte = 0b10110010

	recv1 := flip(sent, 3)          // 1 ビット誤り
	recv2 := flip(flip(sent, 3), 5) // 2 ビット誤り（3 番目と 5 番目）

	fmt.Printf("sent          = %08b  check=%t\n", sent, check(sent))
	fmt.Printf("1bit error    = %08b  check=%t  → 検出できた\n", recv1, check(recv1))
	fmt.Printf("2bit error    = %08b  check=%t  → 見逃した！\n", recv2, check(recv2))
}
```

**期待する出力**:

```
sent          = 10110010  check=true
1bit error    = 10111010  check=false  → 検出できた
2bit error    = 10011010  check=true   → 見逃した！
```

**ポイント**:
- 2 ビット反転 → 1 の個数が「+2 / 0 / -2」のいずれか → **偶奇は変わらない** → `check=true`。
- 受信側は「1 の個数は偶数」しか見ていないので、**2 ビット誤りを正常と誤判定**する。
- これが設問の答えの後半「2 ビットの誤りは検出できない」。

---

## まとめ: 設問 H20 秋 問6 の答え

| 誤りビット数 | 1 の個数の偶奇 | check の結果 | 意味            |
| ------------ | -------------- | ------------ | --------------- |
| 0（正常）    | 偶数のまま     | true         | 正常            |
| **1 ビット** | 偶数 → 奇数    | **false**    | **検出できる**  |
| 2 ビット     | 偶数 → 偶数    | true         | **見逃す**      |
| 3 ビット     | 偶数 → 奇数    | false        | 検出できる      |

- **奇数個の誤りは検出、偶数個の誤りは見逃す。**
- 検出できても**位置は特定できない → 訂正はできない**。
- ゆえに正解は「**1 ビットの誤りは検出できるが、その位置は分からない（訂正できない）。
  2 ビットの誤りは検出できない。**」

### 発展（余力があれば）

- **奇数パリティ**にすると何が変わるか（`% 2 == 1` に揃える）。検出能力は同じ。
- **誤り訂正**したいなら 1 ビットでは足りない → ハミング符号（複数パリティで位置を特定）。
- パリティは XOR の全ビット畳み込み: `b7 ^ b6 ^ ... ^ b0` が偶数パリティビットに一致する。
  `../xor` の交換/結合法則があるから順番を気にせず畳み込める。

---

# 第2部: 二次元パリティで「訂正」する（R5 秋 問4）

> 図のように 16 ビットのデータを 4×4 の正方形状に並べ、行と列にパリティビットを
> 付加する。この方式で何ビットまでの誤りを **訂正** できるか。
>
> **答え: 1 ビット。**

第1部（H20 問6）のパリティは「検出はできるが位置が分からない → 訂正できない」だった。
第2部では **行方向と列方向の 2 つ** のパリティを付ける。1 ビット誤ると
「異常な行」と「異常な列」がそれぞれ 1 つずつ出るので、**その交点** が誤りの位置。
位置が分かる＝反転して**訂正**できる。ここが第1部との決定的な違い。

配置（`data` を 4×4、右端の列に行パリティ、下端の行に列パリティ、右下は全体のパリティ）:

```
d d d d | 行パリティ
d d d d | 行パリティ
d d d d | 行パリティ
d d d d | 行パリティ
---------
列 列 列 列   ← 各列のパリティ
```

パリティ方式は第1部と同じ **偶数パリティ** で統一する。

---

## Module 6: 16 ビットを 4×4 に並べる

**ゴール**: 16 ビットのデータを `[4][4]int`（各セル 0/1）で表し、格子として表示する。

```go
package main

import "fmt"

// printGrid は 4x4 のビット格子を表示する
func printGrid(g [4][4]int) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(g[i][j])
		}
		fmt.Println()
	}
}

func main() {
	// 16 ビットのデータを 4x4 に並べる
	data := [4][4]int{
		{1, 0, 1, 1},
		{0, 1, 0, 1},
		{1, 1, 0, 0},
		{0, 0, 1, 1},
	}
	printGrid(data)
}
```

**期待する出力**:

```
1 0 1 1
0 1 0 1
1 1 0 0
0 0 1 1
```

**ポイント**:
- 第1部は `byte`（1 次元の 8 ビット）だったが、今回は 2 次元なので `[4][4]int` にする。
- 各セルは 0 か 1。ここに「行方向」と「列方向」の 2 つのパリティを後で足していく。

---

## Module 7: 行と列にパリティを付ける（送信側エンコード）

**ゴール**: 4×4 データに行パリティ（右端列）と列パリティ（下端行）を足して 5×5 にする。
できあがった 5×5 は **すべての行・すべての列が偶数** になる。

```go
package main

import "fmt"

var data = [4][4]int{
	{1, 0, 1, 1},
	{0, 1, 0, 1},
	{1, 1, 0, 0},
	{0, 0, 1, 1},
}

// encode は 4x4 データに行・列パリティを付けて 5x5 にする
func encode(data [4][4]int) [5][5]int {
	var g [5][5]int
	// 元データをコピー
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			g[i][j] = data[i][j]
		}
	}
	// 行パリティ: 各行(0..3)の 1 の個数の偶奇を右端(列 4)に入れる
	for i := 0; i < 4; i++ {
		sum := 0
		for j := 0; j < 4; j++ {
			sum += data[i][j]
		}
		g[i][4] = sum % 2
	}
	// 列パリティ: 各列(0..4)の 1 の個数の偶奇を下端(行 4)に入れる
	//   列 4 も含めるので、右下 g[4][4] は「全体のパリティ」になる
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 4; i++ {
			sum += g[i][j]
		}
		g[4][j] = sum % 2
	}
	return g
}

// printGrid5 は 5x5 のビット格子を表示する
func printGrid5(g [5][5]int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
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
}
```

**期待する出力**（右端の列＝行パリティ、下端の行＝列パリティ、右下＝全体のパリティ）:

```
1 0 1 1 1
0 1 0 1 0
1 1 0 0 0
0 0 1 1 0
0 0 0 1 1
```

**ポイント**:
- 行パリティ列 `[1,0,0,0]` … 1 行目だけ 1 の個数が奇数（3 個）なので 1。
- 列パリティ行 `[0,0,0,1,1]` … 各列を偶数に揃える。列 4 まで含めるので右下も計算される。
- できた 5×5 は **どの行も・どの列も 1 の個数が偶数**。これが「正常な状態」の指紋。

---

## Module 8: 受信側で行・列パリティを再検査（誤りなし）

**ゴール**: 受け取った 5×5 の「各行・各列が偶数か」を再計算する。誤りが無ければ全部 OK。

```go
package main

import "fmt"

var data = [4][4]int{
	{1, 0, 1, 1},
	{0, 1, 0, 1},
	{1, 1, 0, 0},
	{0, 0, 1, 1},
}

func encode(data [4][4]int) [5][5]int {
	var g [5][5]int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			g[i][j] = data[i][j]
		}
	}
	for i := 0; i < 4; i++ {
		sum := 0
		for j := 0; j < 4; j++ {
			sum += data[i][j]
		}
		g[i][4] = sum % 2
	}
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 4; i++ {
			sum += g[i][j]
		}
		g[4][j] = sum % 2
	}
	return g
}

// checkRows は各行の 1 の個数が奇数(=異常)なら true を返す
func checkRows(g [5][5]int) [5]bool {
	var bad [5]bool
	for i := 0; i < 5; i++ {
		sum := 0
		for j := 0; j < 5; j++ {
			sum += g[i][j]
		}
		bad[i] = sum%2 != 0
	}
	return bad
}

// checkCols は各列の 1 の個数が奇数(=異常)なら true を返す
func checkCols(g [5][5]int) [5]bool {
	var bad [5]bool
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 5; i++ {
			sum += g[i][j]
		}
		bad[j] = sum%2 != 0
	}
	return bad
}

func main() {
	g := encode(data)
	fmt.Println("行ごとの異常(true=奇数):", checkRows(g))
	fmt.Println("列ごとの異常(true=奇数):", checkCols(g))
}
```

**期待する出力**:

```
行ごとの異常(true=奇数): [false false false false false]
列ごとの異常(true=奇数): [false false false false false]
```

**ポイント**:
- `checkRows` / `checkCols` は「その行・列の 1 の個数が奇数か」を見るだけ（第1部の `check` の 2 次元版）。
- 誤りが無ければ全部 `false`。この「全 false」からのズレ方が、次のモジュールで位置を教えてくれる。

---

## Module 9: 1 ビットの誤り → 交点で位置特定 → 訂正できる（核心）

**ゴール**: 1 ビットだけ反転させると「異常な行 1 つ・異常な列 1 つ」になり、
その **交点** が誤りの位置。反転し返せば訂正できることを見る。

```go
package main

import "fmt"

var data = [4][4]int{
	{1, 0, 1, 1},
	{0, 1, 0, 1},
	{1, 1, 0, 0},
	{0, 0, 1, 1},
}

func encode(data [4][4]int) [5][5]int {
	var g [5][5]int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			g[i][j] = data[i][j]
		}
	}
	for i := 0; i < 4; i++ {
		sum := 0
		for j := 0; j < 4; j++ {
			sum += data[i][j]
		}
		g[i][4] = sum % 2
	}
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 4; i++ {
			sum += g[i][j]
		}
		g[4][j] = sum % 2
	}
	return g
}

func checkRows(g [5][5]int) [5]bool {
	var bad [5]bool
	for i := 0; i < 5; i++ {
		sum := 0
		for j := 0; j < 5; j++ {
			sum += g[i][j]
		}
		bad[i] = sum%2 != 0
	}
	return bad
}

func checkCols(g [5][5]int) [5]bool {
	var bad [5]bool
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 5; i++ {
			sum += g[i][j]
		}
		bad[j] = sum%2 != 0
	}
	return bad
}

func main() {
	g := encode(data)
	fmt.Println("=== 誤りを1ビット注入: (行1,列2) を反転 ===")
	g[1][2] ^= 1 // XOR で 1 ビット反転（第1部の flip と同じ考え方）
	printGrid5(g)

	rows := checkRows(g)
	cols := checkCols(g)
	fmt.Println("異常な行:", rows)
	fmt.Println("異常な列:", cols)

	// 異常な行・列を 1 つずつ探す
	r, c := -1, -1
	for i := 0; i < 5; i++ {
		if rows[i] {
			r = i
		}
	}
	for j := 0; j < 5; j++ {
		if cols[j] {
			c = j
		}
	}
	fmt.Printf("→ 交点 (%d,%d) が誤り。ここを反転して訂正\n", r, c)
	g[r][c] ^= 1 // 訂正

	fmt.Println("訂正後 異常な行:", checkRows(g))
	fmt.Println("訂正後 異常な列:", checkCols(g))
}

func printGrid5(g [5][5]int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Print(g[i][j])
		}
		fmt.Println()
	}
}
```

**期待する出力**:

```
=== 誤りを1ビット注入: (行1,列2) を反転 ===
1 0 1 1 1
0 1 1 1 0
1 1 0 0 0
0 0 1 1 0
0 0 0 1 1
異常な行: [false true false false false]
異常な列: [false false true false false]
→ 交点 (1,2) が誤り。ここを反転して訂正
訂正後 異常な行: [false false false false false]
訂正後 異常な列: [false false false false false]
```

**ポイント**:
- 1 ビット反転 → その行の偶奇が反転し、その列の偶奇も反転 → **異常が行 1 つ・列 1 つ** に出る。
- 異常な行 `1` と異常な列 `2` の **交点 (1,2)** が誤りの位置。第1部で分からなかった「位置」がここで確定する。
- 位置が分かれば `^= 1` で反転し返すだけ → **訂正完了**（訂正後は全 false に戻る）。
- `data` の値や反転する場所を変えても、1 ビットなら必ず「行 1・列 1」で交点が一意に決まる（試す）。

---

## Module 10: 2 ビットの誤り → 位置が定まらず訂正できない（限界）

**ゴール**: 2 ビット反転させると「異常な行 2 つ・列 2 つ」になり、交点候補が 4 つに増える。
どの 2 つが本当の誤りか区別できない → **訂正できない**（検出はできる）。

```go
package main

import "fmt"

var data = [4][4]int{
	{1, 0, 1, 1},
	{0, 1, 0, 1},
	{1, 1, 0, 0},
	{0, 0, 1, 1},
}

func encode(data [4][4]int) [5][5]int {
	var g [5][5]int
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			g[i][j] = data[i][j]
		}
	}
	for i := 0; i < 4; i++ {
		sum := 0
		for j := 0; j < 4; j++ {
			sum += data[i][j]
		}
		g[i][4] = sum % 2
	}
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 4; i++ {
			sum += g[i][j]
		}
		g[4][j] = sum % 2
	}
	return g
}

func checkRows(g [5][5]int) [5]bool {
	var bad [5]bool
	for i := 0; i < 5; i++ {
		sum := 0
		for j := 0; j < 5; j++ {
			sum += g[i][j]
		}
		bad[i] = sum%2 != 0
	}
	return bad
}

func checkCols(g [5][5]int) [5]bool {
	var bad [5]bool
	for j := 0; j < 5; j++ {
		sum := 0
		for i := 0; i < 5; i++ {
			sum += g[i][j]
		}
		bad[j] = sum%2 != 0
	}
	return bad
}

func printGrid5(g [5][5]int) {
	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
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
	fmt.Println("=== 誤りを2ビット注入: (1,2) と (2,3) を反転 ===")
	g[1][2] ^= 1
	g[2][3] ^= 1
	printGrid5(g)
	fmt.Println("異常な行:", checkRows(g))
	fmt.Println("異常な列:", checkCols(g))
}
```

**期待する出力**:

```
=== 誤りを2ビット注入: (1,2) と (2,3) を反転 ===
1 0 1 1 1
0 1 1 1 0
1 1 0 1 0
0 0 1 1 0
0 0 0 1 1
異常な行: [false true true false false]
異常な列: [false false true true false]
```

**ポイント**:
- 異常な行が `{1,2}`、異常な列が `{2,3}` の **2 つずつ** 出る。
- 交点の候補は (1,2)(1,3)(2,2)(2,3) の **4 つ**。実際の誤りは (1,2)(2,3) だが、
  (1,3)(2,2) が誤ったと考えても同じ検査結果になる → **どちらか区別できない**。
- 「変だ」とは分かる（検出できる）が、**位置を一意に決められない → 訂正できない**。
- だからこの方式で確実に訂正できるのは **1 ビットまで**。これが設問の答え。

---

## まとめ: 設問 R5 秋 問4 の答え

| 誤りビット数 | 異常な行 | 異常な列 | 位置の特定       | 訂正           |
| ------------ | -------- | -------- | ---------------- | -------------- |
| 0（正常）    | なし     | なし     | —                | 不要           |
| **1 ビット** | 1 つ     | 1 つ     | **交点で一意**   | **○ 訂正できる** |
| 2 ビット     | 2 つ     | 2 つ     | 候補 4 つで曖昧  | × 訂正できない（検出は可） |

- 1 ビット誤り → 行と列の異常の **交点** で位置が確定 → 反転して訂正。
- 2 ビット誤り → 交点候補が増えて位置が定まらない → 訂正不可。
- ゆえに正解は「**1 ビット**（選択肢ア）」。

### 第1部との対比（ここが理解の核）

| | 第1部 H20 問6（1 次元） | 第2部 R5 問4（二次元） |
| --- | --- | --- |
| パリティの軸 | 1 つ | 行＋列の 2 つ |
| 1 ビット誤り | 検出できる（位置不明） | **検出＋訂正できる**（交点で位置判明） |
| 得られる情報 | 「偶奇」1 ビット | 「異常な行」＋「異常な列」の座標 |

**位置を知るには複数のパリティが要る** — これが誤り「検出」と「訂正」を分ける本質。
この延長線上に、より少ない冗長ビットで位置を特定するハミング符号がある（`learn/hamming/` 候補）。

### 発展（余力があれば）

- **4 ビットが長方形の 4 隅で誤る**と、行も列も偶奇が戻り **検出すらできない**。
  例: (1,2)(1,3)(2,2)(2,3) を反転して `checkRows`/`checkCols` が全 false になるのを確認する。
- 右下の「全体のパリティ」`g[4][4]` だけで何が分かるか（＝第1部の 1 次元パリティ相当）を考える。

---

# 第3部: シリアル通信(UART)のパリティを擬似検証する（発展）

第1部で見た「7bitコード＋パリティ1bit」は、実は **シリアル通信(UART)** で線の上を流れる
フレームそのもの。ここでは Go でフレームを組み立て／分解し、パリティが誤りを検出する
（そして訂正はできない）様子を、通信の文脈で再確認する。試験問題ではなく、第1部の
実応用を体感する発展パート。

UART フレームの構造（アイドルは 1、下位ビット LSB から送る）:

```
スタート   データ(7 or 8bit, LSB→)      パリティ  ストップ
   0    | d0 d1 d2 d3 d4 d5 d6 (d7) |    p    |   1
```

設定は `8E1`（8bit データ / Even / stop1）のように書く。`N`＝パリティなし。
第3部は **ファイル分割** で進める: 共通部品を `uart.go`、`func main()` を `main.go` に置き、
`go run ./part3-uart` で実行する（Go は同じディレクトリの `.go` を 1 パッケージとしてまとめてビルドする）。

---

## Module 11: バイトをビット列に分解する（LSB first）

**ゴール**: 1 バイトを n 個のデータビットに分解する。UART は **下位ビットから** 送るのが肝。

`part3-uart/uart.go`:
```go
package main

// dataBits は byte を n 個のデータビットに分解する（LSB=下位ビットから）
func dataBits(b byte, n int) []int {
	bits := make([]int, n)
	for i := 0; i < n; i++ {
		bits[i] = int((b >> i) & 1) // i 番目の下位ビットを取り出す
	}
	return bits
}
```

`part3-uart/main.go`:
```go
package main

import "fmt"

func main() {
	// 'A' = 0x41 = 0b0100_0001
	fmt.Println("'A' の 8 データビット (LSB first):", dataBits('A', 8))
}
```

実行:
```bash
go run ./part3-uart
```

**期待する出力**:
```
'A' の 8 データビット (LSB first): [1 0 0 0 0 0 1 0]
```

**ポイント**:
- `'A'` は 0x41 =`0100 0001`。MSB から書くと `01000001` だが、UART は **bit0（LSB）から** 送るので `1 0 0 0 0 0 1 0` の順になる。
- `(b >> i) & 1` … 第1部の「右シフトで最下位を見る」と同じ。i を 0→7 と動かし下位から取り出す。
- `[]int`（スライス）は Python のリストに近い可変長の並び。`make([]int, n)` で長さ n を確保。

---

## Module 12: フレームを組み立てる（送信側）

**ゴール**: データにパリティ・スタート・ストップを付けて 1 フレームにする。
`8E1` なら 8bit データ＋1bit パリティで、線上は **9bit（＋start/stop で計 11bit）** になる。

`uart.go` に **追記**（先頭に `import "fmt"` を足す。showFrame で使う）:
```go
import "fmt"

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
	frame := []int{0}              // スタートビット
	data := dataBits(b, n)
	frame = append(frame, data...) // データ（LSB first）
	if mode != "N" {
		frame = append(frame, parityBit(data, mode)) // パリティ
	}
	frame = append(frame, 1)       // ストップビット
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
```

`main.go` を差し替え:
```go
package main

import "fmt"

func main() {
	fmt.Print("8E1  'A': ")
	showFrame(buildFrame('A', 8, "E"))
	fmt.Print("8O1  'A': ")
	showFrame(buildFrame('A', 8, "O"))
	fmt.Print("8N1  'A': ")
	showFrame(buildFrame('A', 8, "N"))
}
```

**期待する出力**:
```
8E1  'A': 0 1 0 0 0 0 0 1 0 0 1   (11bit)
8O1  'A': 0 1 0 0 0 0 0 1 0 1 1   (11bit)
8N1  'A': 0 1 0 0 0 0 0 1 0 1   (10bit)
```

**ポイント**:
- 各行の左端が **スタート(0)**、右端が **ストップ(1)**。その間がデータ(LSB first)＋パリティ。
- `'A'` はデータ中の 1 が 2 個（偶数）→ 偶数パリティ=0、奇数パリティ=1。8N1 はパリティ無しで 1bit 短い。
- `8E1/8O1` は 11bit、`8N1` は 10bit。**8bit＋パリティ＝線上 9bit**（＋start/stop）が確かにできている。
- `append(frame, data...)` の `...` は「スライスを展開して連結」。Go でリストをつなぐ書き方。

---

## Module 13: 受信側で分解・パリティ検査（誤りなし）

**ゴール**: フレームからデータを取り出し、パリティを再計算して照合、バイトに復元する。

`uart.go` に **追記**:
```go
// receive はフレームを分解し、復元バイトとパリティ検査結果(true=正常)を返す
func receive(frame []int, n int, mode string) (byte, bool) {
	data := frame[1 : 1+n] // スタートの次から n 個がデータ
	ok := true
	if mode != "N" {
		p := frame[1+n]                 // 受信したパリティ
		ok = parityBit(data, mode) == p // 再計算と一致するか
	}
	var b byte
	for i := 0; i < n; i++ {
		b |= byte(data[i]) << i // LSB first でバイトに戻す
	}
	return b, ok
}
```

`main.go` を差し替え:
```go
package main

import "fmt"

func main() {
	frame := buildFrame('A', 8, "E")
	b, ok := receive(frame, 8, "E")
	fmt.Print("受信フレーム: ")
	showFrame(frame)
	fmt.Printf("復元バイト = %q, パリティ検査 = %t\n", b, ok)
}
```

**期待する出力**:
```
受信フレーム: 0 1 0 0 0 0 0 1 0 0 1   (11bit)
復元バイト = 'A', パリティ検査 = true
```

**ポイント**:
- `frame[1 : 1+n]` … スライスの一部を切り出す（Python の `frame[1:1+n]` と同じ）。スタートを飛ばして n 個。
- 受信側はパリティを **自分で計算し直して** 届いた p と比べる。一致すれば正常とみなす。
- `b |= byte(data[i]) << i` … LSB first のビットを元のバイトに戻す（第1部 `encode` の逆向き）。
- `receive` が **2 つの値**（バイトと bool）を返すのは Go の多値返却。`b, ok := ...` で受ける。

---

## Module 14: 1 ビットの誤り → パリティエラーを検出

**ゴール**: フレームの途中で 1 ビット化けさせると、受信側のパリティ検査が false になる。

`main.go` を差し替え（`uart.go` はそのまま）:
```go
package main

import "fmt"

func main() {
	frame := buildFrame('A', 8, "E")
	frame[3] ^= 1 // 通信路で 1 ビット化けた（データの一部）

	b, ok := receive(frame, 8, "E")
	fmt.Print("化けたフレーム: ")
	showFrame(frame)
	fmt.Printf("復元バイト = %q, パリティ検査 = %t\n", b, ok)
}
```

**期待する出力**:
```
化けたフレーム: 0 1 0 1 0 0 0 1 0 0 1   (11bit)
復元バイト = 'E', パリティ検査 = false
```

**ポイント**:
- 1 ビット反転 → データの 1 の個数の偶奇が変わる → 受信側の再計算パリティが届いた p と食い違う → **`false`＝検出成功**。
- 復元は `'A'→'E'` と化けているが、`ok=false` なので「この文字は信用できない」と分かる。実機ではここで **再送要求** する。
- ただし **どのビットが化けたかは分からない**（第1部と同じ）。検出止まりで訂正はできない。

---

## Module 15: 2 ビットの誤り → 見逃す（サイレント破損）

**ゴール**: 2 ビット化けると偶奇が元に戻り、パリティ検査が true のまま通ってしまう。

`main.go` を差し替え:
```go
package main

import "fmt"

func main() {
	frame := buildFrame('A', 8, "E")
	frame[3] ^= 1 // 1 つ目の誤り
	frame[5] ^= 1 // 2 つ目の誤り

	b, ok := receive(frame, 8, "E")
	fmt.Print("化けたフレーム: ")
	showFrame(frame)
	fmt.Printf("復元バイト = %q, パリティ検査 = %t\n", b, ok)
}
```

**期待する出力**:
```
化けたフレーム: 0 1 0 1 0 1 0 1 0 0 1   (11bit)
復元バイト = 'U', パリティ検査 = true
```

**ポイント**:
- 2 ビット反転 → 偶奇は元どおり → 受信側は **正常と誤判定（`true`）**。
- しかし復元は `'A'→'U'` と **別の文字に化けている** のに気づけない＝**サイレント破損**。
- これが第1部「偶数個の誤りは見逃す」の、通信での怖さ。パリティだけでは足りず、実務では CRC 等の強い検出符号を上位で併用する。

---

## Module 16: パリティなし(8N1) → そもそも検出しない（まとめ）

**ゴール**: `N`（パリティなし）だと 1 ビット誤りすら検出できないことを確認し、パリティの意味を締める。

`main.go` を差し替え:
```go
package main

import "fmt"

func main() {
	frame := buildFrame('A', 8, "N") // パリティなし
	frame[3] ^= 1                    // 1 ビット化けた

	b, ok := receive(frame, 8, "N")
	fmt.Print("化けたフレーム: ")
	showFrame(frame)
	fmt.Printf("復元バイト = %q, パリティ検査 = %t\n", b, ok)
}
```

**期待する出力**:
```
化けたフレーム: 0 1 0 1 0 0 0 1 0 1   (10bit)
復元バイト = 'E', パリティ検査 = true
```

**ポイント**:
- パリティが無いので `receive` は常に `ok=true`。1 ビット化けても **検出できず素通り**。
- パリティ 1bit を足すコスト（8N1=10bit → 8E1=11bit、約 1 割の速度低下）と引き換えに、1bit 誤りの検出能力を得ている、というトレードオフ。

### まとめ: UART のパリティ（第1部の実応用）

| 設定 | 線上のビット数 | 1bit 誤り | 2bit 誤り |
| --- | --- | --- | --- |
| 8N1 | 10bit | 検出できない | 検出できない |
| 8E1 / 8O1 | 11bit | **検出できる**（位置不明・訂正不可） | 見逃す |

- UART のパリティは第1部の 1 次元パリティそのもの。**1bit 誤りを検出して再送** が基本戦略。
- 偶数/奇数は検出能力は同じ（奇数個検出・偶数個見逃し）。選択は互換性・慣習による。
- 訂正まで欲しければ冗長を増やす（第2部の二次元パリティ／ハミング符号）。
  これで「検出のみ(1次元) → 訂正可(二次元) → 実通信での使われ方(UART)」がひとつながりになる。

---

## コントリビューション

- `main` ブランチへの直接 commit & push は管理者（tkm112345）のみ。
- それ以外の方は fork して Pull Request を作成してください。
  `main` への取り込みはレビュー承認後のみ可能です（ブランチ保護でレビュー必須）。
