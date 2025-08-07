package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var ErrInvalidCode = errors.New("invalid invite code format")

func parseIntQueryParam(r *http.Request, key string, defaultVal int) (int, error) {
	valStr := r.URL.Query().Get(key)
	if valStr == "" {
		return defaultVal, nil
	}
	i, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal, err
	}
	return i, nil
}

func parseUUIDPathParam(r *http.Request, key string) (uuid.UUID, error) {
	valStr := mux.Vars(r)[key]
	if valStr == "" {
		return uuid.Nil, ErrInvalidID
	}
	id, err := uuid.Parse(valStr)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func parseInviteCodePathParam(r *http.Request, key string) (string, error) {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	valStr := strings.ToUpper(mux.Vars(r)[key])
	if valStr == "" {
		return "", ErrInvalidCode
	}
	if len(valStr) != 9 {
		return "", ErrInvalidCode
	}

	for _, char := range valStr {
		if !strings.ContainsRune(charset, char) {
			return "", ErrInvalidCode
		}
	}

	log.Printf("Invite code raw value: '%s' (len=%d)", valStr, len(valStr))

	return valStr, nil
}
