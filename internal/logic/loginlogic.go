package logic

import (
	"context"
	"fmt"

	"go-im-user-server/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/heyehang/go-im-pkg/trand"
	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger

	store *redis.Redis
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),

		store: redis.New(svcCtx.Config.CacheRedis[0].Host),
	}
}

func (l *LoginLogic) Login(in *user_server.LoginReq) (*user_server.LoginReply, error) {
	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err != nil {
		err = errors.WithMessage(err, "FindOneByPhone err")
		l.Logger.Error(err)
		return &user_server.LoginReply{}, err
	}

	if in.Password != user.Password {
		errors.New("user or password err")
		l.Logger.Error(err)
		return &user_server.LoginReply{}, err
	}
	//只做简单的案列
	token := &user_server.Token{
		AccessToken:  trand.RandNString(trand.RandSourceLetterAndNumber, 32),
		AccessExpire: 0,
		RefreshAfter: 0,
	}

	_, err = l.store.SaddCtx(l.ctx, fmt.Sprintf("%s%v", cacheGoImServerUserTokenPrefix, user.Id), token.AccessToken)

	if err != nil {
		err = errors.WithMessage(err, "set cache err")
		l.Logger.Error(err)
		return &user_server.LoginReply{}, err
	}

	return &user_server.LoginReply{Token: token, Id: user.Id}, nil
}
