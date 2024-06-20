package listenbrainz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/navidrome/navidrome/core/external_playlists"
	"github.com/navidrome/navidrome/log"
	"github.com/navidrome/navidrome/model"
)

type listenBrainzError struct {
	Code    int
	Message string
}

func (e *listenBrainzError) Error() string {
	return fmt.Sprintf("ListenBrainz error(%d): %s", e.Code, e.Message)
}

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func newClient(baseURL string, hc httpDoer) *client {
	return &client{baseURL, hc}
}

type client struct {
	baseURL string
	hc      httpDoer
}

type listenBrainzResponse struct {
	Code          int               `json:"code"`
	Message       string            `json:"message"`
	Error         string            `json:"error"`
	Status        string            `json:"status"`
	Valid         bool              `json:"valid"`
	UserName      string            `json:"user_name"`
	PlaylistCount int               `json:"playlist_count"`
	Playlists     []overallPlaylist `json:"playlists,omitempty"`
	Playlist      lbPlaylist        `json:"playlist"`
}

type listenBrainzRequest struct {
	ApiKey string
	Body   *listenBrainzRequestBody
}

type overallPlaylist struct {
	Playlist lbPlaylist `json:"playlist"`
}

type lbPlaylist struct {
	Annotation string       `json:"annotation"`
	Creator    string       `json:"creator"`
	Date       time.Time    `json:"date"`
	Identifier string       `json:"identifier"`
	Title      string       `json:"title"`
	Extension  plsExtension `json:"extension"`
	Tracks     []lbTrack    `json:"track"`
}

type plsExtension struct {
	Extension playlistExtension `json:"https://musicbrainz.org/doc/jspf#playlist"`
}

type playlistExtension struct {
	AdditionalMetadata additionalMeta `json:"additional_metadata"`
	Collaborators      []string       `json:"collaborators"`
	CreatedFor         string         `json:"created_for"`
	LastModified       time.Time      `json:"last_modified_at"`
	Public             bool           `json:"public"`
}

type additionalMeta struct {
	AlgorithmMetadata algoMeta `json:"algorithm_metadata"`
}

type algoMeta struct {
	SourcePatch string `json:"source_patch"`
}

type lbTrack struct {
	Creator    string   `json:"creator"`
	Identifier []string `json:"identifier"`
	Title      string   `json:"title"`
}

type listenBrainzRequestBody struct {
	ListenType listenType   `json:"listen_type,omitempty"`
	Payload    []listenInfo `json:"payload,omitempty"`
}

type listenType string

const (
	Single     listenType = "single"
	PlayingNow listenType = "playing_now"
)

type listenInfo struct {
	ListenedAt    int           `json:"listened_at,omitempty"`
	TrackMetadata trackMetadata `json:"track_metadata,omitempty"`
}

type trackMetadata struct {
	ArtistName     string         `json:"artist_name,omitempty"`
	TrackName      string         `json:"track_name,omitempty"`
	ReleaseName    string         `json:"release_name,omitempty"`
	AdditionalInfo additionalInfo `json:"additional_info,omitempty"`
}

type additionalInfo struct {
	SubmissionClient        string   `json:"submission_client,omitempty"`
	SubmissionClientVersion string   `json:"submission_client_version,omitempty"`
	TrackNumber             int      `json:"tracknumber,omitempty"`
	RecordingMbzID          string   `json:"recording_mbid,omitempty"`
	ArtistMbzIDs            []string `json:"artist_mbids,omitempty"`
	ReleaseMbID             string   `json:"release_mbid,omitempty"`
	DurationMs              int      `json:"duration_ms,omitempty"`
}

type trackInfo struct {
	RecordingName string `json:"recording_name"`
	RecordingMbid string `json:"recording_mbid"`
}

func (c *client) validateToken(ctx context.Context, apiKey string) (*listenBrainzResponse, error) {
	r := &listenBrainzRequest{
		ApiKey: apiKey,
	}
	response, err := c.makeRequest(ctx, http.MethodGet, "validate-token", "", r)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (c *client) updateNowPlaying(ctx context.Context, apiKey string, li listenInfo) error {
	r := &listenBrainzRequest{
		ApiKey: apiKey,
		Body: &listenBrainzRequestBody{
			ListenType: PlayingNow,
			Payload:    []listenInfo{li},
		},
	}

	resp, err := c.makeRequest(ctx, http.MethodPost, "submit-listens", "", r)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		log.Warn(ctx, "ListenBrainz: NowPlaying was not accepted", "status", resp.Status)
	}
	return nil
}

func (c *client) scrobble(ctx context.Context, apiKey string, li listenInfo) error {
	r := &listenBrainzRequest{
		ApiKey: apiKey,
		Body: &listenBrainzRequestBody{
			ListenType: Single,
			Payload:    []listenInfo{li},
		},
	}
	resp, err := c.makeRequest(ctx, http.MethodPost, "submit-listens", "", r)
	if err != nil {
		return err
	}
	if resp.Status != "ok" {
		log.Warn(ctx, "ListenBrainz: Scrobble was not accepted", "status", resp.Status)
	}
	return nil
}

