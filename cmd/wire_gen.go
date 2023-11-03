// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package cmd

import (
	"github.com/google/wire"
	"github.com/navidrome/navidrome/core"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/core/agents/lastfm"
	"github.com/navidrome/navidrome/core/agents/listenbrainz"
	"github.com/navidrome/navidrome/core/artwork"
	"github.com/navidrome/navidrome/core/external_playlists"
	"github.com/navidrome/navidrome/core/ffmpeg"
	"github.com/navidrome/navidrome/core/scrobbler"
	"github.com/navidrome/navidrome/db"
	"github.com/navidrome/navidrome/persistence"
	"github.com/navidrome/navidrome/scanner"
	"github.com/navidrome/navidrome/server"
	"github.com/navidrome/navidrome/server/events"
	"github.com/navidrome/navidrome/server/nativeapi"
	"github.com/navidrome/navidrome/server/public"
	"github.com/navidrome/navidrome/server/subsonic"
	"sync"
)

// Injectors from wire_injectors.go:

func CreateServer(musicFolder string) *server.Server {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	broker := events.GetBroker()
	serverServer := server.New(dataStore, broker)
	return serverServer
}

func CreateNativeAPIRouter() *nativeapi.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	share := core.NewShare(dataStore)
	playlistRetriever := external_playlists.GetPlaylistRetriever(dataStore)
	playlists := core.NewPlaylists(dataStore)
	router := nativeapi.New(dataStore, share, playlistRetriever, playlists)
	return router
}

func CreateSubsonicAPIRouter() *subsonic.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	fileCache := artwork.GetImageCache()
	fFmpeg := ffmpeg.New()
	agentsAgents := agents.New(dataStore)
	externalMetadata := core.NewExternalMetadata(dataStore, agentsAgents)
	artworkArtwork := artwork.NewArtwork(dataStore, fileCache, fFmpeg, externalMetadata)
	transcodingCache := core.GetTranscodingCache()
	mediaStreamer := core.NewMediaStreamer(dataStore, fFmpeg, transcodingCache)
	share := core.NewShare(dataStore)
	archiver := core.NewArchiver(mediaStreamer, dataStore, share)
	players := core.NewPlayers(dataStore)
	scanner := GetScanner()
	broker := events.GetBroker()
	playlists := core.NewPlaylists(dataStore)
	playTracker := scrobbler.GetPlayTracker(dataStore, broker)
	router := subsonic.New(dataStore, artworkArtwork, mediaStreamer, archiver, players, externalMetadata, scanner, broker, playlists, playTracker, share)
	return router
}

func CreatePublicRouter() *public.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	fileCache := artwork.GetImageCache()
	fFmpeg := ffmpeg.New()
	agentsAgents := agents.New(dataStore)
	externalMetadata := core.NewExternalMetadata(dataStore, agentsAgents)
	artworkArtwork := artwork.NewArtwork(dataStore, fileCache, fFmpeg, externalMetadata)
	transcodingCache := core.GetTranscodingCache()
	mediaStreamer := core.NewMediaStreamer(dataStore, fFmpeg, transcodingCache)
	share := core.NewShare(dataStore)
	archiver := core.NewArchiver(mediaStreamer, dataStore, share)
	router := public.New(dataStore, artworkArtwork, mediaStreamer, share, archiver)
	return router
}

func CreateLastFMRouter() *lastfm.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	router := lastfm.NewRouter(dataStore)
	return router
}

func CreateListenBrainzRouter() *listenbrainz.Router {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	router := listenbrainz.NewRouter(dataStore)
	return router
}

func createScanner() scanner.Scanner {
	sqlDB := db.Db()
	dataStore := persistence.New(sqlDB)
	playlists := core.NewPlaylists(dataStore)
	fileCache := artwork.GetImageCache()
	fFmpeg := ffmpeg.New()
	agentsAgents := agents.New(dataStore)
	externalMetadata := core.NewExternalMetadata(dataStore, agentsAgents)
	artworkArtwork := artwork.NewArtwork(dataStore, fileCache, fFmpeg, externalMetadata)
	cacheWarmer := artwork.NewCacheWarmer(artworkArtwork, fileCache)
	broker := events.GetBroker()
	playlistRetriever := external_playlists.GetPlaylistRetriever(dataStore)
	scannerScanner := scanner.New(dataStore, playlists, cacheWarmer, broker, playlistRetriever)
	return scannerScanner
}

// wire_injectors.go:

var allProviders = wire.NewSet(core.Set, artwork.Set, subsonic.New, nativeapi.New, public.New, persistence.New, lastfm.NewRouter, listenbrainz.NewRouter, events.GetBroker, db.Db)

// Scanner must be a Singleton
var (
	onceScanner     sync.Once
	scannerInstance scanner.Scanner
)

func GetScanner() scanner.Scanner {
	onceScanner.Do(func() {
		scannerInstance = createScanner()
	})
	return scannerInstance
}
