package usecase

import (
	"fmt"
	"io"
	"net/http"
)

func (u *UseCase) SyncChain() error {
	// Route Node
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api/blocks", u.configs.Host), nil)
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Request Error. StatusCode: %d", response.StatusCode)
	}

	payload, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	u.repo.BlockChain.UnmarshalAndReplaceBlock(payload, u.timeProvider)

	return nil
}