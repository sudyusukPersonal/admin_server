package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// ルートハンドラー
func HandleRoot(c echo.Context) error {
	fmt.Println("アクセス")
	return c.String(http.StatusOK, "アクセスがありました")
}