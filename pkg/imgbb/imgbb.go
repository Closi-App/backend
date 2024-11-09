package imgbb

import (
	imgbb "github.com/JohnNON/ImgBB"
	"github.com/spf13/viper"
	"net/http"
)

func NewImgbb(cfg *viper.Viper) *imgbb.Client {
	return imgbb.NewClient(
		&http.Client{
			Timeout: cfg.GetDuration("imgbb.timeout"),
		},
		cfg.GetString("imgbb.api_key"),
	)
}
