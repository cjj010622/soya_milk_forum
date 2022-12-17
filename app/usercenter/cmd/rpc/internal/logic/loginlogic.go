package logic

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"soya_milk_forum/app/usercenter/cmd/rpc/usercenter"
	"soya_milk_forum/common/errs"

	"soya_milk_forum/app/usercenter/cmd/rpc/internal/svc"
	"soya_milk_forum/app/usercenter/cmd/rpc/pb"

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

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, bson.D{{"telephone_number", in.TelephoneNumber},
		{"is_deleted", 0}})

	if err != nil && err != mongo.ErrNoDocuments {
		logx.Error(err, fmt.Sprintf("telephone_number:%s,err:%v", in.TelephoneNumber, err))
		return nil, errs.NewErrCode(errs.DB_ERROR)
	}

	if err == mongo.ErrNoDocuments {
		return nil, errs.NewErrCode(errs.USER_NOT_EXIST)
	}

	if user.Password != in.Password {
		return nil, errs.NewErrCode(errs.PASSWORD_ERR)
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: 1,
	})

	if err != nil {
		return nil, err
	}

	return &usercenter.LoginResp{
		AccessToken: tokenResp.AccessToken,
	}, nil

}
