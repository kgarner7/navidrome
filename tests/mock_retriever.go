package tests

import (
	"context"

	"github.com/navidrome/navidrome/core/external_playlists"
)

type noopRetriever struct {
}

func MockRetriever() external_playlists.PlaylistRetriever {
	r := noopRetriever{}
	return &r
}

// GetAvailableAgents implements external_playlists.PlaylistRetriever.
func (n *noopRetriever) GetAvailableAgents(ctx context.Context, userId string) []external_playlists.PlaylistSourceInfo {
	return nil
}

// GetPlaylists implements external_playlists.PlaylistRetriever.
func (n *noopRetriever) GetPlaylists(ctx context.Context, offset int, count int, userId string, agent string, playlistType string) (*external_playlists.ExternalPlaylists, error) {
	return nil, nil
}

// ImportPlaylists implements external_playlists.PlaylistRetriever.
func (n *noopRetriever) ImportPlaylists(ctx context.Context, update bool, userId string, agent string, mapping external_playlists.ImportMap) error {
	return nil
}

// SyncPlaylist implements external_playlists.PlaylistRetriever.
func (n *noopRetriever) SyncPlaylist(ctx context.Context, playlistId string) error {
	return nil
}

// SyncRecommended implements external_playlists.PlaylistRetriever.
func (n *noopRetriever) SyncRecommended(ctx context.Context, userrId string, agent string) error {
	return nil
}

var _ external_playlists.PlaylistRetriever = (*noopRetriever)(nil)
