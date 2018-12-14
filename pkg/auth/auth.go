package auth

const (
//secret = "superSecretSecret"
)

//JWTPayload holds the payload portion of the web token
type JWTPayload struct {
	UserId int   `json:"uId"`
	Exp    int64 `json:"exp"`
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateAuthForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserID   string `json:"userId"`
}

// CreateToken takes user and returns token for that user
//func CreateToken(userId int) string {
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
//		"userId":    userId,
//		"createdAt": time.Now().Unix(),
//	})
//
//	// Sign and get the complete encoded token as a string using the secret
//	tokenString, err := token.SignedString([]byte(secret))
//	if err != nil {
//		log.Println("couldn't create token:", err)
//	}
//
//	return tokenString
//}
//
//// GetTypeFromRequest returns Type and TypeID of the user
//func GetTypeFromRequest(req *nova.Request) (string, string) {
//	headerToken := req.Request.Header.Get("Authorization")
//	to := GetTokenFromHeader(headerToken)
//	if to == "" {
//		return "", ""
//	}
//
//	return ValidateAndGetType(to)
//
//}

// GetTokenFromHeader takes header string and strips of Bearer
//func GetTokenFromHeader(header string) string {
//	if len(header) > 7 {
//		return header[7:]
//	}
//
//	return ""
//}
//
//// ValidateAndGetType checks token and returns userID from claims
//func ValidateAndGetType(to string) (string, string) {
//	t := ""
//	typeId := ""
//
//	tk, err := jwt.Parse(to, func(token *jwt.Token) (interface{}, error) {
//		return []byte(secret), nil
//	})
//	if err != nil {
//		return "", ""
//	}
//
//	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
//		t = claims["type"].(string)
//		typeId = claims["typeId"].(string)
//	}
//
//	return t, typeId
//}

// CreatePWHash creates hash from password with random salt using bcrypt
//func CreatePWHash(pw string) (string, error) {
//	pwHash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
//	if err != nil {
//		return "", err
//	}
//
//	return string(pwHash), nil
//}
//
//// ValidatePassword checks password entered to validate it matches stored hash
//func ValidatePassword(password, hash string) error {
//	println("checking " + password + " " + hash)
//	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
//}
