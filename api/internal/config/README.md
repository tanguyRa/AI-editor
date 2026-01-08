# config

```tree
config/
├── README.md
└── config.go
    ├── type Config {Environment: string, Encryption: EncryptionConfig}
    ├── type EncryptionConfig {Key: string}
    ├── func Load() (*Config, error)
    ├── func loadFromFile(path string, config *Config) error
    ├── func loadFromEnv(config *Config)
    ├── func setDefaults() *Config
    └── func validate(config *Config) error
```
