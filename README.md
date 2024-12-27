# Go によるブロックチェーン基盤作成

報酬やジェネシスブロックは適当  
いったんバックエンド完成

- Done
  - 簡易な Block, Chain, Wallet, Transaction の実装
  - Transaction Pool, Mining
  - handler
  - Mine Transactions
  - P2P 実現のための PubSub
  - コンテナ仮想化
- TODO
  - Validate
    - Blockchain.isValidChain
      - ハッシュが難易度要件を満たしている
    - Wallet.createTransactionTransaction.verifyTransaction
      - 負の金額が不可能
  - Frontend by React
  - Walletsパッケージの歪な依存改善または整理

## スタック

| 名前    | Version | 概要                                           |
| ------- | ------- | ---------------------------------------------- |
| Go      | 1.23    | 開発言語                                       |
| go-chi  | v5.1    | [middleware](https://github.com/go-chi/chi/v5) |
| zerolog | 1.33    | [Logger](https://github.com/rs/zerolog)        |
| mock    | 0.5.0   | [Mock](https://github.com/uber-go/mock)
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

### web/infra/wallets

- WalletとTransactionの実装

### web/handler

- API 受け口, json の marshal/unmarshal

### web/server

- RESTAPI の設定

### web/usecase

- ブロック等に直接触れないロジック部

## Dependencies

server -> handler -> usecase -> repository -> infra

WalletとTransactionが密結合でimport cycle errorが発生する  
それらをラップする構造体を作成し強引に依存関係を保っている  
→ API実行でサンプルだから簡便なレイヤードにしたが、失敗したかもしれない  
infra層間ではblockが一番下になっている  
