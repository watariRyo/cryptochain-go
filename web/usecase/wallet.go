package usecase

import "github.com/watariRyo/cryptochain-go/web/domain/model"

func (u *UseCase) GetWalletInfo() (*model.WalletInfo, error) {
	publicKey := u.repo.Wallets.GetWallet().PublicKey
	balance, err := u.repo.Wallets.CaluculateBalance(u.GetBlock(), publicKey)
	if err != nil {
		return nil, err
	}
	return &model.WalletInfo{
		Address: publicKey,
		Balance: balance,
	}, nil
}
