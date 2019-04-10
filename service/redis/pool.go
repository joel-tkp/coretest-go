package redis

import (
    "io/ioutil"
    "gopkg.in/yaml.v2" 
    "github.com/gomodule/redigo/redis"
    "os"
    "os/signal"
    "syscall"
    "time"
)

var Config struct {
    Redis    RedisConfig    `yaml:"redis"`
}

type RedisConfig struct {
    Address string `yaml:"address"`
}

var (
    Pool *redis.Pool
)

func init() {
    // read config from config directory
    out, err := ioutil.ReadFile("config/coretest.config.yml")
    if err != nil {
        panic(err)
    }
    if err := yaml.Unmarshal(out, &Config); err != nil {
        panic(err)
    }
    redisHost := Config.Redis.Address
    if redisHost == "" {
        redisHost = "127.0.0.1:6379" // default
    }
    Pool = newPool(redisHost)
    cleanupHook()
}

func newPool(server string) *redis.Pool {

    return &redis.Pool{

        MaxIdle:     3,
        IdleTimeout: 240 * time.Second,

        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", server)
            if err != nil {
                return nil, err
            }
            return c, err
        },

        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            _, err := c.Do("PING")
            return err
        },
    }
}

func cleanupHook() {

    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)
    signal.Notify(c, syscall.SIGKILL)
    go func() {
        <-c
        Pool.Close()
        os.Exit(0)
    }()
}