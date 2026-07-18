# ハンズオン: パリティビットを Go で実験する

応用情報技術者試験の **パリティビット問題（H20 秋 問6）** を「腑に落とす」ための写経ハンズオン。

> 7 ビットのコードと 1 ビットのパリティビットからなる 8 ビットのデータで発生した
> 誤りに関する記述として、適切なものはどれか
>
> **答え: 1 ビットの誤りは検出できるが、その位置は分からない（訂正できない）。**
> **2 ビットの誤りは検出できない。**

この「なぜ？」を、小さな Go プログラムで 1 つずつ目で確認する。
XOR ハンズオン（`../xor`）の続き。パリティは「1 の個数の偶奇」であり、
その計算は XOR の交換法則・結合法則の応用そのもの。

## 進め方

- `main.go` を各モジュールのコードに**写経（打ち込み）**して書き換える
- `go run main.go` で実行し、「期待する出力」と一致するか確認する
- コードより **「なぜそうなるか」** に注目する（下の解説を読む）

パリティ方式は本ハンズオンでは **偶数パリティ**（1 の個数を偶数に揃える）で統一する。
パリティビットは **最下位ビット（一番右）** に付ける。

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

## コントリビューション

- `main` ブランチへの直接 commit & push は管理者（tkm112345）のみ。
- それ以外の方は fork して Pull Request を作成してください。
  `main` への取り込みはレビュー承認後のみ可能です（ブランチ保護でレビュー必須）。
