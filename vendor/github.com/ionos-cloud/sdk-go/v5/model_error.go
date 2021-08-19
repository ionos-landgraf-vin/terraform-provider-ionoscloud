/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// Error struct for Error
type Error struct {
	// HTTP status code of the operation
	HttpStatus *int32 `json:"httpStatus,omitempty"`
	Messages *[]ErrorMessage `json:"messages,omitempty"`
}



// GetHttpStatus returns the HttpStatus field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *Error) GetHttpStatus() *int32 {
	if o == nil {
		return nil
	}


	return o.HttpStatus

}

// GetHttpStatusOk returns a tuple with the HttpStatus field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Error) GetHttpStatusOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}


	return o.HttpStatus, true
}

// SetHttpStatus sets field value
func (o *Error) SetHttpStatus(v int32) {


	o.HttpStatus = &v

}

// HasHttpStatus returns a boolean if a field has been set.
func (o *Error) HasHttpStatus() bool {
	if o != nil && o.HttpStatus != nil {
		return true
	}

	return false
}



// GetMessages returns the Messages field value
// If the value is explicit nil, the zero value for []ErrorMessage will be returned
func (o *Error) GetMessages() *[]ErrorMessage {
	if o == nil {
		return nil
	}


	return o.Messages

}

// GetMessagesOk returns a tuple with the Messages field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Error) GetMessagesOk() (*[]ErrorMessage, bool) {
	if o == nil {
		return nil, false
	}


	return o.Messages, true
}

// SetMessages sets field value
func (o *Error) SetMessages(v []ErrorMessage) {


	o.Messages = &v

}

// HasMessages returns a boolean if a field has been set.
func (o *Error) HasMessages() bool {
	if o != nil && o.Messages != nil {
		return true
	}

	return false
}


func (o Error) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.HttpStatus != nil {
		toSerialize["httpStatus"] = o.HttpStatus
	}
	

	if o.Messages != nil {
		toSerialize["messages"] = o.Messages
	}
	
	return json.Marshal(toSerialize)
}

type NullableError struct {
	value *Error
	isSet bool
}

func (v NullableError) Get() *Error {
	return v.value
}

func (v *NullableError) Set(val *Error) {
	v.value = val
	v.isSet = true
}

func (v NullableError) IsSet() bool {
	return v.isSet
}

func (v *NullableError) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableError(val *Error) *NullableError {
	return &NullableError{value: val, isSet: true}
}

func (v NullableError) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableError) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


