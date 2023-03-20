package core

import (
	"github.com/google/wire"
	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/core/external_playlists"
	"github.com/navidrome/navidrome/core/ffmpeg"
	"github.com/navidrome/navidrome/core/scrobbler"
)

var Set = wire.NewSet(
	NewMediaStreamer,
	GetTranscodingCache,
	NewArchiver,
	NewExternalMetadata,
	NewPlayers,
	NewShare,
	NewPlaylists,
	agents.New,
	ffmpeg.New,
	scrobbler.GetPlayTracker,
	external_playlists.GetPlaylistRetriever,
)
