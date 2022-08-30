package logic

import (
	"context"

	"go-im-user-server/rpc/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user_server.LoginReq) (*user_server.LoginReply, error) {
	// todo: add your logic here and delete this line

	return &user_server.LoginReply{}, nil
}
