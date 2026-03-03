package swaggercombine

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-openapi/spec"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	c := New("My API", "ignite")
	require.NotNil(t, c.spec)
	require.Equal(t, "ignite", c.spec.ID)
	require.Equal(t, "2.0", c.spec.Swagger)
	require.Equal(t, "My API", c.spec.Info.Title)
}

func TestMergeDefinitionsAndTags(t *testing.T) {
	c := New("My API", "ignite")
	in := &spec.Swagger{
		SwaggerProps: spec.SwaggerProps{
			Definitions: spec.Definitions{
				"MyType": spec.Schema{
					SchemaProps: spec.SchemaProps{
						Type: spec.StringOrArray{"object"},
					},
				},
			},
			Tags: []spec.Tag{
				{TagProps: spec.TagProps{Name: "tag-a"}},
			},
		},
	}

	out := c.mergeDefinitions(in)
	require.Nil(t, out.Definitions)
	require.Contains(t, c.spec.Definitions, "MyType")

	out = c.mergeTags(out)
	require.Nil(t, out.Tags)
	require.Len(t, c.spec.Tags, 1)
	require.Equal(t, "tag-a", c.spec.Tags[0].Name)
}

func TestAddSpecAndCombine(t *testing.T) {
	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.json")
	specJSON := `{
  "swagger":"2.0",
  "info":{"title":"A","version":"1.0"},
  "paths":{
    "/hello":{
      "get":{"operationId":"GetHello","responses":{"200":{"description":"ok"}}}
    }
  }
}`
	require.NoError(t, os.WriteFile(specPath, []byte(specJSON), 0o600))

	c := New("My API", "ignite")
	require.NoError(t, c.AddSpec("mod1-", specPath, true))

	outPath := filepath.Join(dir, "combined", "swagger.json")
	require.NoError(t, c.Combine(outPath))

	raw, err := os.ReadFile(outPath)
	require.NoError(t, err)

	var out map[string]any
	require.NoError(t, json.Unmarshal(raw, &out))
	require.Equal(t, "2.0", out["swagger"])
	paths, ok := out["paths"].(map[string]any)
	require.True(t, ok)
	require.Contains(t, paths, "/hello")
}
