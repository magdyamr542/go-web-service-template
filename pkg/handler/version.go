package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/magdyamr542/go-web-service-template/pkg/api"
)

// Injected at build time
var (
	Version      string
	Description  = "Manage learning resources based on tags"
	BuildTime    string
	CommitTime   string
	CommitSHA    string
	CommitBranch string
)

type versionHandler struct {
}

func newVersionHandler() *versionHandler {
	return &versionHandler{}
}

func (h *versionHandler) GetVersion(ctx echo.Context) error {
	version := api.VersionInfo{
		BuildTime:    BuildTime,
		CommitBranch: CommitBranch,
		CommitSHA:    CommitSHA,
		CommitTime:   CommitTime,
		Description:  Description,
		Version:      Version,
	}
	return ctx.JSON(http.StatusOK, version)
}
