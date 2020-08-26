package pixiv

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dghubble/sling"
)

const (
	apiBase = "https://app-api.pixiv.net/"
)

//AppPixiv AppPixivAPI
type AppPixiv struct {
	sling   *sling.Sling
	BaseAPI *BasePixiv
}

//AppPixivAPI AppPixiv
func AppPixivAPI() *AppPixiv {
	s := sling.New().Base(apiBase).Set("User-Agent", "PixivIOSApp/7.6.2 (iOS 12.2; iPhone9,1)").Set("App-Version", "7.6.2").Set("App-OS-VERSION", "12.2").Set("App-OS", "ios")
	baseAPI := BasePixivAPI()
	return &AppPixiv{
		sling:   s,
		BaseAPI: baseAPI,
	}
}

/*
type _Params struct {
	_ _ `url:"_,omitemtpy"`
}

func (api *AppPixiv) _(_) (*_, error) {
	path := ""
	data := &_{}
	erro := &PixivResponseError{}
	params := &_{
		_,
	}
	if err := api.request(path, params, data, erro, true); err != nil {
		return nil, err
	}
	if erro.Error.UserMessage != "" {
		return nil, errors.New(erro.Error.UserMessage)
	}
	return &_, nil
}

*/

type UserImages struct {
	Medium string `json:"medium"`
}
type User struct {
	ID         uint64 `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Comment    string `json:"comment"`
	IsFollowed bool   `json:"is_followed"`

	ProfileImages UserImages `json:"profile_image_urls"`
}
type UserDetail struct {
	User *User `json:"user"`
	// TODO:
	// Profile
	// ProfilePublicity
	// Workspace
}
type Tag struct {
	Name string `json:"name"`
}
type Images struct {
	SquareMedium string `json:"square_medium"`
	Medium       string `json:"medium"`
	Large        string `json:"large"`
	Original     string `json:"original"`
}
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url"`
}
type MetaPage struct {
	Images Images `json:"image_urls"`
}
type Illust struct {
	ID          uint64   `json:"id"`
	Title       string   `json:"title"`
	Type        string   `json:"type"`
	Images      Images   `json:"image_urls"`
	Caption     string   `json:"caption"`
	Restrict    int      `json:"restrict"`
	User        User     `json:"user"`
	Tags        []Tag    `json:"tags"`
	Tools       []string `json:"tools"`
	CreateData  string   `json:"create_data"`
	PageCount   int      `json:"page_count"`
	Width       int      `json:"width"`
	Height      int      `json:"height"`
	SanityLevel int      `json:"sanity_level"`
	// TODO:
	// Series `json:"series"`
	MetaSinglePage MetaSinglePage `json:"meta_single_page"`
	MetaPages      []MetaPage     `json:"meta_pages"`
	TotalView      int            `json:"total_view"`
	TotalBookmarks int            `json:"total_bookmarks"`
	IsBookmarked   bool           `json:"is_bookmarked"`
	Visible        bool           `json:"visible"`
	IsMuted        bool           `json:"is_muted"`
	TotalComments  int            `json:"total_comments"`
}

type IllustsResponse struct {
	Illusts []Illust `json:"illusts"`
	NextURL string   `json:"next_url"`
}
type IllustResponse struct {
	Illust Illust `json:"illust"`
}

//PixivResponseError PixivResponseError
type PixivResponseError struct {
	Error PixivError `json:"error"`
}
type UserMessageDetail struct {
}
type PixivError struct {
	Message            string            `json:"message"`
	Reason             string            `json:"reason"`
	UserMessage        string            `json:"user_message"`
	UserMessageDetails UserMessageDetail `json:"user_message_details"`
}
type ParentComment struct {
}
type IllustComments struct {
	ID            int           `json:"id"`
	Comment       string        `json:"comment"`
	Date          time.Time     `json:"date"`
	User          User          `json:"user"`
	ParentComment ParentComment `json:"parent_comment"`
	//_       _         `json:"_"`
}
type IllustCommentsResponse struct {
	TotalComments int              `json:"total_comments"`
	Comments      []IllustComments `json:"comments"`
	NextURL       string           `json:"next_url"`
}

