package logic

import (
	"context"
	"fmt"

	"go-im-user-server/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"

	goredis "github.com/go-redis/redis/v8"
)

type GetDeviceTokensByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger

	store *redis.Redis
}

func NewGetDeviceTokensByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceTokensByUserIDLogic {
	return &GetDeviceTokensByUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),

		store: redis.New(svcCtx.Config.CacheRedis[0].Host),
	}
}

func (l *GetDeviceTokensByUserIDLogic) GetDeviceTokensByUserID(in *user_server.GetDeviceTokensByUserIDReq) (*user_server.GetDeviceTokensByUserIDReply, error) {
	var userDeviceTokenMap = make(map[int64]*user_server.ListDeviceToken, 0)

	cmdRes := make([]*goredis.StringSliceCmd, 0, len(in.Ids))
	err := l.store.PipelinedCtx(l.ctx, func(p redis.Pipeliner) error {
		for _, v := range in.Ids {
			cmdRes = append(cmdRes, p.SMembers(l.ctx, fmt.Sprintf("%s%v", cacheGoImServerUserDeviceTokenPrefix, v)))
		}
		return nil
	})
	if err != nil {
		l.Logger.Error(err)
		return &user_server.GetDeviceTokensByUserIDReply{}, err
	}
	for i, v := range cmdRes {
		deviceTokens, err := v.Result()
		if err != nil {
			l.Logger.Error(err)
			return &user_server.GetDeviceTokensByUserIDReply{}, err
		}
		userDeviceTokenMap[in.Ids[i]] = &user_server.ListDeviceToken{Values: deviceTokens}
	}

	return &user_server.GetDeviceTokensByUserIDReply{UserDeviceTokens: userDeviceTokenMap}, nil
}
