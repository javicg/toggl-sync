package cmd

type MockInputController struct {
	TextInput      string
	TextInputError error
	Password       string
	PasswordError  error
}

func (mr MockInputController) RequestTextInput(string) (string, error) {
	return mr.TextInput, mr.TextInputError
}

func (mr MockInputController) RequestPassword(string) (string, error) {
	return mr.Password, mr.PasswordError
}
