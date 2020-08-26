# pixiv

Pixiv API for Golang (with Auth supported)

Forked from [everpcpc/pixiv](https://github.com/everpcpc/pixiv)
Inspired by [pixivpy](https://github.com/upbit/pixivpy)

I applied some changes to meet my purpose.

## Useage

```golang
api := pixiv.AppPixivAPI()
api.BaseAPI.HookAccessToken(func(token string) { fmt.Println(token) })
api.BaseAPI.Login("#Username","@Password")
user, err := app.UserDetail(uid)
illusts, next, err := app.UserIllusts(uid, "illust", 0)
illusts, next, err := app.UserBookmarksIllust(uid, "public", 0, "")
illusts, next, err := app.IllustFollow("public", 0)
```
