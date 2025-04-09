package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

// 管理者ページハンドラー
func HandleAdmin(c echo.Context) error {
	partyID := c.Param("party_id")
	fmt.Printf("政党ID: %s へのアクセス\n", partyID)
	return c.JSON(http.StatusOK, map[string]string{
		"message":  "管理者ページにアクセスしました",
		"party_id": partyID,
	})
}