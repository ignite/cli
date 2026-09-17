package testing_test

import "testing"

func Test_newRealm(t *testing.T) {
	ur := testing.NewUserRealm("g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5")
	const urExpected = "UserRealm{ g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5 }"
	if ur.String() != urExpected {
		t.Errorf("got %q want %q", ur.String(), urExpected)
	}
	cr := testing.NewCodeRealm("gno.land/r/sys/users")
	const crExpected = "CodeRealm{ g1njxh4leja7h52ea0lnq9crx3j6782g77nc7yd4, gno.land/r/sys/users }"
	if cr.String() != crExpected {
		t.Errorf("got %q want %q", cr.String(), crExpected)
	}
}
