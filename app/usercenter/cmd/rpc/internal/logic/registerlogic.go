package logic

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

var ErrUserAlreadyRegisterError = errs.NewErrMsg("user has been registered")

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
		return nil, errors.Wrapf(errs.NewErrCode(errs.DB_ERROR), "telephone_number:%s,err:%v",
			in.TelephoneNumber, err)
	}

	if user != nil {
		return nil, errors.Wrapf(ErrUserAlreadyRegisterError, "Register user exists mobile:%s,err:%v",
			in.TelephoneNumber, err)
	}

	user = &model.User{TelephoneNumber: in.TelephoneNumber, Password: in.Password}

	if in.Username != "" {
		user.Account = in.Username
	} else {
		user.Account = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	}

	_, err = l.svcCtx.UserModel.Add(l.ctx, user)
	if err != nil {
		return nil, errors.Wrapf(errs.NewErrCode(errs.DB_ERROR), "telephone_number:%s,err:%v", in.TelephoneNumber, err)
	}

	return &pb.RegisterResp{}, nil
}
