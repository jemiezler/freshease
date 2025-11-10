package password

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"freshease/backend/ent"
	"freshease/backend/ent/role"
	"freshease/backend/ent/user"
)

type Service struct {
	db        *ent.Client
	jwtSecret []byte
	ttl       time.Duration
}

func NewService(db *ent.Client) *Service {
	ttlMin := getenvIntDefault("JWT_ACCESS_TTL_MIN", 15)
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret" // fallback for development
	}

	return &Service{
		db:        db,
		jwtSecret: []byte(secret),
		ttl:       time.Duration(ttlMin) * time.Minute,
	}
}

func getenvIntDefault(key string, def int) int {
	val := os.Getenv(key)
	if val == "" {
		return def
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return def
	}
	return i
}

// Login authenticates a user with email and password
func (s *Service) Login(ctx context.Context, email, password string) (string, *ent.User, error) {
	// Find user by email
	u, err := s.db.User.Query().
		Where(user.Email(email)).
		WithRole().
		First(ctx)
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Check if user has a password
	if u.Password == nil || *u.Password == "" {
		return "", nil, errors.New("password not set for this user")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.issueJWT(u.ID, u.Email)
	if err != nil {
		return "", nil, err
	}

	return token, u, nil
}

// InitAdmin creates an admin user with the admin role
// This should only be called once during initial setup
func (s *Service) InitAdmin(ctx context.Context, email, password, name string) (*ent.User, error) {
	// Check if admin role exists, create if not
	adminRole, err := s.db.Role.Query().
		Where(role.Name("admin")).
		First(ctx)
	if err != nil {
		// Create admin role if it doesn't exist
		adminRole, err = s.db.Role.Create().
			SetName("admin").
			SetDescription("Administrator role with full system access").
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	// Check if any admin user already exists (by checking for users with admin role)
	adminUsers, err := s.db.User.Query().
		Where(user.HasRoleWith(role.Name("admin"))).
		Count(ctx)
	if err == nil && adminUsers > 0 {
		return nil, errors.New("admin user already exists")
	}

	// Check if user with this email already exists
	existingUser, err := s.db.User.Query().
		Where(user.Email(email)).
		First(ctx)
	if err == nil {
		// User exists, assign admin role
		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		updatedUser, err := s.db.User.UpdateOneID(existingUser.ID).
			SetPassword(string(hashed)).
			SetName(name).
			SetRoleID(adminRole.ID).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		// Reload with role
		return s.db.User.Query().Where(user.ID(updatedUser.ID)).WithRole().First(ctx)
	}

	// Create new admin user
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	newUser, err := s.db.User.Create().
		SetEmail(email).
		SetName(name).
		SetPassword(string(hashed)).
		SetRoleID(adminRole.ID).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Reload with role
	return s.db.User.Query().Where(user.ID(newUser.ID)).WithRole().First(ctx)
}

// issueJWT generates a JWT token for a user
func (s *Service) issueJWT(uid uuid.UUID, email string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   uid.String(),
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(s.ttl).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(s.jwtSecret)
}