type TrendingTagsIllust struct {
	TranslatedName string `json:"translated_name"`
	Tag            string `json:"tag"`
	Illust         Illust `json:"illust"`
}

type TrendingTagsIllustResponse struct {
	TrendTags []TrendingTagsIllust `json:"trend_tags"`
}

type UserPreviews struct {
	User    *User    `json:"user"`
	Illusts []Illust `json:"illusts"`
	IsMuted bool     `json:"is_muted"`
	//novels
}
type UserResponse struct {
	UserPreviews []UserPreviews `json:"user_previews"`
	NextURL      string         `json:"next_url"`
}

type UgoiraMetadataFrame struct {
	Delay int    `json:"delay"`
	File  string `json:"file"`
}
type UgoiraMetadataZipUrls struct {
	Medium   string `json:"medium"`
	Large    string `json:"large"`
	Original string `json:"original"`
}

type UgoiraMetadata struct {
	Frames  []UgoiraMetadataFrame `json:"frames"`
	ZipUrls UgoiraMetadataZipUrls `json:"zip_urls"`
}
type UgoiraMetadataResponse struct {
	UgoiraMetadata UgoiraMetadata `json:"ugoira_metadata"`
}

func (api *AppPixiv) request(path string, params, data interface{}, auth bool) (err error) {
requestStart:
	erro := &PixivResponseError{}
	if auth {
		if _, err := api.BaseAPI.Login("", ""); err != nil {
			return fmt.Errorf("refresh token failed: %v", err)
		}
		_, err = api.sling.New().Get(path).Set("Authorization", "Bearer "+api.BaseAPI.AccessToken).QueryStruct(params).Receive(data, erro)
		if strings.Contains(erro.Error.Message, "invalid_grant") {
			api.BaseAPI.TokenDeadline = time.Now()
			if _, err := api.BaseAPI.Login("", ""); err != nil {
				return fmt.Errorf("refresh token failed: %v", err)
			}
			goto requestStart
		}
	} else {
		_, err = api.sling.New().Get(path).QueryStruct(params).Receive(data, erro)
	}

	switch {
	case erro.Error.UserMessage != "":
		return errors.New(erro.Error.UserMessage)
	case erro.Error.Message != "":
		return errors.New(erro.Error.UserMessage)
	}
	return err
}

type illustDetailParams struct {
	IllustID uint64 `url:"illust_id,omitemtpy"`
}

//IllustDetail gets the detail of an illust
func (api *AppPixiv) IllustDetail(id uint64) (*Illust, error) {
	path := "v1/illust/detail"
	data := &IllustResponse{}
	params := &illustDetailParams{
		IllustID: id,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return &data.Illust, nil
}

type userDetailParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
}

//UserDetail gets users information
func (api *AppPixiv) UserDetail(uid uint64) (*UserDetail, error) {
	path := "v1/user/detail"
	params := &userDetailParams{
		UserID: uid,
		Filter: "for_ios",
	}
	detail := &UserDetail{
		User: &User{},
	}
	if err := api.request(path, params, detail, true); err != nil {
		return nil, err
	}

	return detail, nil
}

