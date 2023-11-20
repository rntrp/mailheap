package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func Load() {
	// https://github.com/joho/godotenv#precedence--conventions
	loadDotEnv()
	tryLoad(v.MAILHEAP_ENV_DIR, ".env."+v.MAILHEAP_ENV+".local")
	if v.MAILHEAP_ENV != "test" {
		tryLoad(v.MAILHEAP_ENV_DIR, ".env.local")
	}
	tryLoad(v.MAILHEAP_ENV_DIR, ".env."+v.MAILHEAP_ENV)
	tryLoad(v.MAILHEAP_ENV_DIR, ".env")
	loadEnv()
	v.print()
}

func tryLoad(path, file string) {
	f := filepath.Join(path, file)
	if err := godotenv.Load(f); err == nil {
		log.Println("Loaded config from " + f)
	}
}

func loadDotEnv() {
	v.MAILHEAP_ENV = os.Getenv("MAILHEAP_ENV")
	if len(v.MAILHEAP_ENV) == 0 {
		v.MAILHEAP_ENV = "development"
	}
	v.MAILHEAP_ENV_DIR = os.Getenv("MAILHEAP_ENV_DIR")
}

func loadEnv() {
	v.MAILHEAP_TEMP_DIR = parseString("MAILHEAP_TEMP_DIR", os.TempDir())
	v.MAILHEAP_DB_LOCATION = parseString("MAILHEAP_DB_LOCATION", filepath.Join(os.TempDir(), "mailheap.db"))
	v.MAILHEAP_SHUTDOWN_TIMEOUT = parseDuration("MAILHEAP_SHUTDOWN_TIMEOUT", 0, time.Second)
	v.MAILHEAP_HTTP_TCP_ADDRESS = parseString("MAILHEAP_HTTP_TCP_ADDRESS", ":8080")
	v.MAILHEAP_HTTP_MAX_REQUEST_SIZE = parseInt64("MAILHEAP_HTTP_MAX_REQUEST_SIZE", -1)
	v.MAILHEAP_HTTP_UPLOAD_MEMBUF_SIZE = parseInt64("MAILHEAP_HTTP_UPLOAD_MEMBUF_SIZE", 10<<20)
	v.MAILHEAP_HTTP_ENABLE_PROMETHEUS = parseBool("MAILHEAP_HTTP_ENABLE_PROMETHEUS", false)
	v.MAILHEAP_HTTP_ENABLE_SHUTDOWN = parseBool("MAILHEAP_HTTP_ENABLE_SHUTDOWN", false)
	v.MAILHEAP_SMTP_USERNAME = parseString("MAILHEAP_SMTP_USERNAME", "username")
	v.MAILHEAP_SMTP_PASSWORD = parseString("MAILHEAP_SMTP_PASSWORD", "password")
	v.MAILHEAP_SMTP_NETWORK_TYPE = parseString("MAILHEAP_SMTP_NETWORK_TYPE", "tcp")
	v.MAILHEAP_SMTP_ADDRESS = parseString("MAILHEAP_SMTP_ADDRESS", ":2525")
	v.MAILHEAP_SMTP_DOMAIN = parseString("MAILHEAP_SMTP_DOMAIN", "localhost")
	v.MAILHEAP_SMTP_READ_TIMEOUT = parseDuration("MAILHEAP_SMTP_READ_TIMEOUT", 10, time.Second)
	v.MAILHEAP_SMTP_WRITE_TIMEOUT = parseDuration("MAILHEAP_SMTP_WRITE_TIMEOUT", 10, time.Second)
	v.MAILHEAP_SMTP_MAX_MSG_BYTES = parseInt64("MAILHEAP_SMTP_MAX_MSG_BYTES", 50*1024*1024)
	v.MAILHEAP_SMTP_MAX_RECIPIENTS = parseInt64("MAILHEAP_SMTP_MAX_RECIPIENTS", 50)
	v.MAILHEAP_SMTP_MAX_LINE_LENGTH = parseInt64("MAILHEAP_SMTP_MAX_LINE_LENGTH", 2000)
	v.MAILHEAP_SMTP_ALLOW_INSECURE_AUTH = parseBool("MAILHEAP_SMTP_ALLOW_INSECURE_AUTH", false)
	v.MAILHEAP_SMTP_DISABLE_AUTH = parseBool("MAILHEAP_SMTP_DISABLE_AUTH", false)
	v.MAILHEAP_SMTP_ENABLE_SMTPUTF8 = parseBool("MAILHEAP_SMTP_ENABLE_SMTPUTF8", false)
	v.MAILHEAP_SMTP_ENABLE_LMTP = parseBool("MAILHEAP_SMTP_ENABLE_LMTP", false)
	v.MAILHEAP_SMTP_ENABLE_REQUIRETLS = parseBool("MAILHEAP_SMTP_ENABLE_REQUIRETLS", false)
	v.MAILHEAP_SMTP_ENABLE_BINARYMIME = parseBool("MAILHEAP_SMTP_ENABLE_BINARYMIME", false)
	v.MAILHEAP_SMTP_ENABLE_DSN = parseBool("MAILHEAP_SMTP_ENABLE_DSN", false)
}

func parseBool(env string, def bool) bool {
	if b, err := strconv.ParseBool(os.Getenv(env)); err == nil {
		return b
	}
	return def
}

func parseString(env, def string) string {
	if s, ok := os.LookupEnv(env); ok {
		return s
	}
	return def
}

func parseInt64(env string, def int64) int64 {
	if i, err := strconv.ParseInt(os.Getenv(env), 10, 64); err == nil {
		return i
	}
	return def
}

func parseDuration(env string, def int64, unit time.Duration) time.Duration {
	return time.Duration(parseInt64(os.Getenv(env), def)) * unit
}
