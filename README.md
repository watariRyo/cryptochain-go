# Go によるブロックチェーン基盤作成

- Done
  - 簡易な Block, Chain, Wallet, Transaction の実装
  - handler
  - P2P 実現のための PubSub
  - コンテナ仮想化
- TODO
  - Transaction Pool
    - Wallet, Transactionの依存修正
  - Mine Transactions
  - Frontend by React

## スタック

| 名前    | Version | 概要                                           |
| ------- | ------- | ---------------------------------------------- |
| Go      | 1.23    | 開発言語                                       |
| go-chi  | v5.1    | [middleware](https://github.com/go-chi/chi/v5) |
| zerolog | 1.33    | [Logger](https://github.com/rs/zerolog)        |
| viper   | 1.19    | [config 管理](https://github.com/spf13/viper)  |
| Redis   | -       | P2P、PubSub に使用                             |
| Docker  | -       | コンテナ仮想化                                 |
| air     | -       | HOT リロード                                   |

## ディレクトリ

### cmd

- エントリ

### configs

- 環境変数の設定、他 Config

### internal/logger

- 共通 logger

### internal/time

- タイムスタンプを扱う interface。
- テスト用の Mock と通常運用とで分けるため

### redis

- Redis の設定ファイル等

### web/domain/model

- メソッド持たない構造体定義

### web/domain/repository

- infra 層の interface

### web/infra/blcok

- ブロックおよびチェーンの作成、追加、検証の実装を格納
  - block  
    ブロックの定義、作成、検証、マイニング調整の実施
  - block_chain  
    チェーンの管理、検証、チェーンへのブロック追加
  - crypto_hash  
    ハッシュ作成
  - genesis  
    ジェネシスブロックの作成

### web/infra/redis

- PubSub 実装

### web/handler

- API 受け口, json の marshal/unmarshal

### web/server

- RESTAPI の設定

### web/usecase

- ブロック等に直接触れないロジック部

## Dependencies

server -> handler -> usecase -> repository -> infra

redis がチェーン持っているのは避けたいが、pubsub の都合難しい、、、  
block を DB と同等とみなして infra 層とした
