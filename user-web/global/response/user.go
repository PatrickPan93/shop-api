package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	strTime := fmt.Sprintf("\"%s\"", time.Time(j).Format("1993-08-30"))
	return []byte(strTime), nil
}

type UserResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
