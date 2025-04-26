package asana

import (
	"fmt"
	"time"
)

type EnumValue struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	EnumValueBase

	Enabled bool `json:"enabled"`
}

type EnumValueBase struct {
	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// The color of the enum option. Defaults to ‘none’.
	Color string `json:"color,omitempty"`
}

type DateValue struct {
	Date     *Date      `json:"date,omitempty"`
	DateTime *time.Time `json:"date_time,omitempty"`
}

type FieldType string

// FieldTypes for CustomField.Type field
const (
	FieldTypeText   FieldType = "text"
	FieldTypeNumber FieldType = "number"
	FieldTypeEnum   FieldType = "enum"

	// https://forum.asana.com/t/new-custom-field-types/32244
	FieldTypeMultiEnum FieldType = "multi_enum"
	FieldTypeDate      FieldType = "date"
	FieldTypeBoolean   FieldType = "boolean"

	// https://forum.asana.com/t/introducing-people-custom-field/187166
	FieldTypePeople FieldType = "people"

	Text   FieldType = FieldTypeText   // Deprecated - use FieldTypeText
	Number FieldType = FieldTypeNumber // Deprecated - use FieldTypeNumber
	Enum   FieldType = FieldTypeEnum   // Deprecated - use FieldTypeEnum
)

type LabelPosition string

const (
	Prefix LabelPosition = "prefix"
	Suffix LabelPosition = "suffix"
)

type Format string

const (
	Currency   = "currency"
	Identifier = "identifier"
	Percentage = "percentage"
	Custom     = "custom"
	None       = "none"
)

type CustomFieldBase struct {
	// ISO 4217 currency code to format this custom field. This will be null
	// if the format is not currency.
	CurrencyCode string `json:"currency_code,omitempty"`

	// This is the string that appears next to the custom field value.
	// This will be null if the format is not custom.
	CustomLabel string `json:"custom_label,omitempty"`

	// Only relevant for custom fields with custom format. This depicts where to place the custom label.
	// This will be null if the format is not custom.
	CustomLabelPosition LabelPosition `json:"custom_label_position,omitempty"`

	// The description of the custom field.
	Description string `json:"description,omitempty"`

	// The format of this custom field.
	Format Format `json:"format,omitempty"`

	// Conditional. This flag describes whether a follower of a task with this
	// field should receive inbox notifications from changes to this field.
	HasNotificationsEnabled *bool `json:"has_notifications_enabled,omitempty"`

	// Read-only. The name of the object.
	Name string `json:"name,omitempty"`

	// Only relevant for custom fields of type ‘Number’. This field dictates the number of places after
	// the decimal to round to, i.e. 0 is integer values, 1 rounds to the nearest tenth, and so on.
	// Must be between 0 and 6, inclusive.
	// For percentage format, this may be unintuitive, as a value of 0.25 has a precision of 0, while a
	// value of 0.251 has a precision of 1. This is due to 0.25 being displayed as 25%.
	// The identifier format will always have a precision of 0.
	Precision *int `json:"precision,omitempty"`

	// The type of the custom field. Must be one of the given values:
	// 'text', 'enum', 'number', 'multi_enum', 'date', 'boolean'
	ResourceSubtype FieldType `json:"resource_subtype"`
}

// Custom Fields store the metadata that is used in order to add user-
// specified information to tasks in Asana. Be sure to reference the Custom
// Fields developer documentation for more information about how custom fields
// relate to various resources in Asana.
// Users in Asana can lock custom fields, which will make them read-only when accessed by other users.
// Attempting to edit a locked custom field will return HTTP error code 403 Forbidden.
type CustomField struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	CustomFieldBase

	// Only relevant for custom fields of type ‘Enum’. This array specifies
	// the possible values which an enum custom field can adopt.
	EnumOptions []*EnumValue `json:"enum_options,omitempty"`

	// Read-only: This flag describes whether this custom field is available to every container
	// in the workspace. Before project-specific custom fields, this field was always true.
	IsGlobalToWorkspace *bool `json:"is_global_to_workspace,omitempty"`

	// Determines whether the custom field is available for editing on a task
	// (i.e. the field is associated with one of the task's parent projects)
	Enabled *bool `json:"enabled,omitempty"`
}

type CustomFieldSetting struct {
	// Read-only. Globally unique ID of the object
	ID string `json:"gid,omitempty"`

	CustomField *CustomField `json:"custom_field"`

	Project *Project `json:"project,omitempty"`

	Important bool `json:"is_important,omitempty"`
}

type AddCustomFieldSettingRequest struct {
	CustomField  string `json:"custom_field"`
	Important    bool   `json:"is_important,omitempty"`
	InsertBefore string `json:"insert_before,omitempty"`
	InsertAfter  string `json:"insert_after,omitempty"`
}

func (p *Project) AddCustomFieldSetting(client *Client, request *AddCustomFieldSettingRequest) (*CustomFieldSetting, error) {
	client.trace("Attach custom field %q to project %q", request.CustomField, p.ID)

	// Custom request encoding
	m := map[string]interface{}{}
	m["custom_field"] = request.CustomField
	m["is_important"] = request.Important

	if request.InsertAfter == "-" {
		m["insert_after"] = nil
	} else if request.InsertAfter != "" {
		m["insert_after"] = request.InsertAfter
	}

	if request.InsertBefore == "-" {
		m["insert_before"] = nil
	} else if request.InsertBefore != "" {
		m["insert_before"] = request.InsertBefore
	}

	result := &CustomFieldSetting{}
	err := client.post(fmt.Sprintf("/projects/%s/addCustomFieldSetting", p.ID), m, result)
	return result, err
}

