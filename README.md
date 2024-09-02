# otdm-linux-app
## 開発中メモ
### ビルド～実行に関して
ビルドして実行を行う場合は
otdm-package下で
```
make build
```
したのちに、
```
sudo cp otdm /usr/local/bin/
```
でコマンドを実行できるようにする。

### 開発中使用
- make コマンド
- go言語
- wireguard

### パッケージ内で使用

### go言語について
#### 用語
- パッケージ
Go言語のパッケージは、関連するGoのコード(関数、方、変数…)をひとまとめにするための単位。何かの機能をまとめて、パッケージとする。
例:エントリポイントになるmainパッケージ

- モジュール

#### コード
- パッケージ宣言
```go
package main

func main(){
    println("Hello Word!")
}
```

- インポート
他のパッケージを使用するためには`import`文を用いる
```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Go!")
}
```
例）
- fmt(使用例:up.go):パッケージ fmt は、C の printf および scanf に類似した関数を使用してフォーマットされた I/O を実装します。
- os(使用例:websocket.go):パッケージ os は、オペレーティング システム機能へのプラットフォームに依存しないインターフェイスを提供します。

- 変数宣言
変数は`var`を用いて宣言を行う
コード例
```go
package main

import "fmt"

func main() {
    var number int = 42
    fmt.Println(number)
}
```

短縮系を用いても可
```go
func main() {
    number := 42
    fmt.Println(number)
}
```
「`:=`」...何だお前！！！
**型推論、「:=」**
intやstringを指定せずともこれを用いてデータを代入し、代入するデータの種類からデータ型を推論する