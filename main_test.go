package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

func TestPrefCd2Str(t *testing.T) {
	conds := []struct {
		input string
		want  string
	}{
		{"01", "hokkaido"},
		{"13", "tokyo"},
		{"47", "okinawa"},
		{"48", ""},
		{"-1", ""},
		{"00", ""},
	}

	for _, c := range conds {
		t.Run(fmt.Sprintf("%s => %s", c.input, c.want), func(t *testing.T) {
			if got := PrefCd2Str(c.input); got != c.want {
				t.Errorf("got: %s, want: %s\n", got, c.want)
			}
		})
	}

}

func TestInitUrlValues(t *testing.T) {
	conds := []struct {
		input string
		want  url.Values
	}{
		{
			"01",
			url.Values{
				"action_kouhyou_splist": []string{"true"},
				"PrefCd":                []string{"01"},
				"OriPrefCd":             []string{"01"},
				"p_count":               []string{"1"},
				"p_offset":              []string{"0"},
				"p_sort_name":           []string{"6"},
				"p_order_name":          []string{"0"},
			},
		},
	}

	for _, c := range conds {
		t.Run(fmt.Sprintf("%s => %s", c.input, c.want), func(t *testing.T) {
			if got := initUrlValues(c.input); !reflect.DeepEqual(got, c.want) {
				t.Errorf("got: %s, want: %s\n", got, c.want)
			}
		})
	}
}

func TestFetchFacilityCount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := `{"total": 100, "list": [{"PrefCd": 13}, {"PrefCd": 33}]}`
		fmt.Fprintln(w, res)
	}))
	defer ts.Close()
	t.Skip("TODO: test code...")

	// res, err := http.PostForm(ts.URL)
}
