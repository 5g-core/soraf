package util

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"strings"

	"gitlabe2.ext.net.nokia.com/5g_core/sor-af/logger"
)

// Serialize - serialize data v to corresponding media type
func Serialize(v interface{}, mediaType string) ([]byte, error) {
	var b []byte
	var err error
	switch KindOfMediaType(mediaType) {
	case MediaKindJSON:
		b, err = json.Marshal(v)
	case MediaKindXML:
		b, err = xml.Marshal(v)
	case MediaKindMultipartRelated:
		b, _, err = MultipartSerialize(v)
	default:
		if err = errors.New("openapi client not supported serialize media type: " + mediaType); err != nil {
			logger.UtilLog.Warnf("Error encode failed: %v", err)
			return nil, err
		}
	}
	return b, err
}

func Deserialize(v interface{}, b []byte, contentType string) (err error) {
	if s, ok := v.(*string); ok {
		*s = string(b)
		return nil
	}

	switch KindOfMediaType(contentType) {
	case MediaKindJSON:
		if err = json.Unmarshal(b, v); err != nil {
			return err
		}
		return nil
	case MediaKindXML:
		if err = xml.Unmarshal(b, v); err != nil {
			return err
		}
		return nil
	case MediaKindMultipartRelated:
		boundary := ""
		for _, part := range strings.Split(contentType, ";") {
			if strings.HasPrefix(part, " boundary=") {
				boundary = part[10:]
			}
		}
		if boundary == "" {
			return errors.New("multipart/related need boundary")
		}
		boundary = strings.Trim(boundary, "\" ")
		if err = MultipartDeserialize(b, v, boundary); err != nil {
			return err
		}
		return nil
	case MediaKindUnsupported:
		return errors.New("undefined response type")
	default:
		if err = errors.New("openapi client not supported serialize media type: " + contentType); err != nil {
			logger.UtilLog.Warnf("Error encode failed: %v", err)
			return err
		}
	}
	return nil
}
