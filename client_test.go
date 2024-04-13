package bilibili

import (
	"net/http"
	"reflect"
	"testing"
)

func TestCookie(t *testing.T) {
	result := []*http.Cookie{
		{Name: "a", Value: "1"},
		{Name: "b", Value: "2"},
	}
	{
		c := New()
		c.SetCookiesString("a=1; b=2")
		if c.GetCookiesString() != "a=1; b=2" {
			t.Fail()
		}
		if !reflect.DeepEqual(c.GetCookies(), result) {
			t.Fail()
		}
	}
	{
		c := New()
		c.SetCookiesString("a=1\nb=2")
		if c.GetCookiesString() != "a=1; b=2" {
			t.Fail()
		}
		if !reflect.DeepEqual(c.GetCookies(), result) {
			t.Fail()
		}
	}
}
