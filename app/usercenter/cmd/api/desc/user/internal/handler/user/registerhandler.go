package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/logic/user"
	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/svc"
	"soya_milk_forum/app/usercenter/cmd/api/desc/user/internal/types"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
