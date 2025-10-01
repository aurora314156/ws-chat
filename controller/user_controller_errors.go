package controller

import "errors"

var ErrUserNotFound = errors.New("user not found")

// ErrNoRowsAffected 用於表示 Update 或 Delete 操作成功執行，但影響的行數為 0。
// 它可以被視為找不到資源的一種情況，或者表示資料沒有發生變化 (304 Not Modified)。
var ErrNoRowsAffected = errors.New("no rows were affected by the operation")

// ErrUserAlreadyExists 用於表示嘗試創建已存在的使用者 (例如 email 或 username 唯一性衝突)。
// 在 HTTP Handler 層，這通常會被轉換為 HTTP 409 Conflict 或 400 Bad Request。
var ErrUserAlreadyExists = errors.New("user already exists (duplicate key)")
