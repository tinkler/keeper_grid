package automation_compatible

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tinkler/keeper_grid/pkg/eth/com"
	"github.com/tinkler/keeper_grid/pkg/pkeystore"
)

/*
A pool which eth_call the automation
*/
type pool struct {
	client *ethclient.Client
	wallet *pkeystore.KeyStore
}

type binder struct {
	stub     *pool
	contract *com.AutomationCompatible
}

func NewHandler(ctx context.Context, rawurl string, wallet *pkeystore.KeyStore) (*pool, error) {
	client, err := ethclient.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return &pool{client: client, wallet: wallet}, nil
}

func (p *pool) NewBinder(addr common.Address) (*binder, error) {
	contract, err := com.NewAutomationCompatible(addr, p.client)
	if err != nil {
		return nil, err
	}
	return &binder{
		stub:     p,
		contract: contract,
	}, nil
}

func (b *binder) CheckUpkeep(opts *bind.CallOpts, checkData []byte) (callable bool, executedata []byte, err error) {
	executedata = make([]byte, 0)
	contract := com.AutomationCompatibleRaw{Contract: b.contract}
	result := []interface{}{
		&callable,
		&executedata,
	}
	results := []interface{}{
		&result,
	}
	err = contract.Call(opts, &results, "checkUpkeep", checkData)
	return
}
