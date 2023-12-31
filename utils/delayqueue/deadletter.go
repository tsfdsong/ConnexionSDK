package delayqueue

import (
	"errors"
	"fmt"
	"github/Connector-Gamefi/ConnectorGoSDK/utils/logger"

	go_redis "github.com/go-redis/redis/v8"
)

const (
	luaRespawnDeadletterScript = `
local deadletter = KEYS[1]
local queue = KEYS[2]
local poolPrefix = KEYS[3]
local limit = tonumber(ARGV[1])
local respawnTTL = tonumber(ARGV[2])

for i = 1, limit do
    local jobID = redis.call("RPOPLPUSH", deadletter, queue)
	if jobID == false then
		return i - 1  -- deadletter is empty
	end
    -- unpack the jobID, and set the TTL
	redis.call("HSET", poolPrefix .. "/" .. jobID, "tries", "1")
    if respawnTTL > 0 then
		redis.call("EXPIRE", poolPrefix .. "/" .. jobID, respawnTTL)
	end
end
return limit  -- deadletter has more data when return value is >= limit
`
	luaDeleteDeadletterScript = `
local deadletter = KEYS[1]
local poolPrefix = KEYS[2]
local limit = tonumber(ARGV[1])

for i = 1, limit do
	local jobID = redis.call("RPOP", deadletter)
	if jobID == false then
		return i - 1
	end
	-- delete the job from the job pool
	redis.call("DEL", poolPrefix .. "/" .. jobID)
end
return limit
`
)

var (
	respawnDeadletterSHA string
	deleteDeadletterSHA  string
)

// Because the DeadLetter is not like Timer which is a singleton,
// DeadLetters are transient objects like Queue. So we have to preload
// the lua scripts separately.
func PreloadDeadLetterLuaScript(redis *RedisInstance) error {
	sha, err := redis.Conn.ScriptLoad(dummyCtx, luaRespawnDeadletterScript).Result()
	if err != nil {
		return fmt.Errorf("failed to preload lua script: %s", err)
	}
	respawnDeadletterSHA = sha

	sha, err = redis.Conn.ScriptLoad(dummyCtx, luaDeleteDeadletterScript).Result()
	if err != nil {
		return fmt.Errorf("failed to preload lua script: %s", err)
	}
	deleteDeadletterSHA = sha
	return nil
}

// DeadLetter is where dead job will be buried, the job can be respawned into ready queue
type DeadLetter struct {
	redis *RedisInstance
	queue queue
}

// NewDeadLetter return an instance of DeadLetter storage
func NewDeadLetter(ns, q string, redis *RedisInstance) (*DeadLetter, error) {
	dl := &DeadLetter{
		redis: redis,
		queue: queue{namespace: ns, queue: q},
	}
	if respawnDeadletterSHA == "" || deleteDeadletterSHA == "" {
		return nil, errors.New("dead letter's lua script is not preloaded")
	}
	return dl, nil
}

func (dl *DeadLetter) Name() string {
	return dl.queue.DeadletterString()
}

// NOTE: this method is not called any where except in tests, but this logic is
// implement in the timer's pump script. please refer to that.
func (dl *DeadLetter) Add(jobID string) error {
	if err := dl.redis.Conn.Persist(dummyCtx, PoolJobKey2(dl.queue, jobID)).Err(); err != nil {
		return err
	}
	return dl.redis.Conn.LPush(dummyCtx, dl.Name(), jobID).Err()
}

func (dl *DeadLetter) Peek() (size int64, jobID string, err error) {
	jobID, err = dl.redis.Conn.LIndex(dummyCtx, dl.Name(), -1).Result()
	switch err {
	case nil:
		// continue
	case go_redis.Nil:
		return 0, "", ErrNotFound
	default:
		return 0, "", err
	}
	size, err = dl.Size()
	if err != nil {
		return 0, "", err
	}
	return size, jobID, nil
}

func (dl *DeadLetter) Delete(limit int64) (count int64, err error) {
	if limit <= 0 {
		return 0, nil
	}
	// Note: we can also use rpop+del to delete deadletter when limit == 1, but may cause some data lose control.
	// Delete is a rarely method, use lua script to handle it maybe just so so.
	poolPrefix := dl.queue.PoolPrefixString()
	var batchSize = BatchSize
	if batchSize > limit {
		batchSize = limit
	}
	for {
		val, err := dl.redis.Conn.EvalSha(dummyCtx, deleteDeadletterSHA, []string{dl.Name(), poolPrefix}, batchSize).Result()
		if err != nil {
			if isLuaScriptGone(err) {
				if err := PreloadDeadLetterLuaScript(dl.redis); err != nil {
					logger.Logrus.WithError(err).Error("Failed to load deadletter lua script")
					// endless-loop if not return
					return count, err
				}
				continue
			}
			return count, err
		}
		n, _ := val.(int64)
		count += n
		if n < batchSize { // Dead letter is empty
			break
		}
		if count >= limit {
			break
		}
		if limit-count < batchSize {
			batchSize = limit - count // This is the last batch, we should't respawn jobs exceeding the limit.
		}
	}
	return count, nil
}

func (dl *DeadLetter) Respawn(limit, ttlSecond int64) (count int64, err error) {
	if limit <= 0 {
		return 0, nil
	}
	// Note: we can also use rpoplpush+hset+expire to respawn deadletter when limit == 1, but may cause some data lose control.
	// Respawn is a rarely method, use lua script to handle it maybe just so so.
	defer func() {
		if err != nil && count != 0 {
			metrics.deadletterRespawnJobs.WithLabelValues(dl.redis.Name).Add(float64(count))
		}
	}()
	queueName := dl.queue.ReadyQueueString()
	poolPrefix := dl.queue.PoolPrefixString()
	var batchSize = BatchSize
	if batchSize > limit {
		batchSize = limit
	}
	for {
		val, err := dl.redis.Conn.EvalSha(dummyCtx, respawnDeadletterSHA, []string{dl.Name(), queueName, poolPrefix}, batchSize, ttlSecond).Result() // Respawn `batchSize` jobs at a time
		if err != nil {
			if isLuaScriptGone(err) {
				if err := PreloadDeadLetterLuaScript(dl.redis); err != nil {
					logger.Logrus.WithError(err).Error("Failed to load deadletter lua script")
					// endless-loop if not return
					return count, err
				}
				continue
			}
			return count, err
		}
		n, _ := val.(int64)
		count += n
		if n < batchSize { // Dead letter is empty
			break
		}
		if count >= limit {
			break
		}
		if limit-count < batchSize {
			batchSize = limit - count // This is the last batch, we should't respawn jobs exceeding the limit.
		}
	}
	return count, nil
}

func (dl *DeadLetter) Size() (size int64, err error) {
	return dl.redis.Conn.LLen(dummyCtx, dl.Name()).Result()
}
