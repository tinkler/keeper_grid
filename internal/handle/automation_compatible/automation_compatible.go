package automation_compatible

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/hibiken/asynq"
	"github.com/tinkler/keeper_grid/pkg/tasks/automation_compatible"
	"github.com/tinkler/moonmist/pkg/mlog"
)

func (h *pool) HandleCheckUpkeep(ctx context.Context, t *asynq.Task) error {
	p, err := automation_compatible.ParseFrom(t)
	if err != nil {
		return err
	}
	binder, err := h.NewBinder(p.AutomationCompatibleAddress)
	if err != nil {
		return fmt.Errorf("bind to contract AutomationCompatible failed: %v: %w", err, asynq.SkipRetry)
	}
	if !h.wallet.HasAddress(p.Keeper) {
		return fmt.Errorf("no permission to keeper %s, %w", p.Keeper.Hex(), asynq.SkipRetry)
	}
	callable, executedata, err := binder.CheckUpkeep(&bind.CallOpts{
		From: p.Keeper,
	}, []byte{})
	if err != nil {
		if errors.Is(err, bind.ErrNoCode) {
			return fmt.Errorf("contract is not exist %s, %w", p.AutomationCompatibleAddress.Hex(), asynq.SkipRetry)
		}
		return err
	}
	mlog.Info(callable, executedata)
	return nil
}
