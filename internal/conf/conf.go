package conf

import "github.com/spf13/viper"

// Config struct will store the configuration values provided by user
type Config struct {
	Realm        string
	Endpoint     string
	ClientID     string
	ClientSecret string
	ClientScope  string
}

var (
	defaults = map[string]interface{}{
		"realm":    "master",
		"endpoint": "localhost",
		"scope":    "openid",
	}
	configName  = "config"
	configPaths = []string{
		".",
		"/opt/kc-pgsql-pam",
		"/etc/",
		"$HOME/.config/",
	}
)

func LoadConfig() (config Config, err error) {
	for k, v := range defaults {
		viper.SetDefault(k, v)
	}
	for _, p := range configPaths {
		viper.AddConfigPath(p)
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")

	viper.SetEnvPrefix("kc_pgsql")  // Becomes "KC_PGSQL"
	viper.BindEnv("Realm")        // KC_PGSQL_REALM
	viper.BindEnv("Endpoint")     // KC_PGSQL_ENDPOINT
	viper.BindEnv("ClientID")     // KC_PGSQL_CLIENTID
	viper.BindEnv("ClientSecret") // KC_PGSQL_CLIENTSECRET
	viper.BindEnv("ClientScope")  // KC_PGSQL_CLIENTSCOPE

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
