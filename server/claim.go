package server

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// getClaim and convert it to map[string]interface{}
func claimToMap(claim interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	v := reflect.ValueOf(claim)
	if v.Kind() == reflect.Map {
		for _, key := range v.MapKeys() {
			strct := v.MapIndex(key)
			// fmt.Println(key.Interface(), strct.Interface())
			m[fmt.Sprintf("%v", key.Interface())] = strct.Interface()
		}
	}

	return m
}

// ClaimGetID - extract id from jwt token claim
func ClaimGetID(c interface{}) (uint, error) {
	claim := claimToMap(c)
	// check if claim has id or not
	if val, ok := claim["id"]; ok {
		i, err := strconv.Atoi(fmt.Sprintf("%v", val))
		return uint(i), err
	}

	// return error if claim has no ID key
	return 0, errors.New("Can not extract id from JWT Token Claim")
}

// ClaimGetUsername - extract username from jwt token claim
func ClaimGetUsername(c interface{}) (string, error) {
	claim := claimToMap(c)
	// check if claim has id or not
	if val, ok := claim["username"]; ok {
		return fmt.Sprintf("%v", val), nil
	}

	// return error if claim has no ID key
	return "", errors.New("Can not extract id from JWT Token Claim")
}
