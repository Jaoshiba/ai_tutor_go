
package entities

import "time"

type SessionDataModel struct {
    SessionID string    `json:"sessionid"`
    UserID    string    `json:"userid"`
    CreatedAt time.Time `json:"createdat"`
    ExpiredAt time.Time `json:"expiredat"`
}
