package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	storage_go "github.com/supabase-community/storage-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

var supabaseClient *storage_go.Client

func ConnectDB() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		dbHost, dbPort, username, password, dbName)

	conn, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn,
		PreferSimpleProtocol: true}), &gorm.Config{})

	if err != nil {
		fmt.Print(err)
	}

	db = conn
}

func InitSupabaseClient() {
	apiURL := fmt.Sprintf("%s/storage/v1", os.Getenv("SUPABASE_URL"))
	apiKey := os.Getenv("STORAGE_ACCESS_KEY")

	supabaseClient = storage_go.NewClient(apiURL, apiKey, nil)
}

func GetSupabaseClient() *storage_go.Client {
	return supabaseClient
}

func GetDB() *gorm.DB {
	return db.Debug()
}
