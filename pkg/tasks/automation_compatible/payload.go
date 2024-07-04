package automation_compatible

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/hibiken/asynq"
	"github.com/tinkler/keeper_grid/pkg/tasks"
	"github.com/tinkler/moonmist/pkg/jsonz/cjson"
	"github.com/tinkler/moonmist/pkg/mlog"
)

type automationCompatiblePayload struct {
	AutomationCompatibleAddress common.Address
	Keeper                      common.Address
	CheckData                   []byte
}

func ParseFrom(t *asynq.Task) (*automationCompatiblePayload, error) {
	var p automationCompatiblePayload
	if err := cjson.Unmarshal(t.Payload(), &p); err != nil {
		return nil, fmt.Errorf("cjson.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	return &p, nil
}

func NewTask(addr string, keeper string) (*asynq.Task, error) {
	if !common.IsHexAddress(addr) {
		return nil, fmt.Errorf("the address of AutomationCompatible is invalid: %s", addr)
	}
	if !common.IsHexAddress(keeper) {
		return nil, fmt.Errorf("the address of keeper is invalid: %s", keeper)
	}
	p, err := cjson.Marshal(automationCompatiblePayload{
		AutomationCompatibleAddress: common.HexToAddress(addr),
		Keeper:                      common.HexToAddress(keeper),
	})
	if err != nil {
		mlog.Error(err)
		return nil, err
	}
	return asynq.NewTask(tasks.ACT_IOTEX, p, asynq.MaxRetry(5)), nil
}
