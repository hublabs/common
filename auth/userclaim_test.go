package auth

import (
	"net/http"
	"testing"

	"github.com/pangpanglabs/goutils/test"
)

var testToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJjb2xsZWFndWUiLCJ0ZW5hbnRDb2RlIjoiaHVibGFicyIsImNvbGxlYWd1ZUlkIjoxMjR9.ZDUrsd1ZtH4pQA14EAJdTQUhEACBfd5TsAXqITpab1s"

func TestUserClaim(t *testing.T) {
	userClaim, err := UserClaim{}.FromToken(testToken)
	test.Ok(t, err)
	test.Equals(t, userClaim.TenantCode, "hublabs")
	test.Equals(t, userClaim.ColleagueId, int64(124))
	test.Equals(t, userClaim.Audience, "colleague")
}

func TestUserClaimFromHttpReq(t *testing.T) {
	req := &http.Request{
		Header: map[string][]string{},
	}
	req.Header.Add("Authorization", testToken)
	req.Header.Add("X-Username", "Jodan")
	req.Header.Add("X-Brand-Code", "NIKE")
	req.Header.Add("X-Store-Id", "1512")
	req.Header.Add("X-Store-Province", "SHH")

	userClaim, err := newUserClaimFromHttpReq(req)
	test.Ok(t, err)

	test.Equals(t, userClaim.BrandCode, "NIKE")
	test.Equals(t, userClaim.StoreId, int64(1512))
	test.Equals(t, userClaim.StoreProvince, "SHH")
	test.Equals(t, userClaim.Username, "Jodan")

}
