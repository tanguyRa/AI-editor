# config

```tree
config/
├── README.md
└── config.go
    ├── type Config {Environment: string, Address: string, Encryption: EncryptionConfig, Database: DatabaseConfig, Polar: PolarConfig}
    ├── type PolarConfig {WebhookSecret: string}
    ├── type EncryptionConfig {Key: string}
    ├── type DatabaseConfig {ConnectionString: string}
    ├── func Load() (*Config, error)
    ├── func loadFromFile(path string, config *Config) error
    ├── func loadFromEnv(config *Config)
    ├── func setDefaults() *Config
    └── func validate(config *Config) error
```
