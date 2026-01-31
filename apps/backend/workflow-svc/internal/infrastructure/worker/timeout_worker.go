package worker

import (
	"context"
	"time"

	"workflow-svc/internal/application/service"
	"workflow-svc/internal/domain/repository"

	"go.uber.org/zap"
)

type TimeoutWorker struct {
	repo        repository.InstanceRepository
	workflowSvc *service.WorkflowService
	interval    time.Duration
	batchSize   int
	logger      *zap.Logger
}

func NewTimeoutWorker(
	repo repository.InstanceRepository,
	workflowSvc *service.WorkflowService,
	interval time.Duration,
	batchSize int,
	logger *zap.Logger,
) *TimeoutWorker {
	return &TimeoutWorker{
		repo:        repo,
		workflowSvc: workflowSvc,
		interval:    interval,
		batchSize:   batchSize,
		logger:      logger,
	}
}

func (w *TimeoutWorker) Start(ctx context.Context) error {
	w.logger.Info("Starting timeout worker",
		zap.Duration("interval", w.interval),
		zap.Int("batch_size", w.batchSize))

	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("Stopping timeout worker")
			return ctx.Err()
		case <-ticker.C:
			if err := w.processTimeouts(ctx); err != nil {
				w.logger.Error("Failed to process timeouts", zap.Error(err))
			}
		}
	}
}

func (w *TimeoutWorker) processTimeouts(ctx context.Context) error {
	instances, err := w.repo.FindPendingTimeouts(ctx, w.batchSize)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return nil
	}

	w.logger.Info("Processing timeout instances", zap.Int("count", len(instances)))

	for _, instance := range instances {
		if err := w.processTimeout(ctx, instance.ID); err != nil {
			w.logger.Error("Failed to process timeout for instance",
				zap.String("instance_id", instance.ID),
				zap.String("workflow", instance.WorkflowName),
				zap.Error(err))
			continue
		}

		w.logger.Info("Processed timeout for instance",
			zap.String("instance_id", instance.ID),
			zap.String("workflow", instance.WorkflowName),
			zap.String("current_state", instance.CurrentState))
	}

	return nil
}

func (w *TimeoutWorker) processTimeout(ctx context.Context, instanceID string) error {
	// Process timeout event
	_, err := w.workflowSvc.ProcessEvent(ctx, instanceID, "timeout", nil)
	return err
}