func (c *client) getPlaylists(ctx context.Context, offset, count int, apiKey, user, plsType string) (*listenBrainzResponse, error) {
	r := &listenBrainzRequest{
		ApiKey: apiKey,
		Body:   nil,
	}

	var endpoint string

	switch plsType {
	case "user":
		endpoint = "user/" + user + "/playlists"
	case "created":
		endpoint = "user/" + user + "/playlists/createdfor"
	case "collab":
		endpoint = "user/" + user + "/playlists/collaborator"
	default:
		return nil, external_playlists.ErrorUnsupportedType
	}

	extra := fmt.Sprintf("?count=%d&offset=%d", count, offset)

	resp, err := c.makeRequest(ctx, http.MethodGet, endpoint, extra, r)

	if err != nil {
		return nil, err
	}

	return resp, err
}

func (c *client) getPlaylist(ctx context.Context, apiKey, plsId string) (*listenBrainzResponse, error) {
	r := &listenBrainzRequest{
		ApiKey: apiKey,
	}

	endpoint := fmt.Sprintf("playlist/%s", plsId)

	resp, err := c.makeRequest(ctx, http.MethodGet, endpoint, "", r)

	if resp != nil && resp.Code == 404 {
		return nil, model.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *client) path(endpoint string) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}
	u.Path = path.Join(u.Path, endpoint)
	return u.String(), nil
}

// https://listenbrainz.readthedocs.io/en/latest/users/api/popularity.html#get--1-popularity-top-recordings-for-artist-(artist_mbid)
// Note that this is popularity by listen. There is (as of June 15, 2024) no way
// to limit the output
func (c *client) getTopSongs(ctx context.Context, mbid string) ([]trackInfo, error) {
	r := &listenBrainzRequest{}
	endpoint := fmt.Sprintf("popularity/top-recordings-for-artist/%s", mbid)

	response, err := c.makeLbzRequest(ctx, http.MethodGet, endpoint, "", r)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	if response.StatusCode != 200 {
		var response listenBrainzResponse
		jsonErr := decoder.Decode(&response)

		if jsonErr != nil {
			return nil, jsonErr
		}
		if response.Code != 0 && response.Code != 200 {
			return nil, &listenBrainzError{Code: response.Code, Message: response.Error}
		}
	}

	var tracks []trackInfo
	jsonErr := decoder.Decode(&tracks)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return tracks, nil
}

const (
	labsBase  = "https://labs.api.listenbrainz.org/"
	algorithm = "session_based_days_9000_session_300_contribution_5_threshold_15_limit_50_skip_30"
)

type artist struct {
	MBID string `json:"artist_mbid"`
	Name string `json:"name"`
}

func (c *client) getSimilarArtists(ctx context.Context, mbid string) ([]artist, error) {
	r := &listenBrainzRequest{}
	url := fmt.Sprintf("%ssimilar-artists/json?artist_mbids=%s&algorithm=%s", labsBase, mbid, algorithm)

	response, err := c.makeRawRequest(ctx, http.MethodGet, url, "", r)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	decoder := json.NewDecoder(response.Body)

	var artists []artist
	jsonErr := decoder.Decode(&artists)

	if jsonErr != nil {
		return nil, jsonErr
	}

	return artists, nil
}

func (c *client) makeRawRequest(ctx context.Context, method string, uri string, query string, r *listenBrainzRequest) (*http.Response, error) {
	if query != "" {
		uri += query
	}

	var req *http.Request

	if r.Body != nil {
		b, _ := json.Marshal(r.Body)
		req, _ = http.NewRequestWithContext(ctx, method, uri, bytes.NewBuffer(b))
	} else {
		req, _ = http.NewRequestWithContext(ctx, method, uri, nil)
	}

	req.Header.Add("Content-Type", "application/json; charset=UTF-8")

	if r.ApiKey != "" {
		req.Header.Add("Authorization", fmt.Sprintf("Token %s", r.ApiKey))
	}

	log.Trace(ctx, fmt.Sprintf("Sending ListenBrainz %s request", req.Method), "url", req.URL)
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *client) makeLbzRequest(ctx context.Context, method string, endpoint string, query string, r *listenBrainzRequest) (*http.Response, error) {
	uri, err := c.path(endpoint)
	if err != nil {
		return nil, err
	}

	return c.makeRawRequest(ctx, method, uri, query, r)
}

func (c *client) makeRequest(ctx context.Context, method string, endpoint string, query string, r *listenBrainzRequest) (*listenBrainzResponse, error) {
	resp, err := c.makeLbzRequest(ctx, method, endpoint, query, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	var response listenBrainzResponse
	jsonErr := decoder.Decode(&response)
	if resp.StatusCode != 200 && jsonErr != nil {
		return nil, fmt.Errorf("ListenBrainz: HTTP Error, Status: (%d)", resp.StatusCode)
	}
	if jsonErr != nil {
		return nil, jsonErr
	}
	if response.Code != 0 && response.Code != 200 {
		return &response, &listenBrainzError{Code: response.Code, Message: response.Error}
	}

	return &response, nil
}
