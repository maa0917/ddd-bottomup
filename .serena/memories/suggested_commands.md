# 推奨開発コマンド

## 依存関係管理
```bash
go mod tidy                    # 依存関係のダウンロードと整理
go mod verify                  # モジュール依存関係の検証
```

## ビルドと実行
```bash
go build -o bin/app .         # アプリケーションのビルド
go run main.go                # アプリケーションの直接実行
```

## テスト
```bash
go test ./...                 # 全テストの実行
go test ./usecase            # 特定パッケージのテスト実行
go test -v ./usecase         # 詳細出力でのテスト実行
go test -run TestName        # 特定テストの実行
```

## コード品質
```bash
go fmt ./...                 # コードフォーマット
go vet ./...                 # 静的解析
```

## アプリケーション起動
- HTTPサーバーはポート8080で起動
- `/health` エンドポイントでヘルスチェック可能

## システムコマンド (macOS/Darwin)
- `ls` - ファイル一覧
- `cd` - ディレクトリ移動  
- `grep` - 文字列検索
- `find` - ファイル検索
- `git` - Gitコマンド

## 開発ワークフロー
1. `go mod tidy` で依存関係更新
2. コード修正
3. `go test ./...` でテスト実行
4. `go fmt ./...` でフォーマット
5. `go vet ./...` で静的解析
6. `go run main.go` で動作確認