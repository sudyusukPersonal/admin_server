package main

import (
	"github.com/labstack/echo/v4"
	"github.com/yourusername/yourproject/handlers"
)

// ルートの設定
func setupRoutes(e *echo.Echo) {
	// 認証不要なエンドポイント
	e.GET("/", handlers.HandleRoot)
	e.POST("/login", handlers.HandleLogin)
	
	// adminハンドラの設定（パラメータ受け取り）
	e.GET("/admin/:party_id", handlers.HandleAdmin)
	
	// policyエンドポイントの追加
	e.GET("/policy", handlers.GetPolicies)
	e.GET("/policy/:party_id", handlers.GetPoliciesByPartyID)
}