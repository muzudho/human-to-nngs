# Manual

## Prepare App

Pre-install:  

```shell
go mod init

go get github.com/pelletier/go-toml
```

Telnet:  

```shell
# Go言語 は 個人作成の同名のライブラリがいっぱいあるので 一番良さそうなのを検索してください。
go get -v -u github.com/reiver/go-telnet
```

Build:  

```shell
go build
```

## Run App

Start:  

```shell
# 白番の例
human-to-nngs

# 黒番の例
human-to-nngs --entry input/humanb.entryConf.toml
```

## Operate App


白番から黒番へ対局を申し込むなら:  

```shell
match playerb B 19 40 0
```

もし `playerb` がまだログインしていなければ、以下のように `5 Error` メッセージが返ってきて　不成立になる。  

```shell
5 No user named "playerb" is logged in.
```

そこで、白番と黒番は 両方ログインしておく必要がある。  


黒番には、以下のような `9 Info` が流れてくる。  

```shell
9 Use <match playera W 19 40 0> or <decline playera> to respond.
```

そこで黒番は以下のように打鍵する。  

```shell
match playera W 19 40 0
```
