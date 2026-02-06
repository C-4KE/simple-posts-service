package cursor

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func Create(commentID int64, parentPath string) string {
	return encodeCursor(strings.Join([]string{parentPath, strconv.FormatInt(commentID, 10)}, "."))
}

func Validate(cursor string) error {
	cursorString, err := decodeCursor(cursor)
	if err != nil {
		return err
	}

	cursorString = strings.ReplaceAll(cursorString, ".", "")
	for _, char := range cursorString {
		if !unicode.IsDigit(char) {
			return errors.New("Cursor " + cursor + " is not valid.")
		}
	}

	return nil
}

func GetPath(cursor string) (string, error) {
	cursorString, err := decodeCursor(cursor)
	if err != nil {
		return "", err
	}

	lastDotIndex := strings.LastIndex(cursorString, ".")
	if lastDotIndex == -1 {
		return "", errors.New("Cursor " + cursor + " is not valid.")
	}

	return cursorString[:lastDotIndex], nil
}

func GetCommentID(cursor string) (string, error) {
	cursorString, err := decodeCursor(cursor)
	if err != nil {
		return "", err
	}

	lastDotIndex := strings.LastIndex(cursorString, ".")
	if lastDotIndex == -1 {
		return "", errors.New("Cursor " + cursor + " is not valid.")
	}

	return cursorString[lastDotIndex:], nil
}

func encodeCursor(cursorString string) string {
	return base64.RawStdEncoding.EncodeToString([]byte(cursorString))
}

func decodeCursor(rawCursor string) (string, error) {
	byteCursor, err := base64.RawStdEncoding.DecodeString(rawCursor)
	return string(byteCursor), err
}
