package Constants

import (
	"Cassandra"
	"Redis"
	"os"
	"strconv"
)

const (
	ASC           = "asc"
	DESC          = "desc"
	SortCreatedAt = "created_at"
	SortUpdatedAt = "updated_at"
	SortTitle     = "title"
)

var (
	CQL      *Cassandra.Cassandra
	Debug, _ = strconv.ParseBool(os.Getenv("VERLOOP_DEBUG"))
)

func init() {
	Redis.RedisInstance = Redis.InitializeRedis()
	CQL = Cassandra.Instance()
	//CurrentStoryTitle = make([]string,2,2)
	SetInitialCount()

}
