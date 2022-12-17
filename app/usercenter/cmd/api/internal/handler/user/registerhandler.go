package user

import (
	"net/http"
	"soya_milk_forum/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
	"soya_milk_forum/app/usercenter/cmd/api/internal/logic/user"
	"soya_milk_forum/app/usercenter/cmd/api/internal/svc"
	"soya_milk_forum/app/usercenter/cmd/api/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		response.HttpResult(r, w, resp, err)
	}
}
