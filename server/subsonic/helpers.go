package subsonic

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/model/request"
	"github.com/navidrome/navidrome/server/subsonic/responses"
	"github.com/navidrome/navidrome/utils"
)

func newResponse() *responses.Subsonic {
	return &responses.Subsonic{Status: "ok", Version: Version, Type: consts.AppName, ServerVersion: consts.Version}
}

func requiredParamString(r *http.Request, param string) (string, error) {
	p := utils.ParamString(r, param)
	if p == "" {
		return "", newError(responses.ErrorMissingParameter, "required '%s' parameter is missing", param)
	}
	return p, nil
}

func requiredParamStrings(r *http.Request, param string) ([]string, error) {
	ps := utils.ParamStrings(r, param)
	if len(ps) == 0 {
		return nil, newError(responses.ErrorMissingParameter, "required '%s' parameter is missing", param)
	}
	return ps, nil
}

func requiredParamInt(r *http.Request, param string) (int, error) {
	p := utils.ParamString(r, param)
	if p == "" {
		return 0, newError(responses.ErrorMissingParameter, "required '%s' parameter is missing", param)
	}
	return utils.ParamInt(r, param, 0), nil
}

type subError struct {
	code     int
	messages []interface{}
}

func newError(code int, message ...interface{}) error {
	return subError{
		code:     code,
		messages: message,
	}
}

func (e subError) Error() string {
	var msg string
	if len(e.messages) == 0 {
		msg = responses.ErrorMsg(e.code)
	} else {
		msg = fmt.Sprintf(e.messages[0].(string), e.messages[1:]...)
	}
	return msg
}

func getUser(ctx context.Context) model.User {
	user, ok := request.UserFrom(ctx)
	if ok {
		return user
	}
	return model.User{}
}

func toArtists(ctx context.Context, artists model.Artists, baseUrl string) []responses.Artist {
	as := make([]responses.Artist, len(artists))
	for i, artist := range artists {
		as[i] = toArtist(ctx, artist, baseUrl)
	}
	return as
}

func toArtist(ctx context.Context, a model.Artist, baseUrl string) responses.Artist {
	artist := responses.Artist{
		Id:             a.ID,
		Name:           a.Name,
		AlbumCount:     a.AlbumCount,
		UserRating:     a.Rating,
		ArtistImageUrl: a.ArtistImageUrl(baseUrl),
	}
	if a.Starred {
		artist.Starred = &a.StarredAt
	}
	return artist
}

func toArtistID3(ctx context.Context, a model.Artist, baseUrl string) responses.ArtistID3 {
	artist := responses.ArtistID3{
		Id:             a.ID,
		Name:           a.Name,
		AlbumCount:     a.AlbumCount,
		ArtistImageUrl: a.ArtistImageUrl(baseUrl),
		UserRating:     a.Rating,
	}
	if a.Starred {
		artist.Starred = &a.StarredAt
	}
	return artist
}

func getBaseArtistImageUrl(r *http.Request) string {
	if conf.Server.BaseDomain == "" {
		return ""
	}

	params := url.Values{}
	params.Add("_", time.Now().UTC().Format(time.RFC3339Nano))
	params.Add("c", utils.ParamString(r, "c"))
	params.Add("f", "json")
	params.Add("s", utils.ParamString(r, "s"))
	params.Add("t", utils.ParamString(r, "t"))
	params.Add("u", utils.ParamString(r, "u"))
	params.Add("v", utils.ParamString(r, "v"))
	return fmt.Sprintf("%s%s/rest/getCoverArt?%s",
		conf.Server.BaseDomain,
		conf.Server.BaseURL,
		params.Encode(),
	)
}

func toGenres(genres model.Genres) *responses.Genres {
	response := make([]responses.Genre, len(genres))
	for i, g := range genres {
		response[i] = responses.Genre{
			Name:       g.Name,
			SongCount:  g.SongCount,
			AlbumCount: g.AlbumCount,
		}
	}
	return &responses.Genres{Genre: response}
}

func getTranscoding(ctx context.Context) (format string, bitRate int) {
	if trc, ok := request.TranscodingFrom(ctx); ok {
		format = trc.TargetFormat
	}
	if plr, ok := request.PlayerFrom(ctx); ok {
		bitRate = plr.MaxBitRate
	}
	return
}

// This seems to be duplicated, but it is an initial step into merging `engine` and the `subsonic` packages,
// In the future there won't be any conversion to/from `engine. Entry` anymore
func childFromMediaFile(ctx context.Context, mf model.MediaFile) responses.Child {
	child := responses.Child{}
	child.Id = mf.ID
	child.Title = mf.Title
	child.IsDir = false
	child.Parent = mf.AlbumID
	child.Album = mf.Album
	child.Year = mf.Year
	child.Artist = mf.Artist
	child.Genre = mf.Genre
	child.Track = mf.TrackNumber
	child.Duration = int(mf.Duration)
	child.Size = mf.Size
	child.Suffix = mf.Suffix
	child.BitRate = mf.BitRate
	if mf.HasCoverArt {
		child.CoverArt = mf.ID
	} else {
		child.CoverArt = "al-" + mf.AlbumID
	}
	child.ContentType = mf.ContentType()
	player, ok := request.PlayerFrom(ctx)
	if ok && player.ReportRealPath {
		child.Path = mf.Path
	} else {
		child.Path = fmt.Sprintf("%s/%s/%s.%s", mapSlashToDash(mf.AlbumArtist), mapSlashToDash(mf.Album), mapSlashToDash(mf.Title), mf.Suffix)
	}
	child.DiscNumber = mf.DiscNumber
	child.Created = &mf.CreatedAt
	child.AlbumId = mf.AlbumID
	child.ArtistId = mf.ArtistID
	child.Type = "music"
	child.PlayCount = mf.PlayCount
	if mf.Starred {
		child.Starred = &mf.StarredAt
	}
	child.UserRating = mf.Rating

	format, _ := getTranscoding(ctx)
	if mf.Suffix != "" && format != "" && mf.Suffix != format {
		child.TranscodedSuffix = format
		child.TranscodedContentType = mime.TypeByExtension("." + format)
	}
	child.BookmarkPosition = mf.BookmarkPosition
	return child
}

func mapSlashToDash(target string) string {
	return strings.ReplaceAll(target, "/", "_")
}

func childrenFromMediaFiles(ctx context.Context, mfs model.MediaFiles) []responses.Child {
	children := make([]responses.Child, len(mfs))
	for i, mf := range mfs {
		children[i] = childFromMediaFile(ctx, mf)
	}
	return children
}

func childFromAlbum(ctx context.Context, al model.Album) responses.Child {
	child := responses.Child{}
	child.Id = al.ID
	child.IsDir = true
	child.Title = al.Name
	child.Name = al.Name
	child.Album = al.Name
	child.Artist = al.AlbumArtist
	child.Year = al.MaxYear
	child.Genre = al.Genre
	child.CoverArt = al.CoverArtId
	child.Created = &al.CreatedAt
	child.Parent = al.AlbumArtistID
	child.ArtistId = al.AlbumArtistID
	child.Duration = int(al.Duration)
	child.SongCount = al.SongCount
	if al.Starred {
		child.Starred = &al.StarredAt
	}
	child.PlayCount = al.PlayCount
	child.UserRating = al.Rating
	return child
}

func childrenFromAlbums(ctx context.Context, als model.Albums) []responses.Child {
	children := make([]responses.Child, len(als))
	for i, al := range als {
		children[i] = childFromAlbum(ctx, al)
	}
	return children
}
