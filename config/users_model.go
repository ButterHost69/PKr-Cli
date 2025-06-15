package config

type UsersConfig struct {
	User        string         `json:"user"`
	ServerLists []ServerConfig `json:"server_lists"`
}

type ServerConfig struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	ServerAlias string `json:"server_alias"`
	ServerIP    string `json:"server_ip"`

	SendWorkspaces []SendWorkspaceFolder `json:"send_workspace"`
	GetWorkspaces  []GetWorkspaceFolder  `json:"get_workspace"`
}

type SendWorkspaceFolder struct {
	WorkspaceName     string `json:"workspace_name"`
	WorkspacePath     string `json:"workspace_path"`
	WorkSpacePassword string `json:"workspace_password"`
}

type GetWorkspaceFolder struct {
	WorkspaceOwnerName string `json:"workspace_owner_name"`
	WorkspaceName      string `json:"workspace_name"`
	WorkspacePassword  string `json:"workspace_password"`
	WorkspacePath      string `json:"workspace_path"`
	LastHash           string `json:"last_hash"`
}
