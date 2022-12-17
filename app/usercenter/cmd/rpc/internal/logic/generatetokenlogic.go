package logic

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"soya_milk_forum/common/ctx_data"
	"soya_milk_forum/common/errs"
	"time"

	"soya_milk_forum/app/usercenter/cmd/rpc/internal/svc"
	"soya_milk_forum/app/usercenter/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, in.UserId)
	if err != nil {
		logx.Error(err, fmt.Sprintf("getJwtToken err userId:%d , err:%v", in.UserId, err))
		return nil, errs.NewErrCode(errs.CREATE_TOKEN_ERR)
	}

	return &pb.GenerateTokenResp{
		AccessToken: accessToken,
	}, nil

}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {

	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims[ctx_data.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}
