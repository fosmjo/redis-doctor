package doctor

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

type Doctor struct {
	client   redis.UniversalClient
	outputer Visitor
}

type Options struct {
	Pattern string
	Count   int
}

func New(redisOpts *redis.UniversalOptions, outputer Visitor) *Doctor {
	return &Doctor{
		client:   redis.NewUniversalClient(redisOpts),
		outputer: outputer,
	}
}

func (d *Doctor) Diagnose(ctx context.Context, symptom string, opts Options) error {
	switch symptom {
	case "slowlog":
		return d.slowlog(ctx, opts.Count)
	case "bigkey":
		return errors.New("not implemented")
	default:
		return errors.New("unknown symptom")
	}
}

func (d *Doctor) slowlog(ctx context.Context, count int) error {
	logs, err := d.client.SlowLogGet(ctx, int64(count)).Result()
	if err != nil {
		return err
	}

	for _, log := range logs {
		err := (*SlowLog)(&log).Accept(d.outputer) // nolint: gosec
		if err != nil {
			return err
		}
	}

	return nil
}
