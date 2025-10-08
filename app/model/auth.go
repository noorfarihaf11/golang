package model
 
import ( 
    "time" 
    "github.com/golang-jwt/jwt/v5" 
) 
 
type User struct { 
    ID        int       `json:"id"` 
    Username  string    `json:"username"` 
    PasswordHash string    `json:"password_hash"` 
    Email     string    `json:"email"` 
    Role      string    `json:"role"` 
    CreatedAt time.Time `json:"created_at"` 
} 
 
type LoginRequest struct { 
    Username string `json:"username"` 
    Password string `json:"password"` 
} 
type LoginResponse struct { 
	User  User   `json:"user"` 
	Token string `json:"token"` 
} 
type JWTClaims struct { 
	UserID   int    `json:"user_id"` 
	Username string `json:"username"` 
	Role     string `json:"role"` 
	jwt.RegisteredClaims 
} 

type RegisterRequest struct {
    Username string `json:"username" validate:"unique,required,min=3,max=50"`
    Email    string `json:"email" validate:"unique,required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

