package command

type LoginCommand struct {
	Username string
	Password string
}

func (c *LoginCommand) CommandName() string {
	return "Login"
}
