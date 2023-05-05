package logic

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"soya_milk_forum/app/usercenter/cmd/rpc/usercenter"
	"soya_milk_forum/app/usercenter/model"
	"soya_milk_forum/common/errs"
	tool "soya_milk_forum/common/tools"

	"soya_milk_forum/app/usercenter/cmd/rpc/internal/svc"
	"soya_milk_forum/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, bson.D{{"telephone_number", in.TelephoneNumber}})
	if err != nil && err != mongo.ErrNoDocuments {
		logx.Error(err, fmt.Sprintf("telephone_number:%s,err:%v", in.TelephoneNumber, err))
		return nil, errs.NewErrCode(errs.DB_ERROR)
	}

	if user != nil {
		logx.Error(fmt.Sprintf("Register user exists mobile:%s,err:%v", in.TelephoneNumber, err))
		return nil, errs.NewErrCode(errs.TELEPHONE_ALREADY_REGIST)
	}

	user = &model.User{TelephoneNumber: in.TelephoneNumber, Password: in.Password}

	if in.Username != "" {
		user.Account = in.Username
	} else {
		user.Account = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	}

	_, err = l.svcCtx.UserModel.Add(l.ctx, user)
	if err != nil {
		logx.Error(err, "telephone_number:%s,err:%v", in.TelephoneNumber, err)
		return nil, errs.NewErrCode(errs.DB_ERROR)
	}

	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&usercenter.GenerateTokenReq{
		UserId: 1,
	})

	if err != nil {
		return nil, err
	}

	return &usercenter.RegisterResp{
		AccessToken: tokenResp.AccessToken,
	}, nil
}
