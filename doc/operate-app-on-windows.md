# Operate app on Windows

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
