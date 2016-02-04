package query

import (
	"errors"

	"gopkg.in/bluesuncorp/validator.v8"
)

// Set of query types we expect to receive.
const (
	TypePipeline = "pipeline"
	TypeTemplate = "template"
)

//==============================================================================

// validate is used to perform model field validation.
var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

//==============================================================================

// Result contains the result of an query set execution.
type Result struct {
	Results interface{} `json:"results"`
	Error   bool        `json:"error"`
}

//==============================================================================

// Query contains the configuration details for a query.
type Query struct {
	Name        string                   `bson:"name" json:"name" validate:"required,min=3"`                                 // Unique name per query document.
	Description string                   `bson:"desc,omitempty" json:"desc,omitempty"`                                       // Description of this specific query.
	Type        string                   `bson:"type" json:"type" validate:"required,min=8"`                                 // TypePipeline, TypeTemplate
	Collection  string                   `bson:"collection,omitempty" json:"collection,omitempty" validate:"required,min=3"` // Name of the collection to use for processing the query.
	Commands    []map[string]interface{} `bson:"commands" json:"commands"`                                                   // Commands to process for the query.
	Continue    bool                     `bson:"continue,omitempty" json:"continue,omitempty"`                               // Indicates that on failure to process the next query.
	Return      bool                     `bson:"return" json:"return"`                                                       // Return the results back to the user with Name as the key.
}

// Validate checks the query value for consistency.
func (q *Query) Validate() error {
	if err := validate.Struct(q); err != nil {
		return err
	}

	if len(q.Commands) == 0 {
		return errors.New("No commands exist")
	}

	switch q.Type {
	case TypePipeline:
		// Place holder since things are good.

	case TypeTemplate:
		if len(q.Commands) > 1 {
			return errors.New("Invalid number of commands")
		}

	default:
		return errors.New("Invalid query type")
	}

	return nil
}

//==============================================================================

// Param contains meta-data about a required parameter for the query.
type Param struct {
	Name      string `bson:"name" json:"name"`             // Name of the parameter.
	Desc      string `bson:"desc" json:"desc"`             // Description about the parameter.
	Default   string `bson:"default" json:"default"`       // Default value for the parameter.
	RegexName string `bson:"regex_name" json:"regex_name"` // Regular expression name.
}

//==============================================================================

// Set contains the configuration details for a rule set.
type Set struct {
	Name        string  `bson:"name" json:"name" validate:"required,min=3"` // Name of the query set.
	Description string  `bson:"desc" json:"desc"`                           // Description of the query set.
	PreScript   string  `bson:"pre_script" json:"pre_script"`               // Name of a script document to prepend.
	PstScript   string  `bson:"pst_script" json:"pst_script"`               // Name of a script document to append.
	Params      []Param `bson:"params" json:"params"`                       // Collection of parameters.
	Queries     []Query `bson:"queries" json:"queries"`                     // Collection of queries.
	Enabled     bool    `bson:"enabled" json:"enabled"`                     // If the query set is enabled to run.
}

// Validate checks the set value for consistency.
func (s *Set) Validate() error {
	if err := validate.Struct(s); err != nil {
		return err
	}

	for _, q := range s.Queries {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// PrepareForInsert replaces the documents for insertion.
func (s *Set) PrepareForInsert() {

	// Fix the commands so it can be inserted.
	for q := range s.Queries {
		for c := range s.Queries[q].Commands {
			prepareForInsert(s.Queries[q].Commands[c])
		}
	}
}

// PrepareForUse replaces the documents back to their orginal form.
func (s *Set) PrepareForUse() {

	// Fix the commands so things are back to their orginal form.
	for q := range s.Queries {
		for c := range s.Queries[q].Commands {
			prepareForUse(s.Queries[q].Commands[c])
		}
	}
}
