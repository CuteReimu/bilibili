package bilibili

import (
	"net/http"
	"strings"
	"testing"
)

func deepEquals(a, b []*http.Cookie) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Name != b[i].Name {
			return false
		}
		if a[i].Value != b[i].Value {
			return false
		}
		if a[i].Domain != b[i].Domain {
			return false
		}
		if a[i].Path != b[i].Path {
			return false
		}
	}
	return true
}

func TestCookie(t *testing.T) {
	{
		result := []*http.Cookie{ // 不测试Expires，因为time.Time转化成string再转化回来会丢失单调时钟
			{Name: "a", Value: "1", Domain: "bilibili.com", Path: "/"},
			{Name: "b", Value: "2", Domain: "bilibili.com", Path: "/"},
		}
		s := make([]string, 0, len(result))
		for _, cookie := range result {
			s = append(s, cookie.String())
		}
		c := New()
		c.SetCookiesString(strings.Join(s, "\n"))
		if !deepEquals(result, c.GetCookies()) {
			t.Fail()
		}
	}
	{
		result := []*http.Cookie{
			{Name: "a", Value: "1"},
			{Name: "b", Value: "2"},
		}
		c := New()
		c.SetRawCookies("a=1; b=2")
		if c.GetCookiesString() != "a=1\nb=2" {
			t.Fail()
		}
		if !deepEquals(result, c.GetCookies()) {
			t.Fail()
		}
	}
}
