package logic

import (
	"context"
	"go-im-user-server/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersByTokensLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUsersByTokensLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersByTokensLogic {
	return &GetUsersByTokensLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUsersByTokensLogic) GetUsersByTokens(in *user_server.GetUsersByTokensReq) (*user_server.GetUsersByTokensReply, error) {
	// todo: add your logic here and delete this line

	return &user_server.GetUsersByTokensReply{}, nil
}
