package pixiv

import (
	"fmt"
	"testing"
)

func TestAppPixivAPI(t *testing.T) {
	api := AppPixivAPI()
	api.BaseAPI.HookAccessToken(func(token string) { fmt.Println(token) })
	api.BaseAPI.SetAuth("", "")
	res, err := api.UgoiraMetadata(83885225)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.UgoiraMetadata)
	}

}
