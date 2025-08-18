package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TheJa750/PrayerPals/internal/database"
	"github.com/google/uuid"
)

var validRoles = []string{"admin", "member"}

func getValidRoles() []string {
	rolesCopy := make([]string, len(validRoles))
	copy(rolesCopy, validRoles)
	return rolesCopy
}

type JsonError struct {
	Message string `json:"error"`
}

type APIConfig struct {
	DBQueries *database.Queries
	JWTSecret string
}

type UserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type UserLoggedIn struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Group struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     uuid.UUID `json:"owner_id"`
	InviteCode  string    `json:"invite_code"` // New field for invite code
}

type PostRequest struct {
	Content string `json:"content"`
}

type Post struct {
	ID           uuid.UUID `json:"id"`
	GroupID      uuid.UUID `json:"group_id"`
	UserID       uuid.UUID `json:"user_id"`
	Content      string    `json:"content"`
	CreatedAt    string    `json:"created_at"`
	Author       string    `json:"author"`        // Username of the post author
	CommentCount int64     `json:"comment_count"` // Number of comments on the post
	Comments     []Comment `json:"comments"`      // Comments associated with the post
}

type Comment struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	GroupID   uuid.UUID `json:"group_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"created_at"`
	Author    string    `json:"author"` // Username of the comment author
}

type PromoteUserRequest struct {
	Role string `json:"role"` // e.g., "admin", "member"
}

type PostFeedRequest struct {
	GroupID uuid.UUID `json:"group_id"`
	Limit   int       `json:"limit"`
	Offset  int       `json:"offset"`
}

type ModerateUserRequest struct {
	Action string `json:"action"` // e.g., "kick", "ban"
	Reason string `json:"reason"` // Reason for the action
}

type UserJoinGroup struct {
	UserID    uuid.UUID `json:"user_id"`
	GroupID   uuid.UUID `json:"group_id"`
	GroupName string    `json:"group_name"`
	Role      string    `json:"role"` // e.g., "admin", "member"
}

type UpdateUserRequest struct {
	Username string `json:"username"` // Optional, can be empty if not updating
	Password string `json:"password"`
}

type GroupMember struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"` // e.g., "admin", "member"
}

type UpdateInviteCodeRequest struct {
	InviteCode string `json:"invite_code"` // New invite code to set for the group
}

func ParseJSON[T any](r *http.Request) (T, error) {
	var data T

	if r.Body == nil {
		return data, errors.New("request body is empty")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&data); err != nil {
		return data, errors.New("invalid request body")
	}

	return data, nil
}

func CreateJSONResponse[T any](data T, w http.ResponseWriter, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return err
	}
	return nil
}