func (p *Project) RemoveCustomFieldSetting(client *Client, customFieldID string) error {
	client.trace("Remove custom field %q from project %q", customFieldID, p.ID)

	// Custom request encoding
	m := map[string]interface{}{
		"custom_field": customFieldID,
	}

	err := client.post(fmt.Sprintf("/projects/%s/removeCustomFieldSetting", p.ID), m, nil)
	return err
}

type ProjectLocalCustomField struct {
	CustomFieldBase

	// Only relevant for custom fields of type ‘Enum’. This array specifies
	// the possible values which an enum custom field can adopt.
	EnumOptions []*EnumValueBase `json:"enum_options,omitempty"`
}

type AddProjectLocalCustomFieldRequest struct {
	CustomField  ProjectLocalCustomField `json:"custom_field"`
	Important    bool                    `json:"is_important,omitempty"`
	InsertBefore string                  `json:"insert_before,omitempty"`
	InsertAfter  string                  `json:"insert_after,omitempty"`
}

func (p *Project) AddProjectLocalCustomField(client *Client, request *AddProjectLocalCustomFieldRequest) (*CustomFieldSetting, error) {
	client.trace("Attach custom field %q to project %q", request.CustomField.Name, p.ID)

	// Custom request encoding
	m := map[string]interface{}{}
	m["custom_field"] = request.CustomField
	m["is_important"] = request.Important

	if request.InsertAfter == "-" {
		m["insert_after"] = nil
	} else if request.InsertAfter != "" {
		m["insert_after"] = request.InsertAfter
	}

	if request.InsertBefore == "-" {
		m["insert_before"] = nil
	} else if request.InsertBefore != "" {
		m["insert_before"] = request.InsertBefore
	}

	result := &CustomFieldSetting{}
	err := client.post(fmt.Sprintf("/projects/%s/addCustomFieldSetting", p.ID), m, result)
	return result, err
}

type CreateCustomFieldRequest struct {
	CustomFieldBase

	// Required: The workspace to create a custom field in.
	Workspace string `json:"workspace"`

	// Conditional. Only relevant for custom fields of type enum.
	// This array specifies the possible values which an enum custom field can adopt.
	EnumOptions []*EnumValueBase `json:"enum_options,omitempty"`
}

func (c *Client) CreateCustomField(request *CreateCustomFieldRequest) (*CustomField, error) {
	c.trace("Create custom field %q in workspace %s", request.Name, request.Workspace)

	result := &CustomField{}
	err := c.post("/custom_fields", request, result)
	return result, err
}

// When a custom field is associated with a project, tasks in that project can
// carry additional custom field values which represent the value of the field
// on that particular task - for instance, the selected item from an enum type
// custom field. These custom fields will appear as an array in a
// custom_fields property of the task, along with some basic information which
// can be used to associate the custom field value with the custom field
// metadata.
type CustomFieldValue struct {
	CustomField

	// A string representation for the value of the custom field. Integrations
	// that don't require the underlying type should use this field to read values.
	// Using this field will future-proof an app against new custom field types.
	DisplayValue *string `json:"display_value,omitempty"`

	// Custom fields of type text will return a text_value property containing
	// the string of text for the field.
	TextValue *string `json:"text_value,omitempty"`

	// Custom fields of type number will return a number_value property
	// containing the number for the field.
	NumberValue *float64 `json:"number_value,omitempty"`

	// Conditional. Only relevant for custom fields of type boolean.
	BooleanValue *bool `json:"boolean_value,omitempty"`

	// Conditional. Only relevant for custom fields of type date.
	DateValue *DateValue `json:"date_value,omitempty"`

	// Conditional. Only relevant for custom fields of type enum.
	// This object is the chosen value of an enum custom field.
	EnumValue *EnumValue `json:"enum_value,omitempty"`

	// Conditional. Only relevant for custom fields of type multi_enum.
	// This object is the chosen values of a multi_enum custom field.
	MultiEnumValues []*EnumValue `json:"multi_enum_values,omitempty"`

	// Conditional. Only relevant for custom fields of type people.
	// This object is the chosen values of a people custom field.
	PeopleValue []*User `json:"people_value,omitempty"`
}

// Fetch loads the full details for this CustomField
func (f *CustomField) Fetch(client *Client, options ...*Options) error {
	client.trace("Loading details for custom field %q", f.ID)

	_, err := client.get(fmt.Sprintf("/custom_fields/%s", f.ID), nil, f, options...)
	return err
}

// CustomFields returns the compact records for all custom fields in the workspace
func (w *Workspace) CustomFields(client *Client, options ...*Options) ([]*CustomField, *NextPage, error) {
	client.trace("Listing custom fields in workspace %s...\n", w.ID)
	var result []*CustomField

	// Make the request
	nextPage, err := client.get(fmt.Sprintf("/workspaces/%s/custom_fields", w.ID), nil, &result, options...)
	return result, nextPage, err
}

// AllCustomFields repeatedly pages through all available custom fields in a workspace
func (w *Workspace) AllCustomFields(client *Client, options ...*Options) ([]*CustomField, error) {
	var allCustomFields []*CustomField
	nextPage := &NextPage{}

	var customFields []*CustomField
	var err error

	for nextPage != nil {
		page := &Options{
			Limit:  50,
			Offset: nextPage.Offset,
		}

		allOptions := append([]*Options{page}, options...)
		customFields, nextPage, err = w.CustomFields(client, allOptions...)
		if err != nil {
			return nil, err
		}

		allCustomFields = append(allCustomFields, customFields...)
	}
	return allCustomFields, nil
}
