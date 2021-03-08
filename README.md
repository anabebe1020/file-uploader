## description

### メイン機能

Slackの特定のチャンネルにファイルを添付すると、
サーバーの指定のディレクトリにファイルをアップロードし、
ダウンロードURLをチャンネルに返します。
-> ダウンロードURL：`BaseURL/HashDir/File.ex`

また、サーバーのディスクが逼迫している場合、アップロードせずにエラーを返します。

### オプション機能

`/list`
-> SlackのSlashCommandから呼び出します。
   アップロード済みのファイルの一覧をチャンネルに返します。

`/del [dir name/file name]`
-> SlackのSlashCommandから呼び出します。
   アップロード済みのファイルの中から`[dir name/file name]`に一致したファイルを削除します。

## preperation

### 1. アップロードディレクトリを用意

実際にアップロードしたファイルを格納するディレクトリを作成しておきます。
```
$ mkdir /sample/path/download
```

### 2. Proxyサーバーを立てる

`/file-uploader`
-> メイン機能、Slackの"file_shared"イベントをキャッチします。

`/file-cmd`
-> オプション機能、SlashCommandの受け口。

### 3. "config.toml"の作成

sample-config.tomlをコピーし、config.tomlを作成します。
```
$ sudo cp -p ./sample-config.toml ./config.toml
```
- [Slack] BotToken Slack Appから取得したbot tokenを指定します。
- [Slack] Channel  該当のチャンネルのIDを指定します。チャンネル名は非推奨です。
- [File]  Path     アップロードするファイルを保存するディレクトリを指定します。
- [File]  DLURL    ダウンロードURLのベースURLを指定します。

## ビルド & 実行

ビルドして出力された実行ファイルを実行します。必要に応じてサービス化を行ってください。
```
$ go build
$ ./FileUploader
```

## 自動起動設定（linux）

サーバーが再起動された時に自動で起動するように設定します。
```
$ cp -p ./service/FileUploader.service /etc/systemd/system/
$ systemctl enable FileUploader.service
$ systemctl start FileUploader.service
```

## 自動削除設定（linux）

アップロードから14日経過したファイルをフォルダごと削除する、というチャックを毎日00:00に行う。
```
$ crontab -e

0 0 * * * root /bin/find /xxx/xxx/download/ -maxdepth 1 -mtime +14 | xargs sudo rm -rf
```