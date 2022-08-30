package logic

import (
	"context"

	"go-im-user-server/rpc/internal/svc"

	"github.com/heyehang/go-im-grpc/user_server"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDeviceTokensByUserIDLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetDeviceTokensByUserIDLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDeviceTokensByUserIDLogic {
	return &GetDeviceTokensByUserIDLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetDeviceTokensByUserIDLogic) GetDeviceTokensByUserID(in *user_server.GetDeviceTokensByUserIDReq) (*user_server.GetDeviceTokensByUserIDReply, error) {
	// todo: add your logic here and delete this line

	return &user_server.GetDeviceTokensByUserIDReply{}, nil
}
