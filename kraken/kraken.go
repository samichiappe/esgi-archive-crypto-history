package kraken

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
)

type AssetPair struct {
	Altname string `json:"altname"`
	Wsname  string `json:"wsname"`
	// Vous pouvez ajouter d'autres champs si nécessaire
}

type KrakenResponse struct {
	Error  []string             `json:"error"`
	Result map[string]AssetPair `json:"result"`
}

func GetAssetPairs() (map[string]AssetPair, error) {
	resp, err := http.Get("https://api.kraken.com/0/public/AssetPairs")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var kr KrakenResponse
	if err := json.Unmarshal(body, &kr); err != nil {
		return nil, err
	}
	if len(kr.Error) > 0 {
		return nil, fmt.Errorf("erreur API Kraken: %v", kr.Error)
	}
	return kr.Result, nil
}

type ServerStatus struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	// D'autres champs peuvent être ajoutés si nécessaire
}

func GetServerStatus() (*ServerStatus, error) {
	resp, err := http.Get("https://api.kraken.com/0/public/SystemStatus")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Error  []string     `json:"error"`
		Result ServerStatus `json:"result"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Error) > 0 {
		return nil, fmt.Errorf("erreur API Kraken (SystemStatus): %v", result.Error)
	}
	return &result.Result, nil
}

type DataAggregate struct {
	Pairs  map[string]AssetPair
	Status *ServerStatus
}

func FetchAllDataConcurrently() (*DataAggregate, error) {
	type result struct {
		pairs  map[string]AssetPair
		status *ServerStatus
		err    error
	}

	resultsCh := make(chan result, 2)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		pairs, err := GetAssetPairs()
		resultsCh <- result{pairs: pairs, err: err}
	}()

	go func() {
		defer wg.Done()
		status, err := GetServerStatus()
		resultsCh <- result{status: status, err: err}
	}()

	wg.Wait()
	close(resultsCh)

	var aggregate DataAggregate
	for res := range resultsCh {
		if res.err != nil {
			return nil, res.err
		}
		if res.pairs != nil {
			aggregate.Pairs = res.pairs
		}
		if res.status != nil {
			aggregate.Status = res.status
		}
	}

	return &aggregate, nil
}
