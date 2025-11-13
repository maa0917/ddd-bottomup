# コーディング規約とスタイル

## パッケージ構成
- `package entity` - ドメインエンティティ
- `package valueobject` - 値オブジェクト 
- `package repository` - リポジトリインターフェース
- `package service` - ドメインサービス
- インポートパス: `ddd-bottomup/domain/valueobject` 等

## 命名規約

### 型命名
- **ID型**: `UserID`, `CircleID`, `ShipmentID` など `{Entity}ID` 形式
- **エンティティ**: `User`, `Circle`, `Shipment` など
- **値オブジェクト**: `Email`, `FullName`, `CircleName` など

### 関数命名
- **コンストラクタ**: `NewUser()`, `NewCircleID()` など `New{Type}` 形式
- **復元関数**: `ReconstructUser()`, `ReconstructUserID()` など `Reconstruct{Type}` 形式
- **アクセサ**: `ID()`, `Name()`, `Email()` など
- **述語**: `Equals()`, `IsPremium()`, `CanAddMember()` など

### フィールド命名
- プライベートフィールド: `id`, `name`, `email`, `isPremium` など小文字開始
- 構造体フィールド: タブインデント使用

## コード構造パターン

### ID型の実装パターン
```go
type UserID struct {
    value string
}

func NewUserID() *UserID {
    return &UserID{value: uuid.New().String()}
}

func ReconstructUserID(value string) (*UserID, error) {
    // バリデーション + 復元
}

func (u *UserID) Value() string { return u.value }
func (u *UserID) Equals(other *UserID) bool { /* 比較 */ }
func (u *UserID) String() string { return u.value }
```

### エンティティの実装パターン
```go
type User struct {
    id        *UserID
    name      *valueobject.FullName
    email     *valueobject.Email
    isPremium bool
}

func NewUser(...) *User { /* 新規作成 */ }
func ReconstructUser(...) *User { /* DB等からの復元 */ }

// アクセサメソッド
func (u *User) ID() *UserID { return u.id }
// 変更メソッド
func (u *User) ChangeName(name *valueobject.FullName) { u.name = name }
```

## エラーハンドリング
- `errors.New()` でカスタムエラーメッセージ
- UUIDパース失敗時のバリデーション
- 空文字チェック等の基本バリデーション

## テストファイル
- `*_test.go` 命名規約
- `TestCreateUserUseCase_Execute_Success` 形式のテスト関数名
- `Test{UseCase}_{Method}_{Scenario}` パターン