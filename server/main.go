package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
)

type Config struct {
	LogLocation string
	RedisHost   string
	RedisPort   string
	MysqlHost   string
	MysqlPort   string
	MysqlUser   string
	MysqlPass   string
	MysqlDb     string
	FsPath      string
}

type App struct {
	Config *Config
	Mysql  *sql.DB
	Redis  *redis.Client
	Logger *log.Logger
}

func LoadConfig(filepath string) (*Config, error) {
	cfg, err := ini.Load(filepath)
	if err != nil {
		return nil, err
	}

	return &Config{
		LogLocation: cfg.Section("").Key("LogLocation").String(),
		RedisHost: cfg.Section("").Key("RedisHost").String(), 
		RedisPort: cfg.Section("").Key("RedisPort").String(),
		MysqlHost: cfg.Section("").Key("MysqlHost").String(),
		MysqlPort: cfg.Section("").Key("MysqlPort").String(),
		MysqlUser: cfg.Section("").Key("MysqlUser").String(),
		MysqlPass: cfg.Section("").Key("MysqlPass").String(),
		MysqlDb: cfg.Section("").Key("MysqlDb").String(),
		FsPath: cfg.Section("").Key("FsPath").String(),
	}, nil
}

func connectMysql(config Config) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.MysqlUser, config.MysqlPass, config.MysqlHost, config.MysqlPort, config.MysqlDb)
	return sql.Open("mysql", dsn)
}

func connectRedis(config Config) (*redis.Client) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})
	return rdb
}

func generateRandomToken(size int) (string, error) {
	token := make([]byte, size)

	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}

	encodedToken := base64.URLEncoding.EncodeToString(token)
	return encodedToken, nil
}

func simulateWorkload() {
	time.Sleep(1 * time.Second)
}

func (app *App) rootHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	id := r.URL.Query().Get("id")

	if id == "" {
		app.Logger.Println("Request with no id found")
		fmt.Fprintln(w, "Ok")
		return
	}

	app.Logger.Printf("Incoming request id=%s\n", id)

	redisKey := "redis_token:" + id
	fsPath := filepath.Join(app.Config.FsPath, id+".txt")

	// --- 1. Check CACHE (Redis + FS) ---
	redisToken, redisErr := app.Redis.Get(ctx, redisKey).Result()
	fsData, fsErr := os.ReadFile(fsPath)

	if redisErr == nil && fsErr == nil {
		app.Logger.Printf("CACHE HIT id=%s\n", id)

		fmt.Fprintf(w, "SUCCESS (cache)\nredis_token=%s\nfilesystem_token=%s\n",
			redisToken, string(fsData))
		return
	}

	app.Logger.Printf("CACHE MISS id=%s (redisErr=%v, fsErr=%v)\n", id, redisErr, fsErr)

	simulateWorkload()

	// --- 2. Try MySQL ---
	var dbRedisToken, dbFsToken string
	query := "SELECT redis_token, filesystem_token FROM unicorns WHERE id=? LIMIT 1"
	err := app.Mysql.QueryRow(query, id).Scan(&dbRedisToken, &dbFsToken)

	if err == nil && dbRedisToken != "" && dbFsToken != "" {
		app.Logger.Printf("MYSQL HIT id=%s\n", id)

		// Rehydrate cache
		app.Redis.Set(ctx, redisKey, dbRedisToken, 0)
		os.WriteFile(fsPath, []byte(dbFsToken), 0644)

		fmt.Fprintf(w, "SUCCESS (mysql)\nredis_token=%s\nfilesystem_token=%s\n",
			dbRedisToken, dbFsToken)
		return
	}

	if err != nil && err != sql.ErrNoRows {
		app.Logger.Printf("DB ERROR id=%s err=%v\n", id, err)
		http.Error(w, fmt.Sprintf("DB error: %v", err), http.StatusInternalServerError)
		return
	}

	// --- 3. Generate new tokens ---
	newRedisToken, err := generateRandomToken(32)
	if err != nil {
		app.Logger.Printf("TOKEN GEN ERROR (redis) id=%s err=%v\n", id, err)
		http.Error(w, "Failed to generate redis token", http.StatusInternalServerError)
		return
	}

	newFsToken, err := generateRandomToken(32)
	if err != nil {
		app.Logger.Printf("TOKEN GEN ERROR (fs) id=%s err=%v\n", id, err)
		http.Error(w, "Failed to generate filesystem token", http.StatusInternalServerError)
		return
	}

	_, err = app.Mysql.Exec(
		"INSERT INTO unicorns (id, redis_token, filesystem_token) VALUES (?, ?, ?)",
		id, newRedisToken, newFsToken,
	)
	if err != nil {
		app.Logger.Printf("DB WRITE ERROR id=%s err=%v\n", id, err)
		http.Error(w, "DB write failed", http.StatusInternalServerError)
		return
	}

	// Cache
	app.Redis.Set(ctx, redisKey, newRedisToken, 0)
	os.WriteFile(fsPath, []byte(newFsToken), 0644)

	app.Logger.Printf("GENERATED id=%s\n", id)
	fmt.Fprintf(w, "GENERATED (new)\nredis_token=%s\nfilesystem_token=%s\n",
		newRedisToken, newFsToken)
}

func main() {
	cfg, err := LoadConfig("server.ini")
	if err != nil {
		log.Fatal("Failed to read server.ini", err)
	}

	logFile, err := os.OpenFile(cfg.LogLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	logger := log.New(logFile, "", log.LstdFlags)

	mysqlDB, err := connectMysql(*cfg)
	if err != nil {
		log.Fatal("MySQL connection failed:", err)
	}

	redisClient := connectRedis(*cfg)

	app := &App{
		Config: cfg,
		Mysql:  mysqlDB,
		Redis:  redisClient,
		Logger: logger,
	}

	http.HandleFunc("/", app.rootHandler)
	app.Logger.Println("Starting server on port 80")
	log.Fatal(http.ListenAndServe("0.0.0.0:80", nil))
}
