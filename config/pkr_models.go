package config

type PKRConfig struct {
	WorkspaceName  string       `json:"workspace_name"`
	AllConnections []Connection `json:"all_connections"`
	LastHash       string       `json:"last_hash"`
}

type Connection struct {
	Username      string `json:"username"`
	ServerName    string `json:"server_name"`
	PublicKeyPath string `json:"public_key_path"`
}