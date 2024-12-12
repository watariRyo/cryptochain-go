# Goによるブロックチェーン基盤作成

- Done
  - 簡易なBlock, Chainの実装
  - handler
  - P2P実現のためのPubSub
  - コンテナ仮想化
- TODO
  - Wallet
  - Keys、署名
  - Transaction, Transaction Pool
  - Frontend by React

## スタック

| 名前 | Version | 概要 |
| --- | --- | --- |
| Go | 1.23 | 開発言語 |
| go-chi | v5.1 | [middleware](https://github.com/go-chi/chi/v5) |
| zerolog | 1.33 | [Logger](https://github.com/rs/zerolog) |
| viper | 1.19 | [config管理](https://github.com/spf13/viper) |
| Redis | - | P2P、PubSubに使用 |
| Docker | - | コンテナ仮想化 |
| air | - | HOTリロード |

## ディレクトリ

### cmd

- エントリ

### configs

- 環境変数の設定、他Config

### internal/logger

- 共通logger

### internal/time

- タイムスタンプを扱うinterface。
- テスト用のMockと通常運用とで分けるため

### redis

- Redisの設定ファイル等

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

### web/redis

- PubSub実装

### web/server

- RESTAPIの設定

## Dependencies

server -> handler -> (UseCase/domain)？ -> infra  

redisがチェーン持っているのは避けたいが、pubsubの都合難しい、、、  
blockをDBと同等とみなして同じ階層とした  
