package frontend

import (
	"github.com/gin-gonic/gin"
	fe "github.com/huhaophp/hblog/internal/entity/frontend"
	srv "github.com/huhaophp/hblog/internal/service"
	"github.com/huhaophp/hblog/internal/service/frontend"
)

var Notice = cNotice{}

type cNotice struct{}

func (*cNotice) HomePage(ctx *gin.Context) {
	s := srv.Context(ctx)

	if !s.Check() {
		s.To("/login").WithError("请登录后，再继续操作").Redirect()
		return
	}

	var req fe.GetRemindListReq
	if err := ctx.ShouldBind(&req); err != nil {
		s.To("/").WithError(err).Redirect()
		return
	}

	noticeService := frontend.NoticeService(ctx)

	if data, err := noticeService.GetList(&req); err != nil {
		s.To("/").WithError(err.Error()).Redirect()
	} else {
		// 提醒未读数量
		remindUnread, _ := frontend.NoticeService(ctx).GetRemindUnread()
		// 系统未读数量
		systemUnread, _ := frontend.NoticeService(ctx).GetSystemUnread()

		data["remindUnread"] = remindUnread
		data["systemUnread"] = systemUnread

		// 更新未读消息状态
		noticeService.ReadAll(req.Type)

		s.View("frontend.notice.home", data)
	}
}