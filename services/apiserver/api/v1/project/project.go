package project

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/aenjoy/iot-lubricant/pkg/model/request"
	"github.com/aenjoy/iot-lubricant/pkg/types/exception"
	exceptionCode "github.com/aenjoy/iot-lubricant/pkg/types/exception/code"
	"github.com/aenjoy/iot-lubricant/services/apiserver/api/v1/helper"
	"github.com/bytedance/sonic"
	"github.com/gin-gonic/gin"
)

func (a Api) AddProject(c *gin.Context) {
	req := helper.RequestBind[request.AddProjectRequest](c)
	if req == nil {
		return
	}
	claims, err := helper.GetClaims(c)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.New(exceptionCode.ErrorGetClaimsFailed, exception.WithMsg("claims is empty")), c)
		return
	}
	projectId, err := a.IProjectService.AddProject(c, claims.User.Id, "", req.ProjectName, req.Description)
	if err != nil {
		helper.FailedWithJson(http.StatusInternalServerError,
			exception.ErrNewException(err, exceptionCode.ErrorAddProjectFailed), c)
		return
	}
	if req.DSNLinkerInfo != nil {
		var dsn string
		switch req.DataBaseType {
		case "mysql":
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
				req.DSNLinkerInfo.User, req.DSNLinkerInfo.Pass, req.DSNLinkerInfo.Host, req.DSNLinkerInfo.Port,
				req.DSNLinkerInfo.Db)
		case "TDEngine":
			dsn, _ = sonic.MarshalString(req.DSNLinkerInfo)
		default:
			helper.FailedWithJson(http.StatusBadRequest, exception.New(exceptionCode.UnsupportedDatabaseType), c)
			return
		}

		dsn = base64.StdEncoding.EncodeToString([]byte(dsn))
		err = a.IProjectService.AddDataStoreEngine(c, projectId, dsn, req.DataBaseType, req.Description, req.StoreTable)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptionCode.ErrorAddProjectFailed), c)
			return
		}
	}
	if req.Washer != nil {
		_, err := a.IProjectService.AddWasher(c, req.Washer)
		if err != nil {
			helper.FailedWithJson(http.StatusInternalServerError,
				exception.ErrNewException(err, exceptionCode.AddWasherFailed), c)
			return
		}
	}
	helper.SuccessJson(projectId, c)
}

func (a Api) RemoveProject(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}
