# Example

```go
package workerhandler

import (
	"context"
	"fmt"
	"time"

	"example.service/internal/modules/push-notif/usecase"

	"github.com/mrapry/go-lib/codebase/factory/types"
	"github.com/mrapry/go-lib/golibhelper"
	"github.com/mrapry/go-lib/logger"
	"github.com/mrapry/go-lib/tracer"
)

// CronHandler struct
type CronHandler struct {
	uc usecase.PushNotifUsecase
}

// NewCronHandler constructor
func NewCronHandler(uc usecase.PushNotifUsecase) *CronHandler {
	return &CronHandler{
		uc: uc,
	}
}

// MountHandlers return group map topic key to handler func
func (h *CronHandler) MountHandlers() map[string]types.WorkerHandlerFunc {
	return map[string]types.WorkerHandlerFunc{
		golibhelper.CronJobKeyToString("push-notif", "10s"):            h.handlePushNotif,
		golibhelper.CronJobKeyToString("heavy-push-notif", "22:43:07"): h.handleHeavyPush,
	}
}

func (h *CronHandler) handlePushNotif(ctx context.Context, message []byte) error {
	trace := tracer.StartTrace(ctx, "CronDelivery-HandlePushNotif")
	defer trace.Finish()

	logger.LogI("processing")
	logger.LogI("done")
	return nil
}

func (h *CronHandler) handleHeavyPush(ctx context.Context, message []byte) error {
	trace := tracer.StartTrace(ctx, "CronDelivery-HandleHeavyPush")
	defer trace.Finish()

	fmt.Println("processing")
	time.Sleep(30 * time.Second)
	fmt.Println("done")
	return nil
}

```