type userIllustsParams struct {
	UserID uint64 `url:"user_id,omitempty"`
	Filter string `url:"filter,omitempty"`
	Type   string `url:"type,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

// UserIllusts type: [illust, manga]
func (api *AppPixiv) UserIllusts(uid uint64, _type string, offset int) ([]Illust, int, error) {
	path := "v1/user/illusts"
	params := &userIllustsParams{
		UserID: uid,
		Filter: "for_ios",
		Type:   _type,
		Offset: offset,
	}
	data := &IllustsResponse{}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type userBookmarkIllustsParams struct {
	UserID        uint64 `url:"user_id,omitempty"`
	Restrict      string `url:"restrict,omitempty"`
	Filter        string `url:"filter,omitempty"`
	MaxBookmarkID int    `url:"max_bookmark_id,omitempty"`
	Tag           string `url:"tag,omitempty"`
}

// UserBookmarksIllust restrict: [public, private]
func (api *AppPixiv) UserBookmarksIllust(uid uint64, restrict string, maxBookmarkID int, tag string) ([]Illust, int, error) {
	path := "v1/user/bookmarks/illust"
	params := &userBookmarkIllustsParams{
		UserID:        uid,
		Restrict:      "public",
		Filter:        "for_ios",
		MaxBookmarkID: maxBookmarkID,
		Tag:           tag,
	}
	data := &IllustsResponse{}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

type illustFollowParams struct {
	Restrict string `url:"restrict,omitempty"`
	Offset   int    `url:"offset,omitempty"`
}

// IllustFollow restrict: [public, private]
func (api *AppPixiv) IllustFollow(restrict string, offset int) ([]Illust, int, error) {
	path := "v2/illust/follow"
	params := &illustFollowParams{
		Restrict: restrict,
		Offset:   offset,
	}
	data := &IllustsResponse{}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data.Illusts, next, err
}

//IllustCommentsParams is used by function IllustComments
type IllustCommentsParams struct {
	IllustID             uint64 `url:"illust_id,omitemtpy"`
	Offset               int    `url:"offset,omitempty"`
	IncludeTotalComments bool   `url:"include_total_comments,omitempty"`
}

//IllustComments get the comments of an Illust
func (api *AppPixiv) IllustComments(IllustID uint64, offset int, IncludeTotalComments bool) (*IllustCommentsResponse, int, error) {
	path := "v1/illust/comments"
	data := &IllustCommentsResponse{}
	params := &IllustCommentsParams{
		IllustID:             IllustID,
		Offset:               offset,
		IncludeTotalComments: IncludeTotalComments,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

type IllustRelatedParams struct {
	IllustID uint64 `url:"illust_id,omitemtpy"`
	Offset   int    `url:"offset,omitempty"`
	Filter   string `url:"filter,omitempty"`
}

//IllustRelated get the related illusts of an Illust
func (api *AppPixiv) IllustRelated(IllustID uint64, offset int) (*IllustsResponse, int, error) {
	path := "v2/illust/related"
	data := &IllustsResponse{}
	params := &IllustRelatedParams{
		IllustID: IllustID,
		Offset:   offset,
		Filter:   "for_ios",
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

// IllustRecommendedParams IllustRecommended
type IllustRecommendedParams struct {
	ContentType                  string `url:"content_type,omitemtpy"`
	includeRankingLabel          bool   `url:"include_ranking_label,omitempty"`
	Filter                       string `url:"filter,omitempty"`
	MaxBookmarkIDForRecommend    uint64 `url:"max_bookmark_id_for_recommend,omitempty"`
	MinBookmarkIDForRecentIllust uint64 `url:"min_bookmark_id_for_recent_illust,omitempty"`
	Offset                       int    `url:"offset,omitempty"`
	IncludeRankingIllusts        bool   `url:"include_ranking_illusts,omitempty"`
}

//IllustRecommended ("illust", true, 0, 0, 0, true, true)
func (api *AppPixiv) IllustRecommended(contentType string, includeRankingLabel bool, maxBookmarkIDForRecommend uint64, minBookmarkIDForRecentIllust uint64, offset int, includeRankingIllusts bool, reqAuth bool) (*IllustsResponse, int, error) {
	var path string
	if reqAuth {
		path = "v1/illust/recommended"
	} else {
		path = "v1/illust/recommended-nologin"
	}

	data := &IllustsResponse{}
	params := &IllustRecommendedParams{
		ContentType:                  contentType,
		includeRankingLabel:          includeRankingLabel,
		Filter:                       "for_ios",
		MaxBookmarkIDForRecommend:    maxBookmarkIDForRecommend,
		MinBookmarkIDForRecentIllust: minBookmarkIDForRecentIllust,
		Offset:                       offset,
		IncludeRankingIllusts:        includeRankingIllusts,
	}
	if err := api.request(path, params, data, reqAuth); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

// IllustRankingParams IllustRanking
type IllustRankingParams struct {
	Mode   string `url:"mode,omitemtpy"`
	Date   string `url:"date,omitempty"`
	Filter string `url:"filter,omitempty"`
	Offset int    `url:"offset,omitempty"`
}

//IllustRanking ("day", "", 0)
func (api *AppPixiv) IllustRanking(Mode string, Date string, offset int) (*IllustsResponse, int, error) {
	path := "v1/illust/ranking"
	data := &IllustsResponse{}
	params := &IllustRankingParams{
		Filter: "for_ios",
		Mode:   Mode,
		Date:   Date,
		Offset: offset,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

//TrendingTagsIllustParams TrendingTagsIllust
type TrendingTagsIllustParams struct {
	Filter string `url:"filter,omitempty"`
}

//TrendingTagsIllust ("day", "", 0)
func (api *AppPixiv) TrendingTagsIllust() (*TrendingTagsIllustResponse, error) {
	path := "v1/trending-tags/illust"
	data := &TrendingTagsIllustResponse{}
	params := &TrendingTagsIllustParams{
		Filter: "for_ios",
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}

type SearchIllustParams struct {
	Word         string `url:"word,omitempty"`
	SearchTarget string `url:"search_target,omitempty"`
	Sort         string `url:"sort,omitempty"`
	StartDate    string `url:"start_date,omitempty"`
	EndDate      string `url:"end_date,omitempty"`
	Duration     string `url:"duration,omitempty"`
	Offset       int    `url:"offset,omitempty"`
	Filter       string `url:"filter,omitempty"`
}

//SearchIllust ("Sagiri", "partial_match_for_tags", "date_desc", "", "", "", 0)
func (api *AppPixiv) SearchIllust(Word string, searchTarget string, sort string, startDate string, endDate string, duration string, offset int) (*IllustsResponse, int, error) {
	path := "v1/search/illust"
	data := &IllustsResponse{}
	params := &SearchIllustParams{
		Filter:       "for_ios",
		Word:         Word,
		SearchTarget: searchTarget,
		Sort:         sort,
		StartDate:    startDate,
		EndDate:      endDate,
		Duration:     duration,
		Offset:       offset,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

type SearchUserParams struct {
	Word     string `url:"word,omitempty"`
	Sort     string `url:"sort,omitempty"`
	Duration string `url:"duration,omitempty"`
	Offset   int    `url:"offset,omitempty"`
	Filter   string `url:"filter,omitempty"`
}

//SearchUser ("Quan_", "date_desc", "", 0)
func (api *AppPixiv) SearchUser(Word string, sort string, duration string, offset int) (*UserResponse, int, error) {
	path := "v1/search/user"
	data := &UserResponse{}
	params := &SearchUserParams{
		Filter:   "for_ios",
		Word:     Word,
		Sort:     sort,
		Duration: duration,
		Offset:   offset,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, 0, err
	}
	next, err := parseNextPageOffset(data.NextURL)
	return data, next, err
}

type UgoiraMetadataParams struct {
	IllustID uint64 `url:"illust_id,omitempty"`
}

//UgoiraMetadata ("day", "", 0)
func (api *AppPixiv) UgoiraMetadata(illustID uint64) (*UgoiraMetadataResponse, error) {
	path := "v1/ugoira/metadata"
	data := &UgoiraMetadataResponse{}
	params := &UgoiraMetadataParams{
		IllustID: illustID,
	}
	if err := api.request(path, params, data, true); err != nil {
		return nil, err
	}
	return data, nil
}
