package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strconv"
	"strings"
)

func setBody(body interface{}, contentType string) (bodyBuf *bytes.Buffer, err error) {
	if bodyBuf == nil {
		bodyBuf = &bytes.Buffer{}
	}

	if reader, ok := body.(io.Reader); ok {
		_, err = bodyBuf.ReadFrom(reader)
	} else if b, ok := body.([]byte); ok {
		_, err = bodyBuf.Write(b)
	} else if s, ok := body.(string); ok {
		_, err = bodyBuf.WriteString(s)
	} else if s, ok := body.(*string); ok {
		_, err = bodyBuf.WriteString(*s)
	} else {
		b, err = Serialize(body, contentType)
		if err != nil {
			return nil, err
		}
		_, err = bodyBuf.Write(b)
	}

	if err != nil {
		return nil, err
	}

	if bodyBuf.Len() == 0 {
		err = fmt.Errorf("Invalid body type %s\n", contentType)
		return nil, err
	}
	return bodyBuf, nil
}
func MultipartDeserialize(b []byte, v interface{}, boundary string) (err error) {

	body := bytes.NewReader(b)
	r := multipart.NewReader(body, boundary)
	val := reflect.Indirect(reflect.ValueOf(v))

	contentIDIndex := make(map[string]int)

	for {
		var part *multipart.Part
		multipartBody := make([]byte, 1000)

		// if no remian part, break this loop
		if nextPart, err := r.NextPart(); err == io.EOF {
			break
		} else {
			part = nextPart
		}

		contentType := part.Header.Get("Content-Type")
		var n int
		n, err = part.Read(multipartBody)
		if err == nil {
			return
		}
		multipartBody = multipartBody[:n]

		kind := KindOfMediaType(contentType)

		if kind == MediaKindJSON {
			value := val.Field(0)
			if value.Kind() == reflect.Ptr {
				ptr := reflect.New(value.Type().Elem())
				value.Set(ptr)
			}
			if err = json.Unmarshal(multipartBody, value.Interface()); err != nil {
				return err
			}
			structType := val.Type()
			for i := 1; i < structType.NumField(); i++ {
				_, ref, class := parseMultipartFieldParameters(structType.Field(i).Tag.Get("multipart"))
				if ref != "" {
					if contentID, err := getContentID(val, ref, class); err != nil {
						return err
					} else if contentID != "" {
						contentIDIndex[contentID] = i
					}
				}
			}
		} else {
			contentID := part.Header.Get("Content-ID")
			if index, ok := contentIDIndex[contentID]; ok {
				value := val.Field(index)
				value.SetBytes(multipartBody)
			} else {
				return fmt.Errorf("multipart binary data need Content-ID")
			}
		}
	}

	return nil
}

