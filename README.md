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

### ジャーナルの確認
```
journalctl -t otdm-package
```
で可能

### 開発中使用
- make コマンド
- go言語
- wireguard

### パッケージ内で使用