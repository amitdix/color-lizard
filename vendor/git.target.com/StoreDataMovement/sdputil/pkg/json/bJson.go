package json

import (
	ej "encoding/json"
	"errors"
)

type JSON map[string]interface{}

func Parse(s string) (JSON, error) {
	var raw JSON
	err := ej.Unmarshal([]byte(s), &raw)
	return raw, err
}

func (b JSON) GetJson(key string) (JSON, error) {
	val, ok0 := b[key]
	if !ok0 {
		return nil, errors.New("bJson object doesn't have key : " + key)
	}
	j, ok1 := val.(map[string]interface{})
	if !ok1 {
		return nil, errors.New("Couldn't parse value for \"" + key + "\" to json")
	}
	return j, nil
}

func (b JSON) GetString(key string) (string, error) {
	val, ok0 := b[key]
	if !ok0 {
		return "", errors.New("bJson object doesn't have key : " + key)
	}
	j, ok1 := val.(string)
	if !ok1 {
		return "", errors.New("bJson member at key \"" + key + "\" isn't string")
	}
	return j, nil
}

func (b JSON) GetFloat(key string) (float64, error) {
	val, ok0 := b[key]
	if !ok0 {
		return 0.0, errors.New("bJson object doesn't have key : " + key)
	}
	j, ok1 := val.(float64)
	if !ok1 {
		return 0.0, errors.New("bJson member at key \"" + key + "\" isn't float")
	}
	return j, nil
}

func (b JSON) GetInt(key string) (int, error) {
	f, e := b.GetFloat(key)
	if e != nil {
		return 0.0, e
	}
	i := int(f)
	if float64(i) != f {
		return 0.0, errors.New("bJson member at key \"" + key + "\" isn't int")
	}
	return i, nil
}

func (b JSON) GetInt64(key string) (int64, error) {
	f, e := b.GetFloat(key)
	if e != nil {
		return 0.0, e
	}
	i := int64(f)
	if float64(i) != f {
		return 0.0, errors.New("bJson member at key \"" + key + "\" isn't int")
	}
	return i, nil
}
