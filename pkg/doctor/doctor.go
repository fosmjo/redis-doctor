package doctor

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/fosmjo/redis-doctor/pkg/proto"
)

type Doctor struct {
	client   redis.UniversalClient
	outputer Visitor
}

type Options struct {
	Pattern     string
	Type        string
	Length      int64
	Cardinality int64
	Batch       int
	Limit       int
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
		return d.slowlog(ctx, opts.Limit)
	case "bigkey":
		return d.bigkey(ctx, opts)
	default:
		return errors.New("unknown symptom")
	}
}

func (d *Doctor) slowlog(ctx context.Context, limit int) error {
	logs, err := d.client.SlowLogGet(ctx, int64(limit)).Result()
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

func (d *Doctor) bigkey(ctx context.Context, opts Options) error {
	keys := make([]string, 0, opts.Batch)
	types := make([]string, 0, opts.Batch)
	cards := make([]int64, 0, opts.Batch)
	results := make([]*proto.DebugObjectResult, 0, opts.Batch)

	iterator := d.iterator(ctx, opts.Pattern, opts.Batch, opts.Type)

	for count := 0; count < opts.Limit; {
		keys, err := d.keys(ctx, iterator, opts.Batch, keys[:0])
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			break
		}
		if len(keys) > opts.Limit-count {
			keys = keys[:opts.Limit-count]
		}

		types, err = d.types(ctx, keys, opts.Type, types[:0])
		if err != nil {
			return err
		}

		cards, err = d.card(ctx, keys, types, cards[:0])
		if err != nil {
			return err
		}

		results, err = d.debugObject(ctx, keys, results[:0])
		if err != nil {
			return err
		}

		for i := range keys {
			_isBigKey := func() bool {
				if opts.Length != 0 {
					return isBigKey(results[0].SerializedLength, opts.Length)
				}
				return isBigKey(cards[i], opts.Cardinality)
			}()

			if !_isBigKey {
				continue
			}

			bk := &BigKey{
				Key:              keys[i],
				Type:             types[i],
				Encoding:         results[i].Encoding,
				SerializedLength: results[i].SerializedLength,
				Cardinality:      cards[i],
			}
			err = d.outputer.VisitBigKey(bk)
			if err != nil {
				return err
			}

			count++
		}
	}

	return nil
}

func (d *Doctor) iterator(
	ctx context.Context, pattern string, count int, _type string,
) *redis.ScanIterator {
	if _type == "" {
		return d.client.Scan(ctx, 0, pattern, int64(count)).Iterator()
	} else {
		return d.client.ScanType(ctx, 0, pattern, int64(count), _type).Iterator()
	}
}

func (d *Doctor) keys(
	ctx context.Context, iterator *redis.ScanIterator, batch int, keys []string,
) ([]string, error) {
	for iterator.Next(ctx) {
		keys = append(keys, iterator.Val())

		if len(keys) == batch {
			break
		}
	}

	return keys, iterator.Err()
}

func (d *Doctor) types(
	ctx context.Context, keys []string, _type string, types []string,
) ([]string, error) {
	if _type != "" {
		for range keys {
			types = append(types, _type)
		}
		return types, nil
	}

	return d._type(ctx, keys, types)
}

func (d *Doctor) _type(ctx context.Context, keys []string, types []string) ([]string, error) {
	cmds, err := d.client.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			for _, key := range keys {
				pipe.Type(ctx, key)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return nil, err
		}

		types = append(types, cmd.(*redis.StatusCmd).Val()) // nolint: forcetypeassert
	}

	return types, nil
}

func (d *Doctor) card(
	ctx context.Context, keys []string, types []string, cards []int64,
) ([]int64, error) {
	cmds, err := d.client.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			for i, key := range keys {
				switch types[i] {
				case "string":
					pipe.StrLen(ctx, key)
				case "list":
					pipe.LLen(ctx, key)
				case "hash":
					pipe.HLen(ctx, key)
				case "set":
					pipe.SCard(ctx, key)
				case "zset":
					pipe.ZCard(ctx, key)
				default:
					return fmt.Errorf("unsupported redis data type: %s", types[i])
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return nil, err
		}

		cards = append(cards, cmd.(*redis.IntCmd).Val()) // nolint: forcetypeassert
	}

	return cards, nil
}

func (d *Doctor) debugObject(
	ctx context.Context, keys []string, results []*proto.DebugObjectResult,
) ([]*proto.DebugObjectResult, error) {
	cmds, err := d.client.Pipelined(
		ctx,
		func(pipe redis.Pipeliner) error {
			for _, key := range keys {
				pipe.DebugObject(ctx, key)
			}
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	for _, cmd := range cmds {
		if err := cmd.Err(); err != nil {
			return nil, err
		}

		raw := cmd.(*redis.StringCmd).Val() // nolint: forcetypeassert
		result, err := proto.ParseDebugObjectResult(raw)
		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func isBigKey(value, threshold int64) bool {
	return value >= threshold
}
