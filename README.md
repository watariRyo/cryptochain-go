# Goによるブロックチェーン基盤作成

## スタック

| 名前 | Version | 概要 |
| --- | --- | --- |
| Go | 1.23 | 開発言語 |
| go-chi | v5.1 | [middleware](github.com/go-chi/chi/v5) |
| zerolog | 1.33 | [Logger](github.com/rs/zerolog) |

## ディレクトリ

### internal/logger

- 共通logger

### internal/time

- タイムスタンプを扱うinterface。
- テスト用のMockと通常運用とで分けるため

### web/blcok

- ブロックおよびチェーンの作成、追加、検証の実装を格納
  - block  
  ブロックの定義、作成、検証、マイニング調整の実施
  - block_chain  
  チェーンの管理、検証、チェーンへのブロック追加
  - crypto_hash  
  ハッシュ作成
  - genesis  
  ジェネシスブロックの作成

### web/handler

- API受け口, 実装別だし未定

### web/server

- RESTAPIの設定
