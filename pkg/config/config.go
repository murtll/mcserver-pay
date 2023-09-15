package config

import (
	"fmt"
	"strings"

	"github.com/murtll/mcserver-pay/pkg/util"
)

var ListenAddr = util.GetStrOrDefault("LISTEN_ADDR", ":8020")
var HealthPath = util.GetStrOrDefault("HEALTH_PATH", "/_healthz")

var PostgresString = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
	util.GetStrOrDefault("POSTGRES_HOST", "localhost"),
	util.GetStrOrDefault("POSTGRES_USER", "test"),
	util.GetStrOrDefault("POSTGRES_PASSWORD", "test"),
	util.GetStrOrDefault("POSTGRES_DB", "test"),
	util.GetIntOrDefault("POSTGRES_PORT", 5432),
	util.GetStrOrDefault("POSTGRES_SSL", "disable"))

var FkTrustedIps = strings.Split(
	util.GetStrOrDefault("FK_TRUSTED_IPS", "168.119.157.136,168.119.60.227,138.201.88.124,178.154.197.79"), ",")

var FkSigningKey = util.GetStrOrDefault("FK_SIGNING_KEY", "")
var FkMerchantID = util.GetIntOrDefault("FK_MERCHANT_ID", 0)

var Version = util.GetStrOrDefault("APP_VERSION", "0.1.0")

var ApiUrl = util.GetStrOrDefault("API_URL", "http://localhost:3001")
