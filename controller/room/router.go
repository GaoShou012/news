package room

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"wchatv1/common"
	"wchatv1/config"
	proto_room "wchatv1/proto/room"
)

func MakePassToken(ctx *gin.Context) {
	params := &proto_room.MakePassTokenReq{}
	err := ctx.BindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	rsp,err:= config.RoomServiceConfig.ServiceClient().MakePassToken(context.TODO(),params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    "200",
		"msg":     "success",
		"success": true,
		"retry":   false,
		"token":   rsp.Token,
		"data":    "",
	})
}

func Select(ctx *gin.Context) {
	var params struct {
		common.Page
	}
	if err := ctx.BindJSON(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if err := params.PageCheck(); err != nil {
		common.ResponseError(ctx, err)
		return
	}
}
