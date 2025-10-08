package middlewares

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Auth0Middleware struct {
	Domain   string
	Audience string
	jwks     *JWKS
	mu       sync.RWMutex
}

type JWKS struct {
	Keys      []JWK
	ExpiresAt time.Time
}

type JWK struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

func NewAuth0(domain, audience string) *Auth0Middleware {
	return &Auth0Middleware{
		Domain:   domain,
		Audience: audience,
		jwks:     nil,
	}
}

// RequireAuth0 validates Auth0 JWT tokens
func (m *Auth0Middleware) RequireAuth0() gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}

		tokenString := strings.TrimSpace(h[len("Bearer "):])

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify signing method
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Get key ID from token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, errors.New("kid header not found")
			}

			// Get public key from JWKS
			publicKey, err := m.getPublicKey(kid)
			if err != nil {
				return nil, err
			}

			return publicKey, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}

		// Verify audience if configured
		if m.Audience != "" {
			aud, _ := claims["aud"].(string)
			if aud != m.Audience {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid audience"})
				return
			}
		}

		// Verify issuer
		expectedIssuer := m.Domain
		if !strings.HasSuffix(expectedIssuer, "/") {
			expectedIssuer += "/"
		}
		iss, _ := claims["iss"].(string)
		if iss != expectedIssuer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid issuer"})
			return
		}

		// Extract Auth0 user ID (sub claim)
		auth0ID, _ := claims["sub"].(string)
		if auth0ID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing sub claim"})
			return
		}

		// Set Auth0 ID in context for downstream handlers
		c.Set("auth0_id", auth0ID)
		c.Set("user_claims", claims)

		c.Next()
	}
}

// getPublicKey fetches the public key from Auth0 JWKS endpoint
func (m *Auth0Middleware) getPublicKey(kid string) (*rsa.PublicKey, error) {
	// Check if we have cached JWKS and it's not expired
	m.mu.RLock()
	if m.jwks != nil && time.Now().Before(m.jwks.ExpiresAt) {
		m.mu.RUnlock()
		return m.findKey(kid)
	}
	m.mu.RUnlock()

	// Fetch new JWKS
	m.mu.Lock()
	defer m.mu.Unlock()

	jwksURL := fmt.Sprintf("%s.well-known/jwks.json", m.Domain)
	if !strings.HasSuffix(m.Domain, "/") {
		jwksURL = fmt.Sprintf("%s/.well-known/jwks.json", m.Domain)
	}

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch JWKS: status %d", resp.StatusCode)
	}

	var jwks struct {
		Keys []JWK `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	// Cache JWKS for 1 hour
	m.jwks = &JWKS{
		Keys:      jwks.Keys,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	return m.findKey(kid)
}

// findKey locates the JWK by kid and converts it to RSA public key
func (m *Auth0Middleware) findKey(kid string) (*rsa.PublicKey, error) {
	for _, key := range m.jwks.Keys {
		if key.Kid == kid {
			return m.jwkToPublicKey(key)
		}
	}
	return nil, errors.New("key not found")
}

// jwkToPublicKey converts JWK to RSA public key
func (m *Auth0Middleware) jwkToPublicKey(jwk JWK) (*rsa.PublicKey, error) {
	// Decode modulus
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	// Decode exponent
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)

	publicKey := &rsa.PublicKey{
		N: n,
		E: int(e.Int64()),
	}

	return publicKey, nil
}
