package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	goUtil "gitlab.com/mihael97/Go-utility/src/util"
	"strings"
)

func GetAppName(ctx *gin.Context) (*string, error) {
	path := ctx.Request.URL.Path
	if len(path) == 0 {
		return nil, fmt.Errorf("path is empty")
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "api" {
		return nil, fmt.Errorf("path doesn't start with /api")
	}
	return goUtil.GetPointer(parts[2]), nil
}
