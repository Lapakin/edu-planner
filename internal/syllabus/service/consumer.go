package service

import (
	"context"

	"github.com/Lapakin/edu-planner/internal/adapter/broker"
	"github.com/Lapakin/edu-planner/internal/domain"
	"github.com/Lapakin/edu-planner/internal/logging"
	"github.com/Lapakin/edu-planner/internal/queue"
)

func (s *Services) StartMassageConsumer(ctx context.Context, logger *logging.Logger, sqsCfg *broker.Config) error {
	return queue.StartMassageConsumer(ctx, logger, sqsCfg, func(ctx context.Context, massage *domain.Massage) error {
		if massage.ObjectType == domain.ObjectTypeUser {
			return s.UserSvc.MassageConsumer(ctx, massage)
		}
		return nil
	})
}
