package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	authClient      *auth.Client
	firebaseAPIKey  string // Firebase Web APIキー
)

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

func initFirebase() error {
	ctx := context.Background()
	
	// Firebase認証ファイルのパスを設定
	serviceAccountKeyFile := "config/firebaseAuth.json"
	
	// Firebase Adminの初期化
	opt := option.WithCredentialsFile(serviceAccountKeyFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return fmt.Errorf("firebase.NewApp: %v", err)
	}
	
	// Firestoreクライアントの初期化
	client, err := app.Firestore(ctx)
	if err != nil {
		return fmt.Errorf("app.Firestore: %v", err)
	}
	firestoreClient = client
	
	// Authentication クライアントの初期化
	auth, err := app.Auth(ctx)
	if err != nil {
		return fmt.Errorf("app.Auth: %v", err)
	}
	authClient = auth
	
	// Firebase Web APIキーの設定
	// 注意: 実際の実装では環境変数やConfigから取得するべき
	firebaseAPIKey = "AIzaSyB6XXyfJ0oY11JLBioRoO4jniGBXnxEBWU" // ここは実際のAPIキーに置き換えてください
	
	return nil
}

func getPolicies(c echo.Context) error {
	ctx := context.Background()
	
	// policy_testコレクションから10件のデータを取得
	iter := firestoreClient.Collection("policy_test").Limit(10).Documents(ctx)
	var policies []map[string]string
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("Error fetching documents: %v", err),
			})
		}
		
		// IDのみを含むマップを作成
		policies = append(policies, map[string]string{
			"id": doc.Ref.ID,
		})
	}
	
	return c.JSON(http.StatusOK, policies)
}

func getPoliciesByPartyID(c echo.Context) error {
	ctx := context.Background()
	partyID := c.Param("party_id")
	
	if partyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "党IDが指定されていません",
		})
	}
	
	fmt.Printf("政党ID: %s の政策を取得中\n", partyID)
	
	// policy_testコレクションからparty_idフィールドが一致するドキュメントを取得
	query := firestoreClient.Collection("policy_test").Where("party_id", "==", partyID).Limit(10)
	iter := query.Documents(ctx)
	
	var policies []map[string]interface{}
	
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": fmt.Sprintf("ドキュメント取得エラー: %v", err),
			})
		}
		
		// ドキュメントデータとIDを含むマップを作成
		data := doc.Data()
		data["id"] = doc.Ref.ID
		
		policies = append(policies, data)
	}
	
	// 結果が空の場合の処理
	if len(policies) == 0 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": fmt.Sprintf("政党ID %s の政策が見つかりませんでした", partyID),
			"policies": []interface{}{},
		})
	}
	
	return c.JSON(http.StatusOK, policies)
}

// ログイン処理ハンドラー
func handleLogin(c echo.Context) error {
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

func main() {
	// Firebase初期化
	if err := initFirebase(); err != nil {
		log.Fatalf("Firebase初期化エラー: %v", err)
		os.Exit(1)
	}
	defer firestoreClient.Close()
	
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
			"message": "管理者ページにアクセスしました",
			"party_id": partyID,
		})
	})

	// policyエンドポイントの追加 - policy_testコレクションから10件のデータID取得
	e.GET("/policy", getPolicies)
	
	// 特定の政党IDの政策を取得するエンドポイント
	e.GET("/policy/:party_id", getPoliciesByPartyID)
	
	// ログインエンドポイントの追加
	e.POST("/login", handleLogin)

	// サーバーの起動
	e.Logger.Fatal(e.Start(":8080"))
}