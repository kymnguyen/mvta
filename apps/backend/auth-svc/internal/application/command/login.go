package command

type LoginCommand struct {
	Email    string
	Password string
}

func (c *LoginCommand) CommandName() string {
	return "Login"
}
