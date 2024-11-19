package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/senseyman/go-dconf/manager"
	"github.com/senseyman/go-dconf/repository/redis"
)

type MyTestConfig struct {
	FeeValue    float32
	UserTxLimit int
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	wg := &sync.WaitGroup{}
	var mx sync.Mutex

	initCfg := &MyTestConfig{
		FeeValue:    0.01,
		UserTxLimit: 5,
	}

	redisCli, err := redis.New(ctx, redis.Config{
		Address: "localhost:6379",
	}, "test-app")
	if err != nil {
		panic(err)
	}
	cfgManager := manager.New(redisCli, initCfg, time.Second*5)
	if err := cfgManager.Run(ctx, wg); err != nil {
		log.Printf("err loading init config: %v", err)
	}

	// imitate outside service that changes our configs
	go func() {
		time.Sleep(time.Second * 10)
		i := float32(1)
		for {
			log.Println("updating text config...")
			if err := redisCli.UpdateConfig(ctx, MyTestConfig{
				FeeValue:    initCfg.FeeValue*0.1 + (i / 10),
				UserTxLimit: initCfg.UserTxLimit + int(i),
			}); err != nil {
				log.Printf("err updating test config: %v\n", err)
			}

			i++
			time.Sleep(time.Second * 10)
		}
	}()

	go func() {
		for {
			log.Println("getting test config")
			mx.Lock()
			initCfg = cfgManager.GetConfig()
			mx.Unlock()
			log.Printf("the cfg: %v\n", initCfg)
			time.Sleep(time.Second * 3)
		}
	}()

	wg.Wait()
}
