package main

import (
	"log"
	"ws-chat/db"
	"ws-chat/logger"
	"ws-chat/models"
	"ws-chat/tool"
	usercontroller "ws-chat/userController"
)

func main() {
	logger.Info("========== Init PostgreSQL ==========")
	supaClient, err := db.NewSupabaseClient()
	if err != nil {
		log.Fatalf("Failed to initialize postSQL client: %v", err)
	}
	userID := tool.GetUUID()
	// create
	newUser := models.User{
		Email:          "ricktest1@gmail.com",
		ID:             userID,
		HashedPassword: "pass111",
		FullName:       "rick111",
	}
	usercontroller.CreateUser(supaClient, newUser)
	// create

	// get
	user, err := usercontroller.GetUser(supaClient, userID)
	if err != nil {
		log.Fatalf("Failed to load user info: %v", err)
	}
	log.Printf("get user success: %+v\n", user)
	// get

	// update
	updates := map[string]interface{}{
		"full_name": "rick999",
		"email":     "rick999@gmail.com",
	}

	err = usercontroller.UpdateUser(supaClient, userID, updates)
	if err != nil {
		log.Fatalf("Failed to udpate user info: %v", err)
	}

	log.Printf("update user: %s success！", userID)
	updatedUser, err := usercontroller.GetUser(supaClient, userID)
	if err != nil {
		log.Fatalf("無法取得更新後的使用者資料: %v", err)
	}

	log.Printf("更新後的使用者資料: %+v\n", updatedUser)
	// update

	// delete
	log.Printf("嘗試刪除 ID 為 %s 的使用者...", userID)
	err = usercontroller.DeleteUser(supaClient, userID)
	if err != nil {
		log.Fatalf("無法刪除使用者資料: %v", err)
	}

	log.Printf("成功刪除 ID 為 %s 的使用者！", userID)

	//再次呼叫 GetUser 來確認資料已不存在
	log.Println("\n========== 驗證資料是否已刪除 ==========")
	userAfterDelete, err := usercontroller.GetUser(supaClient, userID)
	if userAfterDelete == nil {
		log.Printf("如預期，無法取得使用者資料： %v", err)
	} else {
		log.Println("錯誤：使用者資料仍存在！")
	}
	// delete
}
