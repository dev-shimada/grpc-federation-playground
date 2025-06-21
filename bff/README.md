# gRPC Federation BFF Server

grpc-federationを使用してBackend For Frontend（BFF）パターンを実装したgRPCサーバーです。複数のマイクロサービスからデータを統合し、単一のAPIエンドポイントとしてクライアントに提供します。

## 🏗️ アーキテクチャ概要

```
┌─────────────┐     ┌─────────────────┐     ┌─────────────────┐
│   Client    │────▶│   BFF Server    │────▶│  User Service   │
│             │     │   (Port 8080)   │     │   (Port 8081)   │
└─────────────┘     │                 │     └─────────────────┘
                    │                 │     
                    │                 │     ┌─────────────────┐
                    │                 │────▶│ Message Service │
                    │                 │     │   (Port 8082)   │
                    └─────────────────┘     └─────────────────┘
```

### 主要コンポーネント

- **BFF Server (Port 8080)**: grpc-federationを使用して複数のサービスからデータを統合
- **User Service (Port 8081)**: ユーザー情報を管理
- **Message Service (Port 8082)**: メッセージ情報を管理

## 🚀 主要機能

### GetMessage API
単一のRPCコールで、メッセージ情報とそれに関連するユーザー情報を同時に取得できます。

- **入力**: `message_id` と `user_id`
- **出力**: メッセージ情報とユーザーの詳細情報を統合したレスポンス
- **特徴**: grpc-federationによる自動的な複数サービスの呼び出しとデータ統合

## 📋 前提条件

### 必須
- Go 1.24.3+
- Buf CLI
- Docker（コンテナ実行の場合）

### 外部サービス
BFFサーバーが正常に動作するためには、以下のサービスが起動している必要があります：
- User Service: `localhost:8081`
- Message Service: `localhost:8082`

## 🔧 セットアップ手順

### 1. リポジトリのクローン
```bash
git clone <repository-url>
cd bff
```

### 2. 依存関係のインストール
```bash
go mod download
```

### 3. Protobufコードの生成
```bash
buf generate
```

### 4. サーバーの起動
```bash
go run main.go
```

サーバーが正常に起動すると、以下のようなログが表示されます：
```json
{"time":"2025-06-21T11:30:00Z","level":"INFO","msg":"Starting BFF server","address":"0.0.0.0:8080"}
```

## 🐳 Docker使用方法

### イメージのビルド
```bash
docker build -t bff-server .
```

### コンテナの実行
```bash
docker run -p 8080:8080 bff-server
```

### ヘルスチェック
Dockerfileにはヘルスチェックが組み込まれており、`grpc-health-probe`を使用してサーバーの状態を監視します。

## 📖 API仕様

### BffService

#### GetMessage
複数のサービスからデータを統合してメッセージ情報を取得します。

**リクエスト**
```protobuf
message GetMessageRequest {
  string message_id = 1;  // 取得するメッセージのID
  string user_id = 2;     // 関連するユーザーのID
}
```

**レスポンス**
```protobuf
message GetMessageResponse {
  Message message = 1;    // 統合されたメッセージ情報
}

message Message {
  User user = 1;          // ユーザー詳細情報
  string text = 2;        // メッセージテキスト
}

message User {
  string id = 1;                              // ユーザーID
  string email = 2;                           // メールアドレス
  string name = 3;                            // ユーザー名
  google.protobuf.Timestamp created_at = 4;   // 作成日時
  google.protobuf.Timestamp updated_at = 5;   // 更新日時
}
```

## 🧪 動作確認とテスト

### grpcurlを使用したAPI動作確認

#### サービス一覧の確認
```bash
grpcurl -plaintext localhost:8080 list
```

#### BffServiceメソッドの確認
```bash
grpcurl -plaintext localhost:8080 list bff.v1.BffService
```

#### GetMessageメソッドの実行
```bash
grpcurl -plaintext -d '{"message_id": "msg123", "user_id": "user456"}' localhost:8080 bff.v1.BffService.GetMessage
```

#### ヘルスチェック
```bash
grpcurl -plaintext -d '{"service": "bff.v1.BffService"}' localhost:8080 grpc.health.v1.Health.Check
```

### 期待されるレスポンス例
```json
{
  "message": {
    "user": {
      "id": "user456",
      "email": "user@example.com",
      "name": "Example User",
      "createdAt": "2025-06-21T10:00:00Z",
      "updatedAt": "2025-06-21T10:00:00Z"
    },
    "text": "Hello, world!"
  }
}
```

