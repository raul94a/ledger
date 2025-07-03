package appRedis

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)


var redisClient *redis.Client

/// .env
func Get() *redis.Client {    
    if redisClient != nil {
        return redisClient
	}
    
    godotenv.Load()
    addr := os.Getenv("REDIS_HOST")
    pwd := os.Getenv("REDIS_PASSWORD")
    port := os.Getenv("REDIS_PORT")
    redisClient = redis.NewClient(&redis.Options{
        Addr:	  fmt.Sprintf("%s:%s",addr,port),
        Password: pwd, // No password set
        DB:		  0,  // Use default DB
        Protocol: 2,  // Connection protocol
    })
	
	
	return redisClient
}