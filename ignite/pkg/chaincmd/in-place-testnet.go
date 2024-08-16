package chaincmd

type InPlaceOption func([]string) []string

func InPlaceWithPrvKey(prvKey string) InPlaceOption {
	return func(s []string) []string {
		if len(prvKey) > 0 {
			return append(s, optionValidatorPrivateKey, prvKey)
		}
		return s
	}
}

func InPlaceWithAccountToFund(accounts string) InPlaceOption {
	return func(s []string) []string {
		if len(accounts) > 0 {
			return append(s, optionAccountToFund, accounts)
		}
		return s
	}
}

func InPlaceWithSkipConfirmation() InPlaceOption {
	return func(s []string) []string {
		return append(s, optionSkipConfirmation)
	}
}
