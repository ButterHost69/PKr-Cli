

package dialer

type CallHandler struct {
	Lipaddr	string
}

type RegisterWorkspaceRequest struct {
	Username      string
	Password      string
	WorkspaceName string
}

type RegisterWorkspaceResponse struct {
	Response int
}

type RegisterUserRequest struct {
	PublicIP	string
	PublicPort	string

	Username	string
	Password	string
}

type RegisterUserResponse struct {
	UniqueUsername	string
	Response		int
}