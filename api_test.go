package pixiv

import (
	"fmt"
	"testing"
)

func TestHello(t *testing.T) {
	api := BasePixivAPI()
	api.SetAuth("", "")
	fmt.Println(api)
}
