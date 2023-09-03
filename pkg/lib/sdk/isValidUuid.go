package sdk

import "github.com/segmentio/ksuid"

func IsValidUuid(u string) bool {
	_, err := ksuid.Parse(u)
	return err == nil
}
