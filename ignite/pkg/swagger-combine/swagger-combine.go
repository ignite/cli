package swaggercombine

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-openapi/analysis"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/spec"

	"github.com/ignite/cli/v29/ignite/pkg/errors"
)

// Config represent swagger-combine config.
type Config struct {
	spec  *spec.Swagger
	specs []*spec.Swagger
}

// New create a mew swagger combine config.
func New(title, name string) *Config {
	return &Config{
		spec: &spec.Swagger{
			SwaggerProps: spec.SwaggerProps{
				ID:      name,
				Swagger: "2.0",
				Info: &spec.Info{
					InfoProps: spec.InfoProps{
						Description: fmt.Sprintf("Chain %s REST API", name),
						Title:       title,
						Contact:     &spec.ContactInfo{ContactInfoProps: spec.ContactInfoProps{Name: name}},
					},
				},
				Definitions: make(spec.Definitions),
			},
		},
		specs: make([]*spec.Swagger, 0),
	}
}

// AddSpec adds a new OpenAPI spec to Config by path in the fs and unique id of spec.
func (c *Config) AddSpec(id, path string, makeUnique bool) error {
	baseDoc, err := loads.Spec(path)
	if err != nil {
		return errors.Wrapf(err, "failed to load spec from path %s", path)
	}

	spec := baseDoc.Spec()
	if makeUnique {
		for i, specPath := range spec.Paths.Paths {
			if specPath.Get != nil {
				specPath.Get.ID = id + specPath.Get.ID
			}
			if specPath.Post != nil {
				specPath.Post.ID = id + specPath.Post.ID
			}
			spec.Paths.Paths[i] = specPath
		}
	}

	c.specs = append(c.specs, c.mergeTags(c.mergeDefinitions(spec)))

	return nil
}

// mergeDefinitions merge spec definitions with main spec and erase the spec definition.
func (c *Config) mergeDefinitions(m *spec.Swagger) *spec.Swagger {
	for k, v := range m.Definitions {
		if _, exists := c.spec.Definitions[k]; exists {
			continue
		}
		c.spec.Definitions[k] = v
	}
	m.Definitions = nil
	return m
}

// mergeTags merge spec tags with main spec and erase the spec tag.
func (c *Config) mergeTags(m *spec.Swagger) *spec.Swagger {
	for _, v := range m.Tags {
		found := false
		for _, vv := range c.spec.Tags {
			if v.Name == vv.Name {
				found = true
				break
			}
		}
		if found {
			continue
		}
		c.spec.Tags = append(c.spec.Tags, v)
	}
	m.Tags = nil
	return m
}

// Combine combines openapi specs into one and saves to out path.
func (c *Config) Combine(out string) error {
	sort.Slice(c.specs, func(a, b int) bool { return c.specs[a].ID < c.specs[b].ID })

	errs := analysis.Mixin(c.spec, c.specs...)
	if len(errs) > 0 {
		return errors.Errorf("invalid mix specs: %s", strings.Join(errs, ", "))
	}
	specJSON, err := c.spec.MarshalJSON()
	if err != nil {
		return err
	}
	// ensure out dir exists.
	outDir := filepath.Dir(out)
	if err := os.MkdirAll(outDir, 0o766); err != nil {
		return err
	}
	if err = os.WriteFile(out, specJSON, 0o600); err != nil {
		return errors.Wrapf(err, "failed to write combined spec to file %s", out)
	}
	return nil
}
