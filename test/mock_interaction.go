package test

import (
	"owl-blogs/domain/model"
	"time"
)

type MockInteractionMetaData struct {
	Str    string
	Number int
	Date   time.Time
}

type MockInteraction struct {
	model.InteractionBase
	metaData *MockInteractionMetaData
}

// Content implements model.Interaction.
func (*MockInteraction) Content() model.InteractionContent {
	return ""
}

// MetaData implements model.Interaction.
func (i *MockInteraction) MetaData() interface{} {
	return i.metaData
}

// SetMetaData implements model.Interaction.
func (i *MockInteraction) SetMetaData(metaData interface{}) {
	i.metaData = metaData.(*MockInteractionMetaData)
}
