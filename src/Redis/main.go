package Redis

import (
	"context"
	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	LockKey             = "counter_lock"
	SentenceCounterKey  = "sentence_counter"
	StoryCounterKey     = "story_counter"
	ParagraphCounterKey = "paragraph_counter"
	WordCounterKey      = "word_counter"
)

var (
	Nil           = redis.Nil
	initialised   uint32
	mutex         sync.Mutex
	RedisInstance *redis.Client
	expiryTime    = time.Second * 10
	sleepDuration = time.Millisecond * 100
	ctx           = context.Background()
)

func InitializeRedis() *redis.Client {
	//using check if thread check

	//check
	if atomic.LoadUint32(&initialised) == 1 {
		return RedisInstance
	}
	//lock
	mutex.Lock()
	defer mutex.Unlock()

	//check
	if initialised == 0 {
		redisOption := redis.Options{}

		redisOption.PoolSize = 1000

		log.Info("delivery redis address is ", "localhost:6379")
		log.Info("delivery redis password is ")
		redisOption.PoolTimeout = time.Duration(1000) //1000ms
		redisOption.Addr = "localhost:6379"
		redisOption.ReadTimeout = 200 * time.Millisecond
		redisOption.Password = ""
		RedisInstance = redis.NewClient(&redisOption)

		pong, err := RedisInstance.Ping(ctx).Result()
		log.Info(pong, " err is ", err)

		//mark initialised
		atomic.StoreUint32(&initialised, 1)
	}
	return RedisInstance
}

func IncrementStoryCounter() int64 {
	count, err := RedisInstance.Incr(ctx, StoryCounterKey).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"key": SentenceCounterKey,
		}).Error("error incrementing story counter")
	} else {
		log.Debug("current story count is ", count)
	}
	return count
}

func IncrementParagraphCounter() int64 {
	count, err := RedisInstance.Incr(ctx, ParagraphCounterKey).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"key": ParagraphCounterKey,
		}).Error("error incrementing paragraph counter")
	} else {
		log.Debug("current paragraph count is ", count)
	}
	return count
}

func IncrementSentenceCounter() int64 {
	count, err := RedisInstance.Incr(ctx, SentenceCounterKey).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"key": SentenceCounterKey,
		}).Error("error incrementing sentence counter")
	} else {
		log.Debug("current sentence count is ", count)
	}
	return count
}

func IncrementWordCounter() int64 {
	count, err := RedisInstance.Incr(ctx, WordCounterKey).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
			"key": WordCounterKey,
		}).Error("error incrementing word counter")
	} else {
		log.Debug("current word count is ", count)
	}
	return count
}

func LoadCount(key string) int64 {
	count, err := RedisInstance.Get(ctx, key).Result()
	if err != nil {
		if err != Nil {
			log.WithFields(log.Fields{
				"err": err,
				"key": key,
			}).Error("error loading count")
		}
	}
	cnt, _ := strconv.ParseInt(count, 10, 64)
	return cnt
}

func SetCounterCount(key, val string) bool {
	_, err := RedisInstance.Set(ctx, key, val, 0).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"err":   err,
			"key":   key,
			"value": val,
		}).Error("err setting counter")
		return false
	}
	return true
}

func AcquireLock(key string) bool {
	log.Debug("Getting a Lock")
	result, err := RedisInstance.SetNX(ctx, key, "1", expiryTime).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("while acquiring redis lock")
		return false
	} else {
		log.Debug("successfully acquired a lock")
	}
	return result
}

func ReleaseLock(key string) bool {
	val, err := RedisInstance.Del(ctx, key).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("while releasing redis lock")
	} else if val == 1 {
		return true
	}
	return false
}

func ExecuteQueryLock(key string) bool {
	var (
		cnt = 0
	)
	// Loop till key does not exists
	for cnt < 5 && CheckLockKeyExistance(key) {
		//sleep for 1 sec
		time.Sleep(sleepDuration)
		cnt++
		log.Debug("Sleeping for 1s waiting for acquiring lock for key = ", key)
	}
	if cnt < 5 {
		AcquireLock(key)
	}
	return true
	//defer func() {
	//	// Keep in mind lock auto expires in 1 minute
	//	if success && key != "" {
	//		log.Debug("Releasing lock for ", key, " success = ", releaseLock(key))
	//	}
	//}()
}

func CheckLockKeyExistance(key string) bool {
	log.Debug("Checking for key ", key)
	val, err := RedisInstance.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("while checking lock existence")
		return false

	} else {
		if val == "1" {
			return true
		} else {
			return false
		}
	}
}

func CheckKeyExistance(key string) bool {
	log.Debug("Checking for key ", key)
	_, err := RedisInstance.Get(ctx, key).Result()
	if err == redis.Nil {
		return false
	} else if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("while checking key existence")
		return false

	} else {
		return true
	}
}
