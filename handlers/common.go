package handlers

import (
	"cloud.google.com/go/firestore"
	"time"
)

// グローバル変数（依存関係）
var (
	firestoreClient *firestore.Client
	firebaseAPIKey  string
)

// 共通構造体定義
// ログインリクエスト用の構造体
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Firebaseからの認証レスポンス用の構造体
type FirebaseAuthResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"` // ユーザーID
	DisplayName  string `json:"displayName,omitempty"`
	Registered   bool   `json:"registered"`
	Error        struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

// セッション情報用の構造体
type Session struct {
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// 依存関係設定用関数
func SetDependencies(client *firestore.Client, apiKey string) {
	firestoreClient = client
	firebaseAPIKey = apiKey
}