## 🛠️ 技術仕様

### 使用技術スタック
- **言語**: Go 1.24.3
- **フレームワーク**: gRPC
- **データ統合**: grpc-federation
- **プロトコル**: Protocol Buffers (proto3)
- **ビルドツール**: Buf
- **コンテナ**: Docker

### ポート設定
| サービス | ポート | 説明 |
|---------|--------|------|
| BFF Server | 8080 | メインのBFFサーバー |
| User Service | 8081 | ユーザー管理サービス |
| Message Service | 8082 | メッセージ管理サービス |

### grpc-federationの特徴
- **自動サービス呼び出し**: proto定義に基づいて複数のサービスを自動的に呼び出し
- **データ統合**: 複数のレスポンスを単一のレスポンスに統合
- **エラーハンドリング**: 上流サービスのエラーを適切に処理
- **パフォーマンス**: 並列処理による効率的なデータ取得

## 📁 プロジェクト構造

```
├── proto/                    # Protocol Buffers定義
│   ├── bff/v1/              # BFFサービス定義
│   ├── user/v1/             # Userサービス定義
│   └── message/v1/          # Messageサービス定義
├── gen/                     # 生成されたgRPCコード
│   ├── bff/v1/
│   ├── user/v1/
│   └── message/v1/
├── main.go                  # メインサーバー実装
├── buf.yaml                 # Buf設定ファイル
├── buf.gen.yaml            # コード生成設定
├── go.mod                   # Go modules設定
├── go.sum                   # 依存関係ハッシュ
├── Dockerfile              # コンテナイメージ定義
└── README.md               # プロジェクトドキュメント
```

## 🐛 トラブルシューティング

### よくある問題と解決方法

#### 1. サーバー起動エラー
```
failed to connect to user service: connection refused
```
**解決方法**: User ServiceとMessage Serviceが起動していることを確認してください。

#### 2. protobuf生成エラー
```
buf generate failed
```
**解決方法**: 
```bash
buf mod update
buf generate
```

#### 3. ポート競合
```
failed to listen: address already in use
```
**解決方法**: ポート8080が使用されていないことを確認するか、`main.go`のPort定数を変更してください。

#### 4. gRPC呼び出しエラー
**確認事項**:
- 外部サービスが正常に応答しているか
- ネットワーク接続に問題がないか
- リクエストデータが正しい形式か

### ログ確認
サーバーはJSON形式でログを出力します。問題発生時は以下の情報を確認してください：
- エラーメッセージ
- 外部サービスへの接続状況
- リクエスト/レスポンス内容

## 🔄 開発とカスタマイズ

### 新しいサービスの追加

1. **proto定義の追加**
   ```protobuf
   // proto/newservice/v1/newservice.proto
   service NewService {
     rpc GetData(GetDataRequest) returns (GetDataResponse);
   }
   ```

2. **BFF proto定義の更新**
   ```protobuf
   // proto/bff/v1/bff.proto に追加
   import "newservice/v1/newservice.proto";
   
   option (grpc.federation.file) = {
     import: [
       "user/v1/user.proto",
       "message/v1/message.proto",
       "newservice/v1/newservice.proto"  // 追加
     ]
   };
   ```

3. **クライアントファクトリの拡張**
   ```go
   // main.go に追加
   func (f *BffServiceClientFactory) NewService_V1_NewServiceClient(cfg bffv1.BffServiceClientConfig) (newservicev1.NewServiceClient, error) {
     return newservicev1.NewNewServiceClient(f.newServiceConn), nil
   }
   ```

### 新しいRPCメソッドの追加

1. proto定義でメソッドを定義
2. grpc-federation設定を追加
3. `buf generate` でコード再生成
4. サーバーの再起動

### パフォーマンス最適化

- **接続プーリング**: gRPC接続の再利用
- **タイムアウト設定**: 適切なタイムアウト値の設定
- **ロードバランシング**: 複数インスタンスでの負荷分散
- **メトリクス監視**: Prometheusメトリクスの追加

## 📚 参考リンク

- [grpc-federation Documentation](https://github.com/mercari/grpc-federation)
- [gRPC-Go Documentation](https://grpc.io/docs/languages/go/)
- [Protocol Buffers Guide](https://protobuf.dev/)
- [Buf Documentation](https://buf.build/docs/)

## 📄 ライセンス

このプロジェクトのライセンス情報については、LICENSEファイルを参照してください。
