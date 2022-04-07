package cosmosaccount

// CreateTestAccount creates an account for testing purposes within memory keyring backend
func CreateTestAccount(name string) (acc Account, mnemonic string, err error) {
	r, err := NewInMemory()
	if err != nil {
		return Account{}, "", err
	}
	return r.Create(name)
}
