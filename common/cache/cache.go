package cache

import (
	"context"
	"encoding"
	"errors"
	"fmt"
	"time"

	aero "github.com/aerospike/aerospike-client-go/v7"
	aerolog "github.com/aerospike/aerospike-client-go/v7/logger"
	"k8s.io/klog/v2"
)

// AeroRepo represnts internal struct that
// holds necessary component for Aerospike's operations
type AeroRepo[T any] struct {
	client *aero.Client
	ns     string
	set    string
	bin    string
}

// NewAero initializes a new Aerospike client for a given type T
func NewAero[T any](config *Config) (*AeroRepo[T], error) {
	if config.LogLevel != "" {
		switch {
		case config.LogLevel == "Info":
			aerolog.Logger.SetLevel(aerolog.INFO)
		case config.LogLevel == "Error":
			aerolog.Logger.SetLevel(aerolog.ERR)
		case config.LogLevel == "Warn":
			aerolog.Logger.SetLevel(aerolog.WARNING)
		case config.LogLevel == "Debug":
			aerolog.Logger.SetLevel(aerolog.DEBUG)
		default:
			return nil, errors.New("not a supported aero log level")
		}
	}

	client, err := aero.NewClient(config.URL, config.Port)
	if err != nil {
		return nil, err
	}

	return &AeroRepo[T]{
		client: client,
		ns:     config.Namespace,
		set:    config.Set,
		bin:    config.Bin,
	}, nil
}

// Set implments `Set()` method from cache's interface
func (r *AeroRepo[T]) Set(ctx context.Context, value *T, exp time.Duration) error {
	key, err := aero.NewKey(r.ns, r.set, value)
	if err != nil {
		klog.ErrorS(err, "cache: aero set error")
		return err
	}

	bins := aero.BinMap{r.bin: value}
	policy := aero.NewWritePolicy(0, uint32(exp.Seconds()))
	return r.client.Put(policy, key, bins)
}

// Get implements `Get()` method from cache's interface
func (r *AeroRepo[T]) Get(ctx context.Context, id string) (T, error) {
	var result T
	key, err := aero.NewKey(r.ns, r.bin, id)
	if err != nil {
		klog.ErrorS(err, "cache: aero get error")
		return result, err
	}

	getPolicy := aero.NewPolicy()
	record, err := r.client.Get(getPolicy, key, r.bin)
	if err != nil {
		if errors.Is(err, aero.ErrKeyNotFound) {
			klog.ErrorS(err, "cache: key not found")
			return result, err
		}
		klog.ErrorS(err, "cache: aero get error")
		return result, err
	}

	if record == nil {
		err := fmt.Errorf("cache: aero record is empty")
		klog.ErrorS(err, "cache: aero record is empty")
		return result, err
	}

	value, ok := record.Bins[r.bin]
	if !ok {
		err := fmt.Errorf("bin %s not found in record", r.bin)
		klog.ErrorS(err, "cache: bin not found in record")
		return result, err
	}

	var deserializationErr error
	result, deserializationErr = deserializeRecord[T](value)
	if deserializationErr != nil {
		klog.ErrorS(deserializationErr, "cache: deserialization")
		return result, deserializationErr
	}

	return result, nil
}

// deserializeRecord will deserialize an Aerospike's record to a type T (type T should implment `encoding` interface)
func deserializeRecord[T any](value any) (T, error) {
	var result T

	switch v := value.(type) {
	case []byte:
		if unmarshaler, ok := any(&result).(encoding.BinaryUnmarshaler); ok {
			err := unmarshaler.UnmarshalBinary(v)
			if err != nil {
				return result, err
			}
		} else {
			return result, fmt.Errorf("type %T does not implement encoding.BinaryUnmarshaler", &result)
		}
	case string:
		if unmarshaler, ok := any(&result).(encoding.TextUnmarshaler); ok {
			err := unmarshaler.UnmarshalText([]byte(v))
			if err != nil {
				return result, err
			}
		} else {
			return result, fmt.Errorf("type %T does not implement encoding.TextUnmarshaler", &result)
		}
	default:
		return result, fmt.Errorf("value type %T is not supported", value)
	}

	return result, nil
}
