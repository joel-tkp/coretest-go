package user

import (
	"io/ioutil"
	"gopkg.in/yaml.v2" 
	_ "github.com/lib/pq" // backend-db driver
	"github.com/jmoiron/sqlx" // backend-db wrapper extension
	"References/coretest/service/user"
	"References/coretest/service/redis"
	"testing"
	"reflect"
	"fmt"
)

// config
var Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

type DatabaseConfig struct {
	Master   string `yaml:"master"`
	Follower string `yaml:"follower"`
}

type RedisConfig struct {
	Address string `yaml:"address"`
}

var userres *Resource
var userIdTested int64
var userTested user.User

func init() {
	out, err := ioutil.ReadFile("../../config/coretest.config.yml")
	if err != nil {
		fmt.Println("Init , error initiation when read config file!")
	}
	if err := yaml.Unmarshal(out, &Config); err != nil {
		fmt.Println("Init , error initiation when unmarshal config file!")
	}
	masterDB, err := sqlx.Connect("postgres", Config.Database.Master)
	if err != nil {
		fmt.Println("Init , error initiation when connect to master database!")
	}
	followerDB, err := sqlx.Connect("postgres", Config.Database.Follower)
	if err != nil {
		fmt.Println("Init , error initiation when connect to follower database!")
	}

	// Init Redis Cache Service
	redis.InitService(Config.Redis.Address)

	// user-resource init
	userres = New(masterDB, followerDB)
}

// Test Create Data
func TestCreate(t *testing.T) {
	u := user.User{
		Name: "Test", 
		Email: "test@email.com", 
		IsActive: true, 
		IdempotencyKey: "anyIdempotencyKey",
	}
	res := userres.Create(u)
	if !reflect.DeepEqual(res, u) {
		t.Error("TestCreate , expected: ", u, " ,got: ", res)
	}
	userIdTested = 1
	userTested = res
}

// Test Get Data by ID
func TestGet(t *testing.T) {
	res,err := userres.Get(userIdTested)
	if err != nil {
		t.Error("TestGet , error occured!")
	}
	if !reflect.DeepEqual(res, userTested) {
		t.Error("TestGet , expected: ", userTested, " ,got: ", res)
	}
}

// Test Get List of Data
func TestGetList(t *testing.T) {
	res, err := userres.GetList(false, -1, 1, "", "", "")
	if err != nil {
		t.Error("TestGetList , error occured!")
	}
	if len(res) == 0 {
		t.Error("TestGetList , empty result!")
	}
}

// Test Update Data
func TestUpdate(t *testing.T) {
	u := user.User{
		ID: userIdTested,
		Name: "Test", 
		Email: "test@email.com", 
		IsActive: true, 
		IdempotencyKey: "anyIdempotencyKey",
	}
	res := userres.Update(u)
	if !reflect.DeepEqual(res, u) {
		t.Error("TestUpdate , expected: ", u, " ,got: ", res)
	}
}

// Test Delete Data
func TestDelete(t *testing.T) {
	userres.Delete(userIdTested)
	res,err := userres.Get(userIdTested)
	if err != nil {
		t.Error("TestDelete , error occured!")
	}
	if res.Name != "" {
		t.Error("TestDelete , related user not deleted!")
	}
}
