package structtag_test

import (
	"sort"
	"strings"
	"testing"

	"github.com/fatih/structtag"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	test := []struct {
		name string
		tag  string
		exp  []*structtag.Tag
	}{
		{
			name: "empty tag",
			tag:  "",
		},
		{
			name: "tag with one key (valid)",
			tag:  `json:""`,
			exp: []*structtag.Tag{
				{
					Key: "json",
				},
			},
		},
		{
			name: "tag with one key and dash name",
			tag:  `json:"-"`,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "-",
				},
			},
		},
		{
			name: "tag with key and name",
			tag:  `json:"foo"`,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with key, name and option",
			tag:  `json:"foo,omitempty"`,
			exp: []*structtag.Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{"omitempty"},
				},
			},
		},
		{
			name: "tag with multiple keys",
			tag:  `json:"" hcl:""`,
			exp: []*structtag.Tag{
				{
					Key: "json",
				},
				{
					Key: "hcl",
				},
			},
		},
		{
			name: "tag with multiple keys and names",
			tag:  `json:"foo" hcl:"foo"`,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with multiple keys and names",
			tag:  `json:"foo" hcl:"foo"`,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "foo",
				},
			},
		},
		{
			name: "tag with multiple keys and different names",
			tag:  `json:"foo" hcl:"bar"`,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:  "hcl",
					Name: "bar",
				},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo,omitempty" structs:"bar,omitnested"`,
			exp: []*structtag.Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{"omitempty"},
				},
				{
					Key:     "structs",
					Name:    "bar",
					Options: []string{"omitnested"},
				},
			},
		},
		{
			name: "tag with multiple keys, different names and options",
			tag:  `json:"foo" structs:"bar,omitnested" hcl:"-"`,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
				{
					Key:     "structs",
					Name:    "bar",
					Options: []string{"omitnested"},
				},
				{
					Key:  "hcl",
					Name: "-",
				},
			},
		},
		{
			name: "tag with quoted name",
			tag:  `json:"foo,bar:\"baz\""`,
			exp: []*structtag.Tag{
				{
					Key:     "json",
					Name:    "foo",
					Options: []string{`bar:"baz"`},
				},
			},
		},
		{
			name: "tag with trailing space",
			tag:  `json:"foo" `,
			exp: []*structtag.Tag{
				{
					Key:  "json",
					Name: "foo",
				},
			},
		},
	}

	for _, ts := range test {
		t.Run(ts.name, func(t *testing.T) {
			tags, err := structtag.Parse(ts.tag)
			require.NoError(t, err)
			require.Equal(t, ts.exp, tags.Tags())

			trimmedInput := strings.TrimSpace(ts.tag)
			require.Equal(t, trimmedInput, tags.String())
		})
	}
}

