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

### go言語について(ChatGPTより生成)
Go言語における以下の概念について、それぞれ解説し、コード例を示します。

### 1. パッケージ

#### 概要
- パッケージはコードを整理する基本単位です。ファイルは必ずどこかのパッケージに属し、そのパッケージはディレクトリで構成されます。

#### 例
ディレクトリ構成：
```
math/
  add.go
main.go
```

`math/add.go`:
```go
package math

func Add(a, b int) int {
    return a + b
}
```

`main.go`:
```go
package main

import (
    "fmt"
    "./math"
)

func main() {
    result := math.Add(2, 3)
    fmt.Println("Result:", result)
}
```

### 2. モジュール

#### 概要
- モジュールはパッケージの集合体です。Go 1.11以降では、モジュールシステムを利用して依存管理を行います。

#### 例
`go.mod`:
```
module example.com/yourproject

go 1.18
```

このファイルはモジュールの名前とGoのバージョンを指定します。`go mod init`コマンドを使うと作成できます。

### 3. 構造体

#### 概要
- 構造体は複数の型を組み合わせた新しい型を定義するために使用します。

#### 例
```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    p := Person{Name: "Alice", Age: 30}
    fmt.Println("Name:", p.Name, "Age:", p.Age)
}
```

### 4. 循環参照とその解決法

#### 概要
- 循環参照とは2つ以上のパッケージがお互いに依存している状態です。Goではパッケージ間の循環参照を許可していません。

#### 解決法
- パッケージの設計を見直し、共通部分を別パッケージに分離する。
- インタフェースを使用して依存を逆転させる。

例として、以下のような2つのパッケージが循環参照している場合：

- `packageA`が`packageB`に依存
- `packageB`が`packageA`に依存

解決方法の一つとして、共通のインタフェースを使って依存関係を逆転させることが考えられます。

### 5. インタフェース

#### 概要
- インタフェースはメソッドの集合で、Goでは型の実装を抽象化する方法を提供します。

#### 例
```go
package main

import "fmt"

// インタフェースの定義
type Speaker interface {
    Speak() string
}

// 構造体1
type Dog struct{}

func (d Dog) Speak() string {
    return "Woof!"
}

// 構造体2
type Cat struct{}

func (c Cat) Speak() string {
    return "Meow!"
}

func main() {
    var s Speaker

    s = Dog{}
    fmt.Println(s.Speak())

    s = Cat{}
    fmt.Println(s.Speak())
}
```

この例では、`Speaker`インタフェースが定義され、そのメソッド`Speak()`を`Dog`と`Cat`の構造体が実装しています。

それぞれの概念について、これらの例が基本的な理解の助けになることを願っています。さらなる詳細が必要であれば、特定のトピックについての質問をお知らせください。