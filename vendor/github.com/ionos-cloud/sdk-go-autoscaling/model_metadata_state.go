/*
 * VM Auto Scaling Service (part of CloudAPI)
 *
 * Provides Endpoints to manage the Autoscaling feature by IONOS cloud
 *
 * API version: 1-SDK.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"fmt"
)

// MetadataState The state of the resource
type MetadataState string

// List of MetadataState
const (
	AVAILABLE MetadataState = "AVAILABLE"
	BUSY MetadataState = "BUSY"
	INACTIVE MetadataState = "INACTIVE"
	SUSPENDED MetadataState = "SUSPENDED"
)

func (v *MetadataState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := MetadataState(value)
	for _, existing := range []MetadataState{ "AVAILABLE", "BUSY", "INACTIVE", "SUSPENDED",   } {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid MetadataState", value)
}

// Ptr returns reference to MetadataState value
func (v MetadataState) Ptr() *MetadataState {
	return &v
}

type NullableMetadataState struct {
	value *MetadataState
	isSet bool
}

func (v NullableMetadataState) Get() *MetadataState {
	return v.value
}

func (v *NullableMetadataState) Set(val *MetadataState) {
	v.value = val
	v.isSet = true
}

func (v NullableMetadataState) IsSet() bool {
	return v.isSet
}

func (v *NullableMetadataState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMetadataState(val *MetadataState) *NullableMetadataState {
	return &NullableMetadataState{value: val, isSet: true}
}

func (v NullableMetadataState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMetadataState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