func TestParseErr(t *testing.T) {
	for _, ts := range []struct {
		name      string
		tag       string
		expectErr error
	}{
		{
			name:      "EOF after key",
			tag:       "json",
			expectErr: structtag.ErrTagSyntax,
		},
		{
			name:      "invalid value EOF after colon",
			tag:       "json:",
			expectErr: structtag.ErrTagSyntax,
		},
		{
			name:      "missing key",
			tag:       ":\"value\"",
			expectErr: structtag.ErrTagKeySyntax,
		},
		{
			name:      "invalid value",
			tag:       "json:name",
			expectErr: structtag.ErrTagValueSyntax,
		},
		{
			name:      "invalid value EOF after reverse-solidus",
			tag:       `json:"\`,
			expectErr: structtag.ErrTagValueSyntax,
		},
		{
			name:      "space after colon",
			tag:       "json: ",
			expectErr: structtag.ErrTagValueSyntax,
		},
	} {
		t.Run(ts.name, func(t *testing.T) {
			require.Error(t, ts.expectErr)
			tags, err := structtag.Parse(ts.tag)
			require.Error(t, err)
			require.Equal(t, ts.expectErr, err)
			require.Nil(t, tags)
		})
	}
}

func TestTags_Get(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo,omitempty" structs:"bar,omitnested"`)
	require.NoError(t, err)

	found, err := tags.Get("json")
	require.NoError(t, err)

	t.Run("String", func(t *testing.T) {
		require.Equal(t, `json:"foo,omitempty"`, found.String())
	})
	t.Run("Value", func(t *testing.T) {
		require.Equal(t, `foo,omitempty`, found.Value())
	})
}

func TestTags_Keys(t *testing.T) {
	t.Run("multiple tags", func(t *testing.T) {
		tags, err := structtag.Parse(`json:"foo,omitempty" structs:"bar,omitnested"`)
		require.NoError(t, err)

		require.Equal(t, []string{"json", "structs"}, tags.Keys())
	})

	t.Run("no tags", func(t *testing.T) {
		tags, err := structtag.Parse(``)
		require.NoError(t, err)

		require.Equal(t, []string{}, tags.Keys())
	})
}

func TestTags_Set(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo,omitempty" structs:"bar,omitnested"`)
	require.NoError(t, err)

	err = tags.Set(&structtag.Tag{
		Key:     "json",
		Name:    "bar",
		Options: []string{},
	})
	require.NoError(t, err)

	found, err := tags.Get("json")
	require.NoError(t, err)

	require.Equal(t, `json:"bar"`, found.String())
}

func TestTags_Set_Append(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo,omitempty"`)
	require.NoError(t, err)

	err = tags.Set(&structtag.Tag{
		Key:     "structs",
		Name:    "bar",
		Options: []string{"omitnested"},
	})
	require.NoError(t, err)

	found, err := tags.Get("structs")
	require.NoError(t, err)

	require.Equal(t, `structs:"bar,omitnested"`, found.String())
	require.Equal(t, `json:"foo,omitempty" structs:"bar,omitnested"`, tags.String())
}

func TestTags_Set_KeyDoesNotExist(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo,omitempty" structs:"bar,omitnested"`)
	require.NoError(t, err)

	err = tags.Set(&structtag.Tag{
		Key:     "",
		Name:    "bar",
		Options: []string{},
	})
	require.Error(t, err, "setting tag with a nonexisting key should error")
	require.Equal(t, structtag.ErrKeyNotSet, err)
}

func TestTags_Set_TagDoesNotExist(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo,omitempty" structs:"bar,omitnested"`)
	require.NoError(t, err)
	tag, err := tags.Get("toml")
	require.Equal(t, structtag.ErrTagNotExist, err)
	require.Nil(t, tag)
}

func TestTags_Delete(t *testing.T) {
	tags, err := structtag.Parse(
		`json:"foo,omitempty" structs:"bar,omitnested" hcl:"-"`,
	)
	require.NoError(t, err)

	tags.Delete("structs")
	require.Equal(t, 2, tags.Len())

	found, err := tags.Get("json")
	require.NoError(t, err)

	require.Equal(t, `json:"foo,omitempty"`, found.String())
	require.Equal(t, `json:"foo,omitempty" hcl:"-"`, tags.String())
}

func TestTags_DeleteOptions(t *testing.T) {
	tags, err := structtag.Parse(
		`json:"foo,omitempty" structs:"bar,omitnested,omitempty" hcl:"-"`,
	)
	require.NoError(t, err)

	tags.DeleteOptions("json", "omitempty")

	require.Equal(t,
		`json:"foo" structs:"bar,omitnested,omitempty" hcl:"-"`,
		tags.String(),
	)

	tags.DeleteOptions("structs", "omitnested")
	require.Equal(t, `json:"foo" structs:"bar,omitempty" hcl:"-"`, tags.String())
}

func TestTags_AddOption(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo" structs:"bar,omitempty" hcl:"-"`)
	require.NoError(t, err)

	tags.AddOptions("json", "omitempty")

	require.Equal(t,
		`json:"foo,omitempty" structs:"bar,omitempty" hcl:"-"`,
		tags.String(),
	)

	// this shouldn't change anything
	tags.AddOptions("structs", "omitempty")

	require.Equal(t,
		`json:"foo,omitempty" structs:"bar,omitempty" hcl:"-"`,
		tags.String(),
	)

	// this should append to the existing
	tags.AddOptions("structs", "omitnested", "flatten")
	require.Equal(t,
		`json:"foo,omitempty" structs:"bar,omitempty,omitnested,flatten" hcl:"-"`,
		tags.String(),
	)
}

func TestTags_String(t *testing.T) {
	tag := `json:"foo" structs:"bar,omitnested" hcl:"-"`

	tags, err := structtag.Parse(tag)
	require.NoError(t, err)

	require.Equal(t, tag, tags.String())
}

func TestTags_Sort(t *testing.T) {
	tags, err := structtag.Parse(`json:"foo" structs:"bar,omitnested" hcl:"-"`)
	require.NoError(t, err)

	sort.Sort(tags)

	require.Equal(t, `hcl:"-" json:"foo" structs:"bar,omitnested"`, tags.String())
}
