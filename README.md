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

### ログの基本構文
```go
err = LogMessage(INFO, "<INFOメッセージ>")
```

```go
errMessage := fmt.Sprintf("<エラーメッセージ> %v\n", err)
		err = LogMessage(ERRO, errMessage)
```



### 各種設定箇所の確認
ログファイル
```
cat /var/log/otdm-package.log
```

otdm statusで表示するjsonファイル
```
cat /tmp/otdm_staus.json
```

wireguardのコンフィグファイル
```
su - && cat /etc/wireguard
```

で可能

### 開発中使用
- make コマンド
- go言語
- wireguard

# README

## github関連のコマンド
<details open>

<summary>github関連のコマンド</summary>

### 大まかな流れ

1. リモートからローカルにpull
2. 作業するブランチを作成し、移動
3. 作業
4. リモートにpush
5. github上でプルリクエスト
6. merge

### 作業

1. リモートからローカルにpull
```bash
git pull
```
2. 作業するブランチを作成し、移動
確認
```bash
git branch
```
もしくは
```bash
git branch -vv
```

作成・移動
```bash
git checkout -b <新規に作成するブランチ名> | git checkout <既存のブランチ名>
```

3. 作業
   
最高の激励の言葉
> Shut the fuck up and write some code

4. リモートにpush
```bash
git add .
```
```bash
git commit -m "<コミット内容>"
```
```bash
git push -u origin <新規(リモートにない)ブランチ名> | git push
```

5. github上でプルリクエスト
webページ上での作業なため、割愛
6. merge
webページ上での作業なため、割愛
7. (必要であれば)リモートのブランチを削除
webページ上での作業なため、割愛
8. ローカルのブランチを削除する場合
ブランチを削除する場合は
```bash
git branch -d <削除するブランチ名>
```

### こんな時は

### ブランチを間違えて作成した、もしくは作成前に作業してしまった

```bash
git stash -u
```
で一度退避したのちに、ブランチを作成、その後
```bash
git stash pop
```
で退避したものを反映

### リモートにしかないブランチを取り込む

```bash
git fetch
git branch -r
git checkout -b <取り込んだ際のブランチ名> <元となるブランチ名>  
```



> ubuntu-103@ubuntu-103:~/sandbox/otdmGui$ git fetch 
> ubuntu-103@ubuntu-103:~/sandbox/otdmGui$ git branch 
> \* main
> ubuntu-103@ubuntu-103:~/sandbox/otdmGui$ git branch -a
> \* main
>   remotes/origin/HEAD -> origin/main
>   remotes/origin/MainEntry
>   remotes/origin/ReworkBase
>   remotes/origin/main
> ubuntu-103@ubuntu-103:~/sandbox/otdmGui$ git checkout -b MainEntry origin/MainEntry
> branch 'MainEntry' set up to track 'origin/MainEntry'.
> Switched to a new branch 'MainEntry'
> ubuntu-103@ubuntu-103:~/sandbox/otdmGui$ git branch -vv
> \* MainEntry ef767f6 [origin/MainEntry] 修正：【未解決】Vue router構築中
>   main│     b2af626 [origin/main] Merge pull request #4 from Otdm-Project/ReworkBase
> ubuntu-103@ubuntu-103:~/sandbox/otdmGui$ 




## 本プロジェクトの構成(一部)


</details>



## REDMEカリカリするときの楽するコピペ
```
<details>

<summary>▶の右に出るメッセージ</summary>

[ここに内容を記述する]

</details>
```

開始タグを`<details open>`にすることで開いた状態を初期状態にする。
タグ前後には一行空白をいれないと、時々マークダウンのタグが無効化されていることがあるので注意