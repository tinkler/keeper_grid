package automation_compatible

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/hibiken/asynq"
)

func waitForReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err != nil {
			if err == ethereum.NotFound {
				if ctx.Err() == nil {
					time.Sleep(1 * time.Second)
					continue
				} else {
					return nil, ctx.Err()
				}
			}
			return nil, err
		}
		if receipt != nil && receipt.Status == types.ReceiptStatusFailed {
			return nil, fmt.Errorf("transaction %s failed, %w", txHash.Hex(), asynq.SkipRetry)
		}
		if receipt != nil && receipt.BlockNumber.Cmp(big.NewInt(0)) > 0 {
			return receipt, nil
		}
	}
}
