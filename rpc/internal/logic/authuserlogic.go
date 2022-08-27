package logic

import (
	"context"

	"go-im-user-server/rpc/internal/svc"
	"go-im-user-server/rpc/types/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type AuthUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAuthUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuthUserLogic {
	return &AuthUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AuthUserLogic) AuthUser(in *user.AuthUserReq) (*user.AuthUserReply, error) {
	// todo: add your logic here and delete this line

	return &user.AuthUserReply{}, nil
}
