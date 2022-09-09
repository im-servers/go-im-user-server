package logic

import (
	"context"
	"fmt"

	"go-im-user-server/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
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
	for _, v := range in.Ids {
		deviceTokens, err := l.store.SmembersCtx(l.ctx, fmt.Sprintf("%s%v", cacheGoImServerUserDeviceTokenPrefix, v))
		if err != nil {
			err = errors.WithMessage(err, "GetSetCtx err")
			l.Logger.Error(err)
			return &user_server.GetDeviceTokensByUserIDReply{}, err
		}
		//l.Logger.Slowf("%v:%s", v, deviceTokens)
		userDeviceTokenMap[v] = &user_server.ListDeviceToken{Values: deviceTokens}
	}

	return &user_server.GetDeviceTokensByUserIDReply{UserDeviceTokens: userDeviceTokenMap}, nil
}
