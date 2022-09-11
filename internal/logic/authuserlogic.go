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

type AuthUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger

	store *redis.Redis
}

func NewAuthUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthUserLogic {
	return &AuthUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
		store:  redis.New(svcCtx.Config.CacheRedis[0].Host),
	}
}

func (l *AuthUserLogic) AuthUser(in *user_server.AuthUserReq) (*user_server.AuthUserReply, error) {
	any, err := l.store.SismemberCtx(l.ctx, fmt.Sprintf("%s%v", cacheGoImServerUserTokenPrefix, in.Id), in.AccessToken)
	if err != nil {
		err = errors.WithMessage(err, "SismemberCtx err")
		return &user_server.AuthUserReply{}, err
	}
	if !any {
		err = errors.New("SismemberCtx not found")
		return &user_server.AuthUserReply{}, err
	}

	_, err = l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil {
		err = errors.WithMessage(err, "FindOne err")
		l.Logger.Error(err)
		return &user_server.AuthUserReply{}, err
	}
	_, err = l.store.SaddCtx(l.ctx, fmt.Sprintf("%s%v", cacheGoImServerUserDeviceTokenPrefix, in.Id), in.DeviceToken)

	if err != nil {
		err = errors.WithMessage(err, "SaddCtx err")
		l.Logger.Error(err)
		return &user_server.AuthUserReply{}, err
	}
	return &user_server.AuthUserReply{AccessToken: in.AccessToken, DeviceToken: in.DeviceToken}, nil
}
