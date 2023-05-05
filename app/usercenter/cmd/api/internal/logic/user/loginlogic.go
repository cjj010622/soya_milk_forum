package user

import (
	"context"
	"github.com/jinzhu/copier"
	"soya_milk_forum/app/usercenter/cmd/rpc/usercenter"
	"soya_milk_forum/common/errs"

	"soya_milk_forum/app/usercenter/cmd/api/internal/svc"
	"soya_milk_forum/app/usercenter/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResp, err := l.svcCtx.UsercenterRpc.Login(l.ctx, &usercenter.LoginReq{
		TelephoneNumber: req.TelephoneNumber,
		Password:        req.Password,
	})

	if err != nil {
		return nil, errs.FormatRpcErr(err)
	}

	resp = new(types.LoginResp)
	_ = copier.Copy(resp, loginResp)

	return resp, nil
}
