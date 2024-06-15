package windIne_orm_config

import "fmt"

type WindIneTimeZone int64

const (
	WindIneORMTimeZoneUTC WindIneTimeZone = iota
	WindIneORMTimeZoneShangHai
)

func (timeZone WindIneTimeZone) String() string {
	switch timeZone {
	case WindIneORMTimeZoneShangHai:
		return fmt.Sprintf("Asia%sShanghai", "%2F")
	case WindIneORMTimeZoneUTC:
		return "UTC"
	default:
		return "UTC"
	}
}
