package repository

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"what-to.com/internal/config"
)

// func TestCreateEntity(t *testing.T) {
// 	// Create a mock database connection
// 	db, err := sql.Open("mock", "")
// 	assert.NoError(t, err)

// 	// Create a new PgRepository instance
// 	repo := &PgRepository{DB: db}

// 	// Create a mock entity
// 	entity := &entity.Entity{Name: "Test Entity"}

// 	// Call the CreateEntity function
// 	err = repo.CreateEntity(entity)

// 	// Assert that there are no errors
// 	assert.NoError(t, err)
// }

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

	// Convert dbConfig struct to map
	configMap := make(map[string]interface{})
	configMap["database"] = make(map[string]interface{})
	reflectValue := reflect.ValueOf(dbConfig)
	reflectType := reflect.TypeOf(dbConfig)
	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Field(i)
		fieldName := reflectType.Field(i).Name
		configMap["database"].(map[string]interface{})[fieldName] = field
	}

	// Call the SetRepoConfig function
	repo.SetRepoConfig(config.ConfigT(configMap["database"].(map[interface{}]interface{})))

	// Assert that the DBConfig fields are set correctly
	assert.Equal(t, dbConfig.Host, repo.GetRepoConfig().Host)
	assert.Equal(t, dbConfig.Port, repo.GetRepoConfig().Port)
	assert.Equal(t, dbConfig.User, repo.GetRepoConfig().User)
	assert.Equal(t, dbConfig.Password, repo.GetRepoConfig().Password)
	assert.Equal(t, dbConfig.DBName, repo.GetRepoConfig().DBName)
}

// func TestGetRepoConfigStr(t *testing.T) {
// 	// Create a new PgRepository instance
// 	repo := &PgRepository{}

// 	// Set the DBConfig fields
// 	repo.DBConfig = DBConfig{
// 		Host:     "localhost",
// 		Port:     5432,
// 		User:     "test_user",
// 		Password: "test_password",
// 		DBName:   "test_db",
// 	}

// 	// Call the GetRepoConfigStr function
// 	configStr := repo.GetRepoConfigStr()

// 	// Assert that the returned string is correct
// 	expectedConfigStr := "host=localhost port=5432 user=test_user password=test_password sslmode=disable"
// 	assert.Equal(t, expectedConfigStr, configStr)
// }
