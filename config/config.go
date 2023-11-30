package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"

	log "github.com/StarkZarn/gophish/logger"
)


// AdminServer represents the Admin server configuration details
type AdminServer struct {
	ListenURL            string   `json:"listen_url"`
	UseTLS               bool     `json:"use_tls"`
	CertPath             string   `json:"cert_path"`
	KeyPath              string   `json:"key_path"`
	CSRFKey              string   `json:"csrf_key"`
	AllowedInternalHosts []string `json:"allowed_internal_hosts"`
	TrustedOrigins       []string `json:"trusted_origins"`
}

// PhishServer represents the Phish server configuration details
type PhishServer struct {
	ListenURL string `json:"listen_url"`
	UseTLS    bool   `json:"use_tls"`
	CertPath  string `json:"cert_path"`
	KeyPath   string `json:"key_path"`
}

// Config represents the configuration information.
type Config struct {
	AdminConf      AdminServer `json:"admin_server"`
	PhishConf      PhishServer `json:"phish_server"`
	DBName         string      `json:"db_name"`
	DBPath         string      `json:"db_path"`
	DBSSLCaPath    string      `json:"db_sslca_path"`
	MigrationsPath string      `json:"migrations_prefix"`
	TestFlag       bool        `json:"test_flag"`
	ContactAddress string      `json:"contact_address"`
	Logging        *log.Config `json:"logging"`
}

// Version contains the current gophish version
var Version = ""

// ServerName is the server type that is returned in the transparency response.
const ServerName = "gophish"

func replaceEnvVars(input string) string {
	re := regexp.MustCompile("\\$\\{([A-Za-z_]+)\\}")
	return re.ReplaceAllStringFunc(input, func(match string) string {
		envVarName := re.FindStringSubmatch(match)[1]
		envVarValue := os.Getenv(envVarName)
		return envVarValue
	})
}

func replaceEnvVarsInMap(m map[string]interface{}) {
	for key, value := range m {
		switch v := value.(type) {
		case string:
			m[key] = replaceEnvVars(v)
		case map[string]interface{}:
			replaceEnvVarsInMap(v)
		case []interface{}:
			replaceEnvVarsInSlice(v)
		}
	}
}

func replaceEnvVarsInSlice(s []interface{}) {
	for i, v := range s {
		switch val := v.(type) {
		case string:
			s[i] = replaceEnvVars(val)
		case map[string]interface{}:
			replaceEnvVarsInMap(val)
		case []interface{}:
			replaceEnvVarsInSlice(val)
		}
	}
}

func unmarshalConfig(data []byte, v interface{}) error {
	return json.Unmarshal([]byte(replaceEnvVars(string(data))), v)
}

// LoadConfig loads the configuration from the specified filepath
func LoadConfig(filepath string) (*Config, error) {
	// Get the config file
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = unmarshalConfig(configFile, config)
	if err != nil {
		return nil, err
	}

	if config.Logging == nil {
		config.Logging = &log.Config{}
	}
	// Choosing the migrations directory based on the database used.
	config.MigrationsPath = config.MigrationsPath + config.DBName
	// Explicitly set the TestFlag to false to prevent config.json overrides
	config.TestFlag = false
	return config, nil
}
