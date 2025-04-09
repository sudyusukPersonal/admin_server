package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echoインスタンスの作成
	e := echo.New()

	// ミドルウェアの設定
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルートハンドラの設定
	e.GET("/", func(c echo.Context) error {
		fmt.Println("アクセス")
		return c.String(http.StatusOK, "アクセスがありました")
	})

	// adminハンドラの設定（パラメータ受け取り）
	e.GET("/admin/:party_id", func(c echo.Context) error {
		partyID := c.Param("party_id")
		fmt.Printf("政党ID: %s へのアクセス\n", partyID)
		return c.JSON(http.StatusOK, map[string]string{
			"message": "管理者ページにアクセスしまし",
			"party_id": partyID,
		})
	})


	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}