package main

import (
	"context"
	"flag"
	"fmt"
	"go-im-user-server/internal/config"
	"math/rand"
	"sync"
	"time"

	"github.com/heyehang/go-im-pkg/trand"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var configFile = flag.String("f", "../etc/userserver.yaml", "the config file")

const (
	cacheGoImServerUserDeviceTokenPrefix = "cache:goImServer:user:devicetoken:id:"
	cacheGoImServerUserTokenPrefix       = "cache:goImServer:user:token:id:"
)

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	// conn := sqlx.NewMysql(c.Mysql.DataSource)
	// UserModel := model.NewUserModel(conn, c.CacheRedis)
	store := redis.New(c.CacheRedis[0].Host)

	bitchCount := 1000
	wg := new(sync.WaitGroup)
	ctx := context.Background()
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// users := make([]*model.User, 0, bitchCount)

			// for i := 0; i < bitchCount; i++ {
			// 	users = append(users, &model.User{
			// 		Id:         0,
			// 		Name:       trand.RandNString(trand.RandSourceLetterAndNumber, 5),
			// 		Age:        trand.RandNInt(0, 80),
			// 		Gender:     0,
			// 		Phone:      trand.RandNString(trand.RandSourceNumber, 15), //测试用，长度变大一点，避免uk冲突
			// 		Password:   "123456",
			// 		CreateAt:   0,
			// 		UpdateAt:   0,
			// 		CreateTime: 0,
			// 		UpdateTime: 0,
			// 	})
			// 	_, err := UserModel.Inserts(ctx, users)
			// 	if err != nil {
			// 		panic(err)
			// 	}
			// }

			store.PipelinedCtx(ctx, func(p redis.Pipeliner) error {
				for i := 1; i <= bitchCount; i++ {
					p.SAdd(ctx, fmt.Sprintf("%s%v", cacheGoImServerUserTokenPrefix, i), trand.RandNString(trand.RandSourceLetterAndNumber, 32))
				}
				res := p.Command(ctx)
				if res.Err() != nil {
					panic(res.Err())
				}
				if len(res.Val()) != bitchCount {
					panic("res vals len err")
				}
				return nil
			})

		}()
	}

	wg.Wait()
}
