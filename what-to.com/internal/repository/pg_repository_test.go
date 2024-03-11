package repository

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"what-to.com/internal/config"
	"what-to.com/internal/entity"
)

func TestCreateEntity(t *testing.T) {
	// Create a mock database connection
	db, err := sql.Open("mock", "")
	assert.NoError(t, err)

	// Create a new PgRepository instance
	repo := &PgRepository{DB: db}

	// Create a mock entity
	entity := &entity.Entity{Name: "Test Entity"}

	// Call the CreateEntity function
	err = repo.CreateEntity(entity)

	// Assert that there are no errors
	assert.NoError(t, err)
}

func TestConnectToRepo(t *testing.T) {
	// Create a new PgRepository instance
	repo := &PgRepository{}

	// Call the ConnectToRepo function
	repo.ConnectToRepo()

	// Assert that the DB field is not nil
	assert.NotNil(t, repo.DB)
}

func TestSetRepoConfig(t *testing.T) {
	// Create a new PgRepository instance
	repo := &PgRepository{}

	// Create a mock DBConfig
	dbConfig := DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test_user",
		Password: "test_password",
		DBName:   "test_db",
	}

	// Call the SetRepoConfig function
	repo.SetRepoConfig(config.ConfigT(dbConfig))

	// Assert that the DBConfig fields are set correctly
	assert.Equal(t, dbConfig.Host, repo.DBConfig.Host)
	assert.Equal(t, dbConfig.Port, repo.DBConfig.Port)
	assert.Equal(t, dbConfig.User, repo.DBConfig.User)
	assert.Equal(t, dbConfig.Password, repo.DBConfig.Password)
	assert.Equal(t, dbConfig.DBName, repo.DBConfig.DBName)
}

func TestGetRepoConfigStr(t *testing.T) {
	// Create a new PgRepository instance
	repo := &PgRepository{}

	// Set the DBConfig fields
	repo.DBConfig = DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "test_user",
		Password: "test_password",
		DBName:   "test_db",
	}

	// Call the GetRepoConfigStr function
	configStr := repo.GetRepoConfigStr()

	// Assert that the returned string is correct
	expectedConfigStr := "host=localhost port=5432 user=test_user password=test_password sslmode=disable"
	assert.Equal(t, expectedConfigStr, configStr)
}
