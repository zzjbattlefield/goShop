package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshlJson() ([]byte, error) {
	stmp := fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32  `json:"id,omitempty"`
	Mobile   string `json:"mobile,omitempty"`
	NickName string `json:"name,omitempty"`
	// BirthDay string `json:"birthDay,omitempty"`
	BirthDay JsonTime `json:"birthDay,omitempty"`
	Gender   string   `json:"gender,omitempty"`
}
