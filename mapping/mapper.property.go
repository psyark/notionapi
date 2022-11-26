package mapping

import "github.com/psyark/notionapi"

type propertyMapper struct {
	propID string
}

func (m *propertyMapper) getPropValue(page notionapi.Page) *notionapi.PropertyValue {
	return page.Properties.Get(m.propID)
}
