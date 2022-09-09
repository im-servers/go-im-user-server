package logic

import (
	"context"
	"crypto/md5"
	"fmt"

	"go-im-user-server/internal/svc"
	"go-im-user-server/model"

	"github.com/heyehang/go-im-grpc/user_server"
	"github.com/heyehang/go-im-pkg/trand"

	"github.com/pkg/errors"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type CreateUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger

	store *redis.Redis
}

func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),

		store: redis.New(svcCtx.Config.CacheRedis[0].Host),
	}
}

func (l *CreateUserLogic) CreateUser(in *user_server.CreateUserReq) (*user_server.CreateUserReply, error) {

	in.Password = "123456"

	user, err := l.svcCtx.UserModel.FindOneByPhone(l.ctx, in.Phone)
	if err == nil {
		//l.Logger.Error(errors.New("user any"))
		return &user_server.CreateUserReply{Id: user.Id, Token: &user_server.Token{
			AccessToken:  trand.RandNString(trand.RandSourceLetterAndNumber, 32),
			AccessExpire: 0,
			RefreshAfter: 0,
		}}, nil

	}

	if err != nil && err.Error() != model.ErrNotFound.Error() {
		err = errors.WithMessage(err, "FindOneByPhone err")
		l.Logger.Error(err)
		return &user_server.CreateUserReply{}, err

	}

	m := md5.New()
	_, _ = m.Write([]byte(in.Password))
	res, err := l.svcCtx.UserModel.Insert(l.ctx, &model.User{
		Id:         0,
		Name:       in.Name,
		Age:        in.Age,
		Gender:     0,
		Phone:      in.Phone,
		Password:   in.Password,
		CreateAt:   0,
		UpdateAt:   0,
		CreateTime: 0,
		UpdateTime: 0,
	})

	if err != nil {
		err = errors.WithMessage(err, "Insert err")
		l.Logger.Error(err)
		return &user_server.CreateUserReply{}, err
	}

	id, _ := res.LastInsertId()
	token := &user_server.Token{
		AccessToken:  trand.RandNString(trand.RandSourceLetterAndNumber, 32),
		AccessExpire: 0,
		RefreshAfter: 0,
	}

	_, err = l.store.SaddCtx(l.ctx, fmt.Sprintf("%s%v", cacheGoImServerUserTokenPrefix, id), token.AccessToken)

	if err != nil {
		err = errors.WithMessage(err, "set cache err")
		l.Logger.Error(err)
		return &user_server.CreateUserReply{}, err
	}

	return &user_server.CreateUserReply{Token: token, Id: id}, nil
}
