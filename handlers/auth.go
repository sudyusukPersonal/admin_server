package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// ログイン処理ハンドラー
func HandleLogin(c echo.Context) error {
	ctx := context.Background()
	
	// リクエストボディからメールアドレスとパスワードを取得
	var loginReq LoginRequest
	if err := c.Bind(&loginReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "リクエスト形式が不正です",
		})
	}
	
	if loginReq.Email == "" || loginReq.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "メールアドレスとパスワードは必須です",
		})
	}
	
	fmt.Printf("ログイン試行: %s\n", loginReq.Email)
	
	// Firebase REST APIで認証
	authURL := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", firebaseAPIKey)
	
	// リクエストボディを作成
	reqBody, err := json.Marshal(map[string]interface{}{
		"email":             loginReq.Email,
		"password":          loginReq.Password,
		"returnSecureToken": true,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "リクエスト作成エラー",
		})
	}
	
	// Firebase Auth APIにリクエスト
	resp, err := http.Post(authURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("認証サービスエラー: %v", err),
		})
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "レスポンス読み取りエラー",
		})
	}
	
	// レスポンスをパース
	var authResp FirebaseAuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "レスポンスパースエラー",
		})
	}
	
	// 認証エラーの処理
	if authResp.Error.Message != "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": fmt.Sprintf("認証失敗: %s", authResp.Error.Message),
		})
	}
	
	// セッションIDを生成（UUID）
	sessionID := uuid.New().String()
	
	// Firestoreにセッション情報を保存
	_, err = firestoreClient.Collection("session_test").Doc(sessionID).Set(ctx, map[string]interface{}{
		"user_id":    authResp.LocalID,
		"created_at": time.Now(),
	})
	
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("セッション保存エラー: %v", err),
		})
	}
	
	fmt.Printf("セッション作成 - ID: %s, ユーザーID: %s\n", sessionID, authResp.LocalID)
	
	// 認証成功
	return c.JSON(http.StatusOK, map[string]string{
		"message":    "認証完了",
		"session_id": sessionID,
	})
}