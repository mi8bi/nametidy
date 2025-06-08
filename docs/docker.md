# nametidy Docker ガイド

nametidyのDocker環境での使用方法について説明します。

## 📋 目次

- [クイックスタート](#クイックスタート)
- [開発環境](#開発環境)
- [本番環境](#本番環境)
- [使用例](#使用例)
- [トラブルシューティング](#トラブルシューティング)

## 🚀 クイックスタート

### 前提条件
- Docker 20.10 以上
- Docker Compose 2.0 以上

### 1. リポジトリのクローン
```bash
git clone https://github.com/mi8bi/nametidy.git
cd nametidy
```

### 2. 即座に実行
```bash
# ヘルプを表示
docker-compose run --rm nametidy --help

# ファイルをクリーンアップ（./filesディレクトリ）
docker-compose run --rm nametidy clean -p /workspace
```

## 🛠 開発環境

### 開発環境の起動
```bash
# 開発コンテナを起動
docker-compose -f docker-compose.dev.yml up -d

# コンテナに接続
docker-compose -f docker-compose.dev.yml exec nametidy-dev bash
```

### コンテナ内での開発
```bash
# アプリケーションをビルド
go build -o nametidy .

# テストを実行
go test ./...

# ベンチマークを実行
go test -bench=. ./...

# アプリケーションを実行
./nametidy clean -p ./test_files
```

### ホットリロード
複数の方法でホットリロードを利用できます：

#### 方法1: Air使用（推奨）
```bash
# コンテナ内でAirを起動（ファイル変更時に自動再ビルド）
air
```

#### 方法2: カスタム開発スクリプト
```bash
# コンテナ内で開発用ウォッチスクリプトを実行
chmod +x scripts/dev-watch.sh
./scripts/dev-watch.sh
```

#### 方法3: 手動ビルド
```bash
# ファイル変更後に手動でリビルド
go build -o nametidy . && ./nametidy --help
```

### 開発環境の停止
```bash
docker-compose -f docker-compose.dev.yml down
```

## 🏭 本番環境

### イメージのビルド
```bash
# 本番用イメージをビルド
docker build -t nametidy:latest .

# または docker-compose でビルド
docker-compose build
```

### アプリケーションの実行
```bash
# Docker run で直接実行
docker run --rm -v "$(pwd)/files:/workspace" nametidy:latest clean -p /workspace

# Docker Compose で実行
docker-compose run --rm nametidy clean -p /workspace
```

### バックグラウンドでのサービス実行
```bash
# サービスとして起動（将来のWeb UI用）
docker-compose up -d

# ログの確認
docker-compose logs -f nametidy

# サービスの停止
docker-compose down
```

## 📖 使用例

### 基本的なファイルクリーンアップ
```bash
# ローカルのfilesディレクトリをクリーンアップ
mkdir -p files
echo "hello world (1).txt" > "files/hello world (1).txt"
echo "test file.doc" > "files/test file.doc"

docker-compose run --rm nametidy clean -p /workspace -v
```

### ドライランモードでの確認
```bash
# 実際には変更せず、変更予定を表示
docker-compose run --rm nametidy clean -p /workspace -d
```

### ファイルに番号を付与
```bash
# ファイルに3桁の番号を付与
docker-compose run --rm nametidy number -p /workspace -n 3
```

### 変更の取り消し
```bash
# 最後の操作を取り消し
docker-compose run --rm nametidy undo -p /workspace
```

### カスタムディレクトリでの実行
```bash
# 任意のディレクトリをマウントして実行
docker run --rm -v "/path/to/your/files:/workspace" nametidy:latest clean -p /workspace
```

## 🐛 トラブルシューティング

### よくある問題と解決方法

#### 1. 権限エラー
```bash
# 権限の問題が発生した場合
sudo chown -R $USER:$USER ./files

# または、コンテナ内でユーザーIDを指定
docker run --rm --user $(id -u):$(id -g) -v "$(pwd)/files:/workspace" nametidy:latest clean -p /workspace
```

#### 2. ファイルが見つからない
```bash
# マウントパスを確認
docker-compose run --rm nametidy ls -la /workspace

# ローカルディレクトリの確認
ls -la ./files
```

#### 3. ビルドエラー
```bash
# Go version関連のエラーの場合
docker-compose -f docker-compose.dev.yml build --no-cache

# キャッシュをクリアしてリビルド
docker system prune -f

# 特定のGo versionを指定してビルド
docker build --build-arg GO_VERSION=1.23 -f Dockerfile.dev -t nametidy-dev:latest .
```

#### 4. Air/ホットリロードの問題
```bash
# Airが使用できない場合、代替のウォッチスクリプトを使用
./scripts/dev-watch.sh

# または手動でのファイル監視
while true; do inotifywait -e modify *.go && go build -o nametidy .; done
```

#### 4. メモリ不足
```bash
# Docker設定でメモリを増やす、または軽量版を使用
docker run --rm --memory=128m -v "$(pwd)/files:/workspace" nametidy:latest clean -p /workspace
```

### ログの確認
```bash
# アプリケーションログ
docker-compose logs nametidy

# 詳細なログ
docker-compose logs -f --tail=100 nametidy

# システムログ
docker system events
```

### パフォーマンスの監視
```bash
# コンテナのリソース使用量
docker stats

# 特定のコンテナの統計
docker stats nametidy-app
```

## 🔧 高度な設定

### 環境変数のカスタマイズ
```bash
# .env ファイルを作成
cat > .env << EOF
TZ=America/New_York
LANG=en_US.UTF-8
LOG_LEVEL=debug
EOF

# 環境変数を指定して実行
docker-compose --env-file .env run --rm nametidy clean -p /workspace
```

### カスタムネットワークの使用
```bash
# カスタムネットワークを作成
docker network create nametidy-net

# ネットワークを指定して実行
docker run --rm --network nametidy-net -v "$(pwd)/files:/workspace" nametidy:latest clean -p /workspace
```

### マルチアーキテクチャビルド
```bash
# ARM64とAMD64の両方をサポート
docker buildx build --platform linux/amd64,linux/arm64 -t nametidy:latest .
```

## 📚 参考リンク

- [Docker公式ドキュメント](https://docs.docker.com/)
- [Docker Compose公式ドキュメント](https://docs.docker.com/compose/)
- [Go Docker Best Practices](https://docs.docker.com/language/golang/)
- [nametidy メインドキュメント](../README.md)

## 🤝 貢献

Docker環境に関する改善提案や問題報告は、GitHubのIssueでお知らせください。

---

> **注意**: このドキュメントは`nametidy`のDocker環境での使用方法を説明しています。アプリケーション自体の使用方法は[メインのREADME](../README.md)を参照してください。