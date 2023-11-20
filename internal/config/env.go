package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"
)

type values struct {
	MAILHEAP_ENV                      string
	MAILHEAP_ENV_DIR                  string
	MAILHEAP_TEMP_DIR                 string
	MAILHEAP_DB_LOCATION              string
	MAILHEAP_SHUTDOWN_TIMEOUT         time.Duration
	MAILHEAP_HTTP_TCP_ADDRESS         string
	MAILHEAP_HTTP_MAX_REQUEST_SIZE    int64
	MAILHEAP_HTTP_UPLOAD_MEMBUF_SIZE  int64
	MAILHEAP_HTTP_ENABLE_PROMETHEUS   bool
	MAILHEAP_HTTP_ENABLE_SHUTDOWN     bool
	MAILHEAP_SMTP_USERNAME            string
	MAILHEAP_SMTP_PASSWORD            string
	MAILHEAP_SMTP_NETWORK_TYPE        string
	MAILHEAP_SMTP_ADDRESS             string
	MAILHEAP_SMTP_DOMAIN              string
	MAILHEAP_SMTP_READ_TIMEOUT        time.Duration
	MAILHEAP_SMTP_WRITE_TIMEOUT       time.Duration
	MAILHEAP_SMTP_MAX_MSG_BYTES       int64
	MAILHEAP_SMTP_MAX_RECIPIENTS      int64
	MAILHEAP_SMTP_MAX_LINE_LENGTH     int64
	MAILHEAP_SMTP_ALLOW_INSECURE_AUTH bool
	MAILHEAP_SMTP_DISABLE_AUTH        bool
	MAILHEAP_SMTP_ENABLE_SMTPUTF8     bool
	MAILHEAP_SMTP_ENABLE_LMTP         bool
	MAILHEAP_SMTP_ENABLE_REQUIRETLS   bool
	MAILHEAP_SMTP_ENABLE_BINARYMIME   bool
	MAILHEAP_SMTP_ENABLE_DSN          bool
}

var v values

func (v *values) print() {
	buf := new(strings.Builder)
	buf.WriteString("Environment has been resolved to:\n")
	val := reflect.Indirect(reflect.ValueOf(v))
	valType := val.Type()
	valNumField := val.NumField()
	for i := 0; i < valNumField; i++ {
		a := valType.Field(i).Name
		b := val.Field(i).Interface()
		buf.WriteString(fmt.Sprintf("%-40s= %v\n", a, b))
	}
	log.Print(buf.String())
}

func GetEnv() string {
	return v.MAILHEAP_ENV
}

func GetEnvDir() string {
	return v.MAILHEAP_ENV_DIR
}

func GetTempDir() string {
	return v.MAILHEAP_TEMP_DIR
}

func GetDBLocation() string {
	return v.MAILHEAP_DB_LOCATION
}

func GetHTTPTCPAddress() string {
	return v.MAILHEAP_HTTP_TCP_ADDRESS
}

func GetHTTPMaxRequestSize() int64 {
	return v.MAILHEAP_HTTP_MAX_REQUEST_SIZE
}

func GetHTTPUploadMemoryBufferSize() int64 {
	return v.MAILHEAP_HTTP_UPLOAD_MEMBUF_SIZE
}

func IsHTTPEnablePrometheus() bool {
	return v.MAILHEAP_HTTP_ENABLE_PROMETHEUS
}

func IsHTTPEnableShutdown() bool {
	return v.MAILHEAP_HTTP_ENABLE_SHUTDOWN
}

func GetShutdownTimeout() time.Duration {
	return v.MAILHEAP_SHUTDOWN_TIMEOUT
}

func GetSMTPUsername() string {
	return v.MAILHEAP_SMTP_USERNAME
}

func GetSMTPPassword() string {
	return v.MAILHEAP_SMTP_PASSWORD
}

func GetSMTPNetworkType() string {
	return v.MAILHEAP_SMTP_NETWORK_TYPE
}

func GetSMTPAddress() string {
	return v.MAILHEAP_SMTP_ADDRESS
}

func GetSMTPDomain() string {
	return v.MAILHEAP_SMTP_DOMAIN
}

func GetSMTPReadTimeout() time.Duration {
	return v.MAILHEAP_SMTP_READ_TIMEOUT
}

func GetSMTPWriteTimeout() time.Duration {
	return v.MAILHEAP_SMTP_WRITE_TIMEOUT
}

func GetSMTPMaxMessageBytes() int64 {
	return v.MAILHEAP_SMTP_MAX_MSG_BYTES
}

func GetSMTPMaxRecipients() int64 {
	return v.MAILHEAP_SMTP_MAX_RECIPIENTS
}

func GetSMTPMaxLineLength() int64 {
	return v.MAILHEAP_SMTP_MAX_LINE_LENGTH
}

func IsSMTPAllowInsecureAuth() bool {
	return v.MAILHEAP_SMTP_ALLOW_INSECURE_AUTH
}

func IsSMTPDisableAuth() bool {
	return v.MAILHEAP_SMTP_DISABLE_AUTH
}

func IsSMTPEnableSMTPUTF8() bool {
	return v.MAILHEAP_SMTP_ENABLE_SMTPUTF8
}

func IsSMTPEnableLMTP() bool {
	return v.MAILHEAP_SMTP_ENABLE_LMTP
}

func IsSMTPEnableREQUIRETLS() bool {
	return v.MAILHEAP_SMTP_ENABLE_REQUIRETLS
}

func IsSMTPEnableBINARYMIME() bool {
	return v.MAILHEAP_SMTP_ENABLE_BINARYMIME
}

func IsSMTPEnableDSN() bool {
	return v.MAILHEAP_SMTP_ENABLE_DSN
}
