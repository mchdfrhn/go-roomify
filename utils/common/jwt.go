// Package common menyediakan fungsi-fungsi utility umum seperti pembuatan dan verifikasi JWT token.
package common

// Import package penanganan error, konfigurasi, model data, waktu, dan library JWT.
import (
	"errors"            // Package untuk membuat error instan
	"go-roomify/config" // Package konfigurasi aplikasi
	"go-roomify/model"  // Package model entitas data
	"time"              // Package manipulasi waktu

	"github.com/dgrijalva/jwt-go" // Library penanganan standar token JWT
)

// JwtClaim merupakan struktur payload klaim custom JWT yang memperluas StandardClaims.
type JwtClaim struct {
	jwt.StandardClaims        // Embedding klaim standar JWT (Issuer, ExpiresAt, IssuedAt, dsb)
	UserId             string `json:"user_id"` // ID pengguna yang disimpan di dalam token
	Role               string `json:"role"`    // Peran pengguna (admin/ga/employee) dalam token
}

// JwtToken merupakan kontrak interface untuk utilitas pembuatan dan verifikasi JWT token.
type JwtToken interface {
	GenerateTokenJwt(user_data model.UserCredentialJwt) (string, error) // Fungsi generate token
	VerifyToken(token_string string) (jwt.MapClaims, error)            // Fungsi verifikasi token
}

// jwtToken merupakan struktur konkrit implementasi interface JwtToken.
type jwtToken struct {
	config config.TokenConfig // Menyimpan konfigurasi token (secret key, issuer, lifetime)
}

// GenerateTokenJwt membuat dan menandatangani string token JWT baru berdasarkan data user.
func (self *jwtToken) GenerateTokenJwt(user_cr model.UserCredentialJwt) (string, error) {
	// Membuat payload klaim baru untuk token JWT
	claims := JwtClaim{
		StandardClaims: jwt.StandardClaims{
			Issuer:    self.config.IssuerName,                       // Mengatur nama issuer token dari config
			ExpiresAt: time.Now().Add(self.config.JwtLifeTime).Unix(), // Mengatur waktu kadaluarsa token (unix timestamp)
		},
		UserId: user_cr.Id,   // Menyimpan ID kredensial user
		Role:   user_cr.Role, // Menyimpan role user
	}

	// Membuat objek token baru dengan algoritma penandatanganan HS256 dan klaim yang ditentukan
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Menandatangani token menggunakan JwtSignatureKey secret key
	signed_token, err := token.SignedString(self.config.JwtSignatureKey)

	// Memeriksa jika terdapat error saat proses pembentukan signature
	if err != nil {
		// Mengembalikan string kosong dan error jika gagal menandatangani
		return "", err
	}
	// Mengembalikan string token JWT yang valid dan nil error
	return signed_token, nil
}

// VerifyToken mengecek keabsahan token JWT dan mengembalikan map claims jika valid.
func (self *jwtToken) VerifyToken(token_string string) (jwt.MapClaims, error) {
	// Mem-parsing string token dan memverifikasi signature kunci rahasia
	token, err := jwt.Parse(token_string, func(token *jwt.Token) (any, error) {
		// Mengembalikan signature key dari konfigurasi untuk verifikasi
		return self.config.JwtSignatureKey, nil
	})
	// Memeriksa jika terjadi error sintaks atau signature saat parsing token
	if err != nil {
		// Mengembalikan nil dan error parsing
		return nil, err
	}

	// Mengonversi klaim token ke dalam bentuk MapClaims dan mengecek validitas token
	claims, ok := token.Claims.(jwt.MapClaims)
	// Memeriksa apakah tipe klaim sesuai dan token dalam keadaan valid
	if !ok || !token.Valid {
		// Mengembalikan error jika parsing klaim gagal atau token tidak valid
		return nil, errors.New("Failed to parse map claims or token is not valid")
	}

	// Memverifikasi apakah nilai issuer di dalam token sesuai dengan issuer di aplikasi
	if !claims.VerifyIssuer(self.config.IssuerName, true) {
		// Mengembalikan error jika issuer tidak cocok
		return nil, errors.New("Failed to Verify Issuer Name")
	}

	// Memverifikasi apakah token belum melewati batas waktu kadaluarsa (expired)
	if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
		// Mengembalikan error jika token telah kedaluwarsa
		return nil, errors.New("Token is expired")
	}

	// Mengembalikan map claims hasil verifikasi
	return claims, nil
}

// NewJwtToken menginisialisasi provider JwtToken baru dengan TokenConfig yang disediakan.
func NewJwtToken(token_config config.TokenConfig) JwtToken {
	// Mengembalikan pointer ke instance struct jwtToken
	return &jwtToken{
		config: token_config,
	}
}
