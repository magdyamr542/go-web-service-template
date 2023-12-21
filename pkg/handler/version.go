package handler

import (
	"net/http"

	"github.com/magdyamr542/go-web-service-template/pkg/api"
	"github.com/magdyamr542/go-web-service-template/pkg/version"

	"github.com/labstack/echo/v4"
)

type versionHandler struct {
}

func newVersionHandler() *versionHandler {
	return &versionHandler{}
}

func (h *versionHandler) GetVersion(ctx echo.Context) error {
	version := api.VersionInfo{
		BuildTime:    version.BuildTime,
		CommitBranch: version.CommitBranch,
		CommitSHA:    version.CommitSHA,
		CommitTime:   version.CommitTime,
		Description:  version.Description,
		Version:      version.Version,
	}
	return ctx.JSON(http.StatusOK, version)
}