func parseMultipartFieldParameters(str string) (contentType string, ref string, class string) {
	for _, part := range strings.Split(str, ",") {
		switch {
		case strings.HasPrefix(part, "contentType:"):
			contentType = part[12:]
		case strings.HasPrefix(part, "ref:"):
			ref = part[4:]
		case strings.HasPrefix(part, "class:"):
			class = part[6:]
		}
	}
	return
}
func getContentID(v reflect.Value, ref string, class string) (contentID string, err error) {
	recursiveVal := v
	if ref[0] == '{' {
		contentID = ref[1 : len(ref)-1]
		return
	}
	if class != "" {
		var lastVal reflect.Value
		for _, part := range strings.Split(class, ".") {
			lastVal = recursiveVal
			recursiveVal = recursiveVal.FieldByName(part)
			if !recursiveVal.IsValid() {
				return "", fmt.Errorf("Do not have reference field %s in %s for multipart", part, lastVal.Type().String())
			}
			if recursiveVal.Kind() == reflect.Ptr {
				if recursiveVal.IsNil() {
					return "", nil
				}
				recursiveVal = recursiveVal.Elem()
			}
		}
		fieldName := recursiveVal.String()
		if i := strings.IndexRune(fieldName, '-'); i != -1 {
			fieldName = fieldName[:i]
		}
		fieldName = fieldName[:1] + strings.ToLower(fieldName[1:]) + "Info"
		recursiveVal = lastVal.FieldByName(fieldName)
		if recursiveVal.Kind() == reflect.Ptr {
			if recursiveVal.IsNil() {
				return "", nil
			}
			recursiveVal = recursiveVal.Elem()
		}
	}

	for _, part := range strings.Split(ref, ".") {
		lastValType := recursiveVal.Type()
		listIndex := -1

		if start := strings.IndexRune(part, '['); start != -1 {
			end := strings.IndexRune(part, ']')
			listIndex, err = strconv.Atoi(part[start+1 : end])
			if err != nil {
				return "", err
			}
			part = part[:start]
			recursiveVal = recursiveVal.FieldByName(part)
		} else if start = strings.IndexRune(part, '('); start != -1 {
			end := strings.IndexRune(part, ')')
			fieldTypeString := part[start+1 : end]
			var i int
			for i = 0; i < lastValType.NumField(); i++ {
				if fieldType := lastValType.Field(i).Type; strings.HasSuffix(fieldType.String(), fieldTypeString) {
					recursiveVal = recursiveVal.Field(i)
					break
				}
			}
			if i == lastValType.NumField() {
				return "", fmt.Errorf("Do not have reference field Type %s in %s for multipart", fieldTypeString, lastValType.String())
			}
		} else {
			recursiveVal = recursiveVal.FieldByName(part)
		}

		if !recursiveVal.IsValid() {
			return "", fmt.Errorf("Do not have reference field %s in %s for multipart", part, lastValType.String())
		}
		if recursiveVal.Kind() == reflect.Ptr {
			if recursiveVal.IsNil() {
				return "", nil
			}
			recursiveVal = recursiveVal.Elem()
		}
		if listIndex >= 0 {
			if listIndex >= recursiveVal.Len() {
				return "", nil
			}
			recursiveVal = recursiveVal.Index(listIndex)
		}
	}
	contentID = recursiveVal.String()
	return
}
func MultipartEncode(v interface{}, body io.Writer) (string, error) {
	val := reflect.Indirect(reflect.ValueOf(v))
	w := multipart.NewWriter(body)

	structType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsNil() {
			continue
		}
		contentType, ref, class := parseMultipartFieldParameters(structType.Field(i).Tag.Get("multipart"))
		h := make(textproto.MIMEHeader)

		if ref != "" {
			if contentID, err := getContentID(val, ref, class); err != nil {
				return "", err
			} else if contentID != "" {
				h.Set("Content-ID", contentID)
			} else {
				return "", fmt.Errorf("ContentID of multipart binaryData in JsonData is empty")
			}
		}
		h.Set("Content-Type", contentType)
		fieldData, err := w.CreatePart(h)
		if err != nil {
			return "", err
		}
		fieldBody, err := setBody(val.Field(i).Interface(), contentType)
		if err != nil {
			return "", err
		}
		_, err = fieldData.Write(fieldBody.Bytes())
		if err != nil {
			return "", err
		}
	}
	err := w.Close()
	if err != nil {
		return "", err
	}
	contentType := "multipart/related; boundary=\"" + w.Boundary() + "\""

	return contentType, nil
}

func MultipartSerialize(v interface{}) ([]byte, string, error) {
	buffer := new(bytes.Buffer)
	val := reflect.Indirect(reflect.ValueOf(v))
	w := multipart.NewWriter(buffer)

	structType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).IsNil() {
			continue
		}
		contentType, ref, class := parseMultipartFieldParameters(structType.Field(i).Tag.Get("multipart"))
		h := make(textproto.MIMEHeader)

		if ref != "" {
			if contentID, err := getContentID(val, ref, class); err != nil {
				return nil, "", err
			} else if contentID != "" {
				h.Set("Content-ID", contentID)
			} else {
				return nil, "", fmt.Errorf("ContentID of multipart binaryData in JsonData is empty")
			}
		}
		h.Set("Content-Type", contentType)
		fieldData, err := w.CreatePart(h)
		if err != nil {
			return nil, "", err
		}
		fieldBody, err := setBody(val.Field(i).Interface(), contentType)
		if err != nil {
			return nil, "", err
		}
		_, err = fieldData.Write(fieldBody.Bytes())
		if err != nil {
			return nil, "", err
		}
	}
	err := w.Close()
	if err != nil {
		return nil, "", err
	}
	contentType := "multipart/related; boundary=\"" + w.Boundary() + "\""

	return buffer.Bytes(), contentType, nil
}
