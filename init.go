package main

import (
	"github.com/yourusername/yourproject/handlers"
)

// ハンドラーの依存関係を設定
func initHandlers() {
	handlers.SetDependencies(firestoreClient, firebaseAPIKey)
}