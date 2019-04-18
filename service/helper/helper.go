package helper

import (
	"github.com/joel-tkp/coretest-go/service/redis" // Redis Service Provider
	"github.com/joel-tkp/coretest-go/pkg/log" // Logger Package
)

// Helper
func RevalidateCache(key string, serviceId string) {
	// begin invalidate cache
	redis.Delete(key)
	log.SetOutputToFile("log/" + serviceId + ".log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", key)
 	// end of cache invalidation
} 

func RevalidateCachePattern(keyPattern string, serviceId string) {
	// begin invalidate cache pattern
	relatedKeys,_ := redis.GetKeys(keyPattern)
	for _, key := range relatedKeys {
		redis.Delete(key)
	}
	log.SetOutputToFile("log/" + serviceId + ".log")
	log.SetLevel(log.InfoLevel)
	log.Infow("Cache Invalidation", "key", keyPattern)
 	// end of cache pattern invalidation
}
