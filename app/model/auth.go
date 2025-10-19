package model
 
import ( 
    "time" 
    "github.com/golang-jwt/jwt/v5" 
    "go.mongodb.org/mongo-driver/bson/primitive"
) 
 
type User struct { 
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Username        string              `bson:"username" json:"username"` 
    PasswordHash    string              `bson:"password_hash" json:"password_hash"` 
    Email           string              `bson:"email" json:"email"` 
    Role            string              `bson:"role" json:"role"` 
    CreatedAt       time.Time           `bson:"created_at" json:"created_at"` 
} 
 
type LoginRequest struct { 
    Username        string              `bson:"username" json:"username"` 
    Password        string              `bson:"password" json:"password"` 
} 
type LoginResponse struct { 
	User            User                `bson:"user" json:"user"` 
	Token           string              `bson:"token" json:"token"` 
} 
type JWTClaims struct { 
	UserID          primitive.ObjectID  `bson:"_id,omitempty" json:"user_id"`
	Username        string              `bson:"username" json:"username"` 
	Role            string              `bson:"role" json:"role"` 
	jwt.RegisteredClaims 
} 

type RegisterRequest struct {
    Username        string              `bson:"username" json:"username" validate:"unique,required,min=3,max=50"`
    Email           string              `bson:"email" json:"email" validate:"unique,required,email"`
    Password        string              `bson:"password" json:"password" validate:"required,min=6"`
    Role            string              `bson:"role" json:"role"` 
}

