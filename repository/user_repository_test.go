package repository_test

import (
	"fmt"
	"os"
	"pojok-baca-api/model"
	"pojok-baca-api/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestPostgresDB(t *testing.T) *gorm.DB {
	t.Helper()

	user := getenv("DB_USER", "postgres")
	pass := getenv("DB_PASS", "mraihan")
	host := getenv("DB_HOST", "localhost")
	port := getenv("DB_PORT", "5432")
	name := getenv("DB_NAME", "db-pojok-baca-test")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, pass, name, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect database: %v", err)
	}

	if err := db.Migrator().DropTable(&model.User{}); err != nil {
		t.Fatalf("failed to drop user table: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}); err != nil {
		t.Fatalf("failed to migrate user table: %v", err)
	}

	return db
}

func getenv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func TestUserRepository_CreateAndGetByEmail_Postgres(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := repository.NewUserRepository(db)

	user := model.User{
		Name:     "Bob",
		Email:    "bob@example.com",
		Password: "securepass123",
		Role:     "admin",
	}

	savedUser, err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, savedUser.ID)

	foundUser, err := repo.GetByEmail(user.Email)
	assert.NoError(t, err)
	assert.Equal(t, savedUser.Email, foundUser.Email)
	assert.Equal(t, savedUser.Name, foundUser.Name)
}

func TestUserRepository_GetByEmail_NotFound_Postgres(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := repository.NewUserRepository(db)

	_, err := repo.GetByEmail("notfound@example.com")
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestUserRepository_GetByID_Postgres(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := repository.NewUserRepository(db)

	user := model.User{
		Name:     "Evan",
		Email:    "evan@example.com",
		Password: "123456",
		Role:     "user",
	}

	savedUser, err := repo.Create(user)
	assert.NoError(t, err)
	assert.NotZero(t, savedUser.ID)

	// Ambil by ID
	foundUser, err := repo.GetByID(savedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, savedUser.ID, foundUser.ID)
	assert.Equal(t, savedUser.Email, foundUser.Email)
}

func TestUserRepository_UpdateDeposit_Postgres(t *testing.T) {
	db := setupTestPostgresDB(t)
	repo := repository.NewUserRepository(db)

	// Buat user tanpa deposit
	user := model.User{
		Name:     "Dina",
		Email:    "dina@example.com",
		Password: "pass321",
		Role:     "user",
	}

	savedUser, err := repo.Create(user)
	assert.NoError(t, err)
	assert.Nil(t, savedUser.Deposit)

	// Update deposit
	newDeposit := 150000
	updatedUser, err := repo.UpdateDeposit(newDeposit, savedUser.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedUser.Deposit)
	assert.Equal(t, newDeposit, *updatedUser.Deposit)

	// Verifikasi langsung dari DB
	checkedUser, err := repo.GetByID(savedUser.ID)
	assert.NoError(t, err)
	assert.Equal(t, newDeposit, *checkedUser.Deposit)
}
