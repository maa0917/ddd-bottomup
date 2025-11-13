# DDD Bottom-Up プロジェクト概要

## プロジェクトの目的
このプロジェクトは、ドメイン駆動設計（DDD）の底上げアプローチを使用したサークル管理システムです。Goで実装された Clean Architecture のサンプルです。

## 技術スタック
- **言語**: Go 1.24
- **依存関係**: 
  - github.com/go-sql-driver/mysql v1.9.3 (MySQL接続)
  - github.com/google/uuid v1.6.0 (UUID生成)
- **データベース**: MySQL (SQLスキーマあり)
- **アーキテクチャ**: Domain-Driven Design + Clean Architecture

## アーキテクチャ構造

### レイヤー構造
1. **ドメイン層** (`/domain`): 
   - エンティティ (User, Circle, CircleMembers, Shipment, Baggage等)
   - バリューオブジェクト (Email, FullName, CircleName, Money)
   - リポジトリインターフェース
   - ドメインサービス

2. **ユースケース層** (`/usecase`):
   - アプリケーションサービス
   - CreateUser, GetUser, UpdateUser, DeleteUser等

3. **インフラストラクチャ層** (`/infrastructure`):
   - リポジトリ実装 (メモリ版とMySQL版)
   - データベーススキーマ

4. **プレゼンテーション層** (`/presentation`):
   - HTTPハンドラー
   - ルーター

### 主要パターン
- **仕様パターン**: 複雑なビジネスルール実装
- **リポジトリパターン**: データアクセスの抽象化
- **エンティティ設計**: 豊富なドメインモデル
- **強い型付け**: 専用ID型 (UserID, CircleID等)

## エンドポイント
- POST /users - ユーザー作成
- GET /users/{id} - ユーザー取得
- PUT /users/{id} - ユーザー更新
- DELETE /users/{id} - ユーザー削除
- GET /health - ヘルスチェック

## データベース
- users テーブル: ユーザー情報
- circles テーブル: サークル情報
- circle_members テーブル: サークルメンバーシップ
- 自動トリガー: メンバー数の自動更新