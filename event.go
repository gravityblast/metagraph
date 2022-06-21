package main

import (
	"fmt"
	"math/big"
	"reflect"
)

type InvalidFieldErr struct {
	Message string
	Field   string
}

func newInvalidFieldErr(msg string, field string) *InvalidFieldErr {
	return &InvalidFieldErr{
		Message: msg,
		Field:   field,
	}
}

func (e *InvalidFieldErr) Error() string {
	return fmt.Sprintf("%s. invalid field %s", e.Message, e.Field)
}

type MetaPointer struct {
	Protocol string
	Pointer  string
}

func parseEventStruct(config *EventConfig, eventMap map[string]interface{}) (*MetaPointer, error) {
	raw := eventMap[config.MetaPointerField]
	if raw == nil {
		return nil, newInvalidFieldErr("cannot find metaPointer field", config.MetaPointerField)
	}

	value := reflect.ValueOf(raw)

	protocol := value.FieldByName(config.ProtocolField)
	pointer := value.FieldByName(config.PointerField)

	if protocol.Kind() == reflect.Invalid {
		return nil, newInvalidFieldErr("cannot find protocol field", config.ProtocolField)
	}

	if pointer.Kind() == reflect.Invalid {
		return nil, newInvalidFieldErr("cannot find pointer field", config.PointerField)
	}

	return &MetaPointer{
		Protocol: fmt.Sprintf("%s", protocol),
		Pointer:  pointer.String(),
	}, nil
}

func parseEventMap(config *EventConfig, eventMap map[string]interface{}) (*MetaPointer, error) {
	var protocol *big.Int
	if config.ProtocolField != "" {
		var okProtocol bool
		protocol, okProtocol = eventMap[config.ProtocolField].(*big.Int)
		if !okProtocol {
			return nil, newInvalidFieldErr("cannot find protocol field", config.ProtocolField)
		}
	}

	pointer, okPointer := eventMap[config.PointerField].(string)
	if !okPointer {
		return nil, newInvalidFieldErr("cannot find pointer field", config.PointerField)
	}

	return &MetaPointer{
		Protocol: protocol.String(),
		Pointer:  pointer,
	}, nil
}

func parseEvent(config *EventConfig, data []byte) (*MetaPointer, error) {
	eventMap := make(map[string]interface{})
	err := config.ABI.UnpackIntoMap(eventMap, config.EventName, data)
	if err != nil {
		return nil, err
	}

	if config.MetaPointerField != "" {
		return parseEventStruct(config, eventMap)
	} else {
		return parseEventMap(config, eventMap)
	}
}
