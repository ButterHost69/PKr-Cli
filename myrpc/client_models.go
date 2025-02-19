package myrpc

type PublicKeyRequest struct {
}

type PublicKeyResponse struct {
	PublicKey []byte
}

type InitWorkspaceConnectionRequest struct {
	WorkspaceName string
	MyUsername    string
	MyPublicKey   []byte

	ServerIP          string
	WorkspacePassword string
}

type InitWorkspaceConnectionResponse struct {
	Response int32 // 200 [Valid / ACK / OK] ||| 4000 [InValid / You Fucked Up Somewhere]
}
