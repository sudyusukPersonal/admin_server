package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/api/iterator"
)

// 全ポリシー取得ハンドラー
func GetPolicies(c echo.Context) error {
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

// 政党別ポリシー取得ハンドラー
func GetPoliciesByPartyID(c echo.Context) error {
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
			"message":  fmt.Sprintf("政党ID %s の政策が見つかりませんでした", partyID),
			"policies": []interface{}{},
		})
	}
	
	return c.JSON(http.StatusOK, policies)
}