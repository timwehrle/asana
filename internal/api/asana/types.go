package asana

import (
	"encoding/json"
	"time"
)

// Date wraps the default time.Time type with appropriate JSON formatting for
// the Asana API when a date is required: '2012-03-26'
type Date time.Time

// Jan 2 15:04:05 2006 MST
const dateLayout = time.DateOnly

// MarshalJSON implements the json.Marshaller interface
func (d *Date) MarshalJSON() ([]byte, error) {
	s := time.Time(*d).Format(dateLayout)
	return json.Marshal(s)
}

// UnmarshalJSON implements the json.Unmarshaller interface
func (d *Date) UnmarshalJSON(value []byte) error {

	var dateString string
	if err := json.Unmarshal(value, &dateString); err != nil {
		return err
	}

	t, err := time.Parse(dateLayout, dateString)
	if err != nil {
		return err
	}

	*d = Date(t)
	return nil
}

// Validator types have a Validate method which is called before posting the
// data to the API
type Validator interface {
	Validate() error
}

// Options - In addition to providing fields and their values in a request,
// you may also specify options to control how your request is interpreted and
// how the response is generated. For GET requests, options are specified as
// URL parameters prefixed with opt_. For POST or PUT requests, options are
// specified in the body, inside the top-level options object (a sibling of
// the data object).
//
// These options can be used in combination in a single request, though some
// of them may conflict in their impact on the response.
type Options struct {
	// Provides the response in “pretty” output. In the case of JSON this
	// means doing proper line breaking and indentation to make it readable.
	// This will take extra time and increase the response size so it is
	// advisable only to use this during debugging.
	Pretty *bool `json:"pretty,omitempty" url:"opt_pretty,omitempty"`

	// In environments that do not support the full range of HTTP verbs, this
	// can be helpful to override the request method sent by the browser,
	// allowing access to the full range of actions in the API. This should
	// only be used when absolutely necessary, as circumventing the browser’s
	// normal interpretation of the HTTP verbs can cause issues. For security
	// reasons, you cannot use this parameter from a GET request, you can only
	// use it to transform a POST request into something else (like a PUT or
	// DELETE).
	Method string `json:"method,omitempty" url:"-"`

	// Some requests return compact representations of objects, to conserve
	// resources and complete the request more efficiently. Other times
	// requests return more information than you may need. This option allows
	// you to list the exact set of fields that the API should be sure to
	// return for the objects. The field names should be provided as paths,
	// described below.   The id of included objects will always be returned,
	// regardless of the field options.
	Fields []string `json:"fields,omitempty" url:"opt_fields,omitempty,comma"`

	// Query results and sub-objects are returned in compact form by default.
	// This option can be used to expand query results or sub-objects to
	// return more detailed information. Be sure you really need the
	// information in the expanded form, as executing a query with many
	// results in expanded form can be costly and return you a lot of data to
	// consume. If the fields option is also used, it will take precedence
	// over the expand option and prevent expansion.
	Expand []string `json:"expand,omitempty" url:"opt_expand,omitempty,comma"`

	// Returns the output in JSON-P format instead of plain JSON, to allow
	// requests to come from within browsers and work around the “same origin
	// policy.” The function named as the value of the opt_jsonp parameter
	// will be called with a single argument, a JavaScript object representing
	// the response.
	JSONP string `json:"jsonp,omitempty" url:"opt_jsonp,omitempty"`

	// The number of objects to return per page. The value must be between 1 and 100.
	Limit int `json:"limit,omitempty" url:"limit,omitempty"`

	// An offset to the next page returned by the API. A pagination request
	// will return an offset token, which can be used as an input parameter to
	// the next request. If an offset is not passed in, the API will return
	// the first page of results.
	//
	// Note: You can only pass in an offset that was returned to you via a
	// previously paginated request.
	Offset string `json:"offset,omitempty" url:"offset,omitempty"`

	// Headers
	Enable  []Feature `json:"-" url:"-"`
	Disable []Feature `json:"-" url:"-"`

	// Filters
	Workspace string `json:"workspace,omitempty" url:"workspace,omitempty"`
	Owner     string `json:"owner,omitempty" url:"owner,omitempty"`

	// Request options
	Debug *bool `json:"-" url:"-"`
}
