package pwdless

import (
	"go.uber.org/zap"
	"time"

	"github.com/dhax/go-base/logging"
)

func (rs *Resource) choresTicker() {
	ticker := time.NewTicker(time.Hour * 1)
	go func() {
		for range ticker.C {
			if err := rs.Store.PurgeExpiredToken(); err != nil {
				logging.Logger.With(zap.String("chore", "purgeExpiredToken")).Error(err.Error())
			}
		}
	}()
}
