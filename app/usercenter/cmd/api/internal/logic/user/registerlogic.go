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

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResp, err := l.svcCtx.UsercenterRpc.Register(l.ctx, &usercenter.RegisterReq{
		TelephoneNumber: req.TelephoneNumber,
		Password:        req.Password,
	})

	if err != nil {
		return nil, errs.FormatRpcErr(err)
	}

	resp = new(types.RegisterResp)
	_ = copier.Copy(resp, registerResp)

	return resp, nil
}
