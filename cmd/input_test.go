package cmd

type MockInputController struct {
	TextInput          string
	FailTextInputAfter int
	TextInputError     error
	Password           string
	FailPasswordAfter  int
	PasswordError      error
}

func (mr *MockInputController) requestTextInput(string) (string, error) {
	if mr.FailTextInputAfter != 0 {
		mr.FailTextInputAfter--
		return mr.TextInput, nil
	}
	return mr.TextInput, mr.TextInputError
}

func (mr *MockInputController) requestPassword(string) (string, error) {
	if mr.FailPasswordAfter != 0 {
		mr.FailPasswordAfter--
		return mr.Password, nil
	}
	return mr.Password, mr.PasswordError
}
