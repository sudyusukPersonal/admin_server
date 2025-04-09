package main

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"google.golang.org/api/option"
)

var (
	firestoreClient *firestore.Client
	authClient      *auth.Client
	firebaseAPIKey  string // Firebase Web APIキー
)

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