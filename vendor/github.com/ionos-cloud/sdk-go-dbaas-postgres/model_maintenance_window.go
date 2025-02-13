/*
 * IONOS DBaaS REST API
 *
 * An enterprise-grade Database is provided as a Service (DBaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.  The API allows you to create additional database clusters or modify existing ones. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 0.0.1
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// MaintenanceWindow A weekly 4 hour-long window, during which maintenance might occur
type MaintenanceWindow struct {
	Time         *string       `json:"time"`
	DayOfTheWeek *DayOfTheWeek `json:"dayOfTheWeek"`
}

// GetTime returns the Time field value
// If the value is explicit nil, the zero value for string will be returned
func (o *MaintenanceWindow) GetTime() *string {
	if o == nil {
		return nil
	}

	return o.Time

}

// GetTimeOk returns a tuple with the Time field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MaintenanceWindow) GetTimeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Time, true
}

// SetTime sets field value
func (o *MaintenanceWindow) SetTime(v string) {

	o.Time = &v

}

// HasTime returns a boolean if a field has been set.
func (o *MaintenanceWindow) HasTime() bool {
	if o != nil && o.Time != nil {
		return true
	}

	return false
}

// GetDayOfTheWeek returns the DayOfTheWeek field value
// If the value is explicit nil, the zero value for DayOfTheWeek will be returned
func (o *MaintenanceWindow) GetDayOfTheWeek() *DayOfTheWeek {
	if o == nil {
		return nil
	}

	return o.DayOfTheWeek

}

// GetDayOfTheWeekOk returns a tuple with the DayOfTheWeek field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *MaintenanceWindow) GetDayOfTheWeekOk() (*DayOfTheWeek, bool) {
	if o == nil {
		return nil, false
	}

	return o.DayOfTheWeek, true
}

// SetDayOfTheWeek sets field value
func (o *MaintenanceWindow) SetDayOfTheWeek(v DayOfTheWeek) {

	o.DayOfTheWeek = &v

}

// HasDayOfTheWeek returns a boolean if a field has been set.
func (o *MaintenanceWindow) HasDayOfTheWeek() bool {
	if o != nil && o.DayOfTheWeek != nil {
		return true
	}

	return false
}

func (o MaintenanceWindow) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Time != nil {
		toSerialize["time"] = o.Time
	}

	if o.DayOfTheWeek != nil {
		toSerialize["dayOfTheWeek"] = o.DayOfTheWeek
	}

	return json.Marshal(toSerialize)
}

type NullableMaintenanceWindow struct {
	value *MaintenanceWindow
	isSet bool
}

func (v NullableMaintenanceWindow) Get() *MaintenanceWindow {
	return v.value
}

func (v *NullableMaintenanceWindow) Set(val *MaintenanceWindow) {
	v.value = val
	v.isSet = true
}

func (v NullableMaintenanceWindow) IsSet() bool {
	return v.isSet
}

func (v *NullableMaintenanceWindow) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMaintenanceWindow(val *MaintenanceWindow) *NullableMaintenanceWindow {
	return &NullableMaintenanceWindow{value: val, isSet: true}
}

func (v NullableMaintenanceWindow) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMaintenanceWindow) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
