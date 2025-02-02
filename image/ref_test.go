package image

import (
	"testing"

	"github.com/cpuguy83/go-docker/errdefs"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

func TestParseRef(t *testing.T) {
	type testCase struct {
		ref      string
		host     string
		locator  string
		tag      string
		errCheck func(error) bool
	}

	testCases := []testCase{
		{ref: "", errCheck: errdefs.IsInvalid},
		{ref: "foo", host: "docker.io", locator: "foo", tag: "latest"},
		{ref: "foo:latest", host: "docker.io", locator: "foo", tag: "latest"},
		{ref: "foo:other", host: "docker.io", locator: "foo", tag: "other"},
		{ref: "foo/bar", host: "docker.io", locator: "foo/bar", tag: "latest"},
		{ref: "foo/bar:latest", host: "docker.io", locator: "foo/bar", tag: "latest"},
		{ref: "foo/bar:other", host: "docker.io", locator: "foo/bar", tag: "other"},
		{ref: "foo/bar/baz:latest", host: "docker.io", locator: "foo/bar/baz", tag: "latest"},
		{ref: "foo/bar/baz:other", host: "docker.io", locator: "foo/bar/baz", tag: "other"},
		{ref: "docker.io/foo/bar", host: "docker.io", locator: "foo/bar", tag: "latest"},
		{ref: "docker.io/foo/bar:latest", host: "docker.io", locator: "foo/bar", tag: "latest"},
		{ref: "foo:5000/bar", host: "foo:5000", locator: "bar", tag: "latest"},
		{ref: "foo:5000/bar:latest", host: "foo:5000", locator: "bar", tag: "latest"},
		{ref: "foo:5000/bar/baz", host: "foo:5000", locator: "bar/baz", tag: "latest"},
		{ref: "foo:5000/bar/baz:latest", host: "foo:5000", locator: "bar/baz", tag: "latest"},
		{ref: "foo:invalid/bar/baz", errCheck: errdefs.IsInvalid},
		{ref: "foo:invalid/bar/baz:latest", errCheck: errdefs.IsInvalid},
	}

	format := func(host, locator, tag string) string {
		return "host=" + host + " locator=" + locator + " tag=" + tag
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.ref, func(t *testing.T) {
			r, err := ParseRef(tc.ref)
			if tc.errCheck == nil {
				tc.errCheck = func(err error) bool {
					return err == nil
				}
			}
			if !tc.errCheck(err) {
				t.Error("unexpected error:", err)
			}
			assert.Check(t, cmp.Equal(tc.host, r.Host), format(r.Host, r.Locator, r.Tag))
			assert.Check(t, cmp.Equal(tc.locator, r.Locator), format(r.Host, r.Locator, r.Tag))
			assert.Check(t, cmp.Equal(tc.tag, r.Tag), format(r.Host, r.Locator, r.Tag))
		})
	}
}
