## description

### メイン機能

Slackの特定のチャンネルにファイルを添付すると、
サーバーの指定のディレクトリにファイルをアップロードし、
ダウンロードURLをチャンネルに返します。
-> ダウンロードURL：`BaseURL/HashDir/File.ex`

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
mkdir /sample/path/download
```

### 2. Proxyサーバーを立てる

`/file-uploader`
-> メイン機能、Slackの"file_shared"イベントをキャッチします。

`/file-cmd`
-> オプション機能、SlashCommandの受け口。

### 3. "config.toml"の作成

sample-config.tomlをコピーし、config.tomlを作成します。
```
cp -p ./sample-config.toml ./config.toml
```
- [Slack] BotToken Slack Appから取得したbot tokenを指定します。
- [Slack] Channel  該当のチャンネルのIDを指定します。チャンネル名は非推奨です。
- [File]  Path     アップロードするファイルを保存するディレクトリを指定します。
- [File]  DLURL    ダウンロードURLのベースURLを指定します。

## ビルド & 実行

ビルドして出力された実行ファイルを実行します。必要に応じてサービス化を行ってください。
```
go build
./FileUploader
```
