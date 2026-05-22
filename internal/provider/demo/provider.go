// Package demo provides an in-memory Provider that returns realistic fake
// data. No Apple credentials or network access are required.
package demo

import (
	"bytes"
	"context"
	"encoding/base64"
	"image"
	"image/color"
	"image/png"
	"strings"
	"time"

	"github.com/simone-vibes/vibez/internal/provider"
)

// Tracks is the built-in demo library shared with the demo Player.
var Tracks = []provider.Track{
	{ID: "d1", Title: "Nights", Artist: "Frank Ocean", Album: "Blonde", Duration: dur(5, 7), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 228, G: 188, B: 89, A: 255}, color.NRGBA{R: 58, G: 34, B: 90, A: 255}, color.NRGBA{R: 245, G: 102, B: 91, A: 255})},
	{ID: "d2", Title: "Pyramids", Artist: "Frank Ocean", Album: "channel ORANGE", Duration: dur(9, 2), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 232, G: 116, B: 44, A: 255}, color.NRGBA{R: 65, G: 31, B: 20, A: 255}, color.NRGBA{R: 251, G: 210, B: 94, A: 255})},
	{ID: "d3", Title: "Novacane", Artist: "Frank Ocean", Album: "nostalgia, ULTRA", Duration: dur(5, 7), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 78, G: 150, B: 187, A: 255}, color.NRGBA{R: 20, G: 30, B: 42, A: 255}, color.NRGBA{R: 238, G: 239, B: 216, A: 255})},
	{ID: "d4", Title: "Redbone", Artist: "Childish Gambino", Album: "Awaken, My Love!", Duration: dur(5, 27), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 172, G: 56, B: 44, A: 255}, color.NRGBA{R: 30, G: 12, B: 28, A: 255}, color.NRGBA{R: 245, G: 178, B: 101, A: 255})},
	{ID: "d5", Title: "Me and Your Mama", Artist: "Childish Gambino", Album: "Awaken, My Love!", Duration: dur(4, 40), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 106, G: 50, B: 129, A: 255}, color.NRGBA{R: 20, G: 14, B: 38, A: 255}, color.NRGBA{R: 242, G: 129, B: 74, A: 255})},
	{ID: "d6", Title: "See You Again", Artist: "Tyler, The Creator", Album: "Flower Boy", Duration: dur(3, 1), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 246, G: 142, B: 177, A: 255}, color.NRGBA{R: 74, G: 41, B: 94, A: 255}, color.NRGBA{R: 142, G: 213, B: 155, A: 255})},
	{ID: "d7", Title: "Garden Shed", Artist: "Tyler, The Creator", Album: "Flower Boy", Duration: dur(3, 32), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 119, G: 190, B: 111, A: 255}, color.NRGBA{R: 39, G: 71, B: 47, A: 255}, color.NRGBA{R: 248, G: 197, B: 97, A: 255})},
	{ID: "d8", Title: "Kill Bill", Artist: "SZA", Album: "SOS", Duration: dur(2, 33), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 68, G: 130, B: 184, A: 255}, color.NRGBA{R: 12, G: 36, B: 62, A: 255}, color.NRGBA{R: 241, G: 225, B: 153, A: 255})},
	{ID: "d9", Title: "Good Days", Artist: "SZA", Album: "Good Days", Duration: dur(4, 39), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 80, G: 172, B: 148, A: 255}, color.NRGBA{R: 20, G: 55, B: 57, A: 255}, color.NRGBA{R: 242, G: 214, B: 135, A: 255})},
	{ID: "d10", Title: "After The Storm", Artist: "Kali Uchis", Album: "Isolation", Duration: dur(3, 57), ArtworkURL: generatedArtworkURL(color.NRGBA{R: 242, G: 157, B: 94, A: 255}, color.NRGBA{R: 54, G: 45, B: 88, A: 255}, color.NRGBA{R: 116, G: 198, B: 184, A: 255})},
}

func generatedArtworkURL(primary, secondary, accent color.NRGBA) string {
	img := image.NewNRGBA(image.Rect(0, 0, 96, 96))
	for y := range 96 {
		for x := range 96 {
			mix := (x*3 + y*2) % 96
			c := blend(primary, secondary, mix)
			if (x/12+y/12)%3 == 0 {
				c = blend(c, accent, 80)
			}
			if x > 18 && x < 78 && y > 18 && y < 78 && (x+y)%17 < 8 {
				c = blend(c, accent, 140)
			}
			img.SetNRGBA(x, y, c)
		}
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		panic(err)
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func blend(a, b color.NRGBA, amount int) color.NRGBA {
	if amount < 0 || amount > 255 {
		panic("invalid blend amount")
	}
	keep := 255 - amount
	return color.NRGBA{
		R: colorByte((int(a.R)*keep + int(b.R)*amount) / 255),
		G: colorByte((int(a.G)*keep + int(b.G)*amount) / 255),
		B: colorByte((int(a.B)*keep + int(b.B)*amount) / 255),
		A: 255,
	}
}

func colorByte(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func dur(m, s int) time.Duration {
	return time.Duration(m)*time.Minute + time.Duration(s)*time.Second
}

// Provider is the demo Provider implementation.
type Provider struct{}

func (Provider) Name() string          { return "Demo" }
func (Provider) IsAuthenticated() bool { return true }

var demoPlaylists = []provider.Playlist{
	{ID: "dp1", Name: "Late Night Coding", TrackCount: 4},
	{ID: "dp2", Name: "Chill Session", TrackCount: 3},
	{ID: "dp3", Name: "Energy Boost", TrackCount: 3},
}

var playlistTracks = map[string][]provider.Track{
	"dp1": {Tracks[0], Tracks[5], Tracks[6], Tracks[7]},
	"dp2": {Tracks[3], Tracks[8], Tracks[9]},
	"dp3": {Tracks[1], Tracks[2], Tracks[4]},
}

func (Provider) Search(_ context.Context, query string) (*provider.SearchResult, error) {
	q := strings.ToLower(query)
	var tracks []provider.Track
	var albums []provider.Album
	seen := map[string]bool{}

	for _, t := range Tracks {
		if strings.Contains(strings.ToLower(t.Title), q) ||
			strings.Contains(strings.ToLower(t.Artist), q) ||
			strings.Contains(strings.ToLower(t.Album), q) {
			tracks = append(tracks, t)
			if !seen[t.Album] {
				seen[t.Album] = true
				albums = append(albums, provider.Album{
					ID:     "da-" + t.ID,
					Title:  t.Album,
					Artist: t.Artist,
				})
			}
		}
	}
	return &provider.SearchResult{Tracks: tracks, Albums: albums}, nil
}

func (Provider) GetLibraryTracks(_ context.Context) ([]provider.Track, error) {
	out := make([]provider.Track, len(Tracks))
	copy(out, Tracks)
	return out, nil
}

func (Provider) GetLibraryPlaylists(_ context.Context) ([]provider.Playlist, error) {
	out := make([]provider.Playlist, len(demoPlaylists))
	copy(out, demoPlaylists)
	return out, nil
}

func (Provider) GetPlaylistTracks(_ context.Context, id string) ([]provider.Track, error) {
	if tracks, ok := playlistTracks[id]; ok {
		out := make([]provider.Track, len(tracks))
		copy(out, tracks)
		return out, nil
	}
	return nil, nil
}

func (Provider) GetAlbumTracks(_ context.Context, _ string) ([]provider.Track, error) {
	out := make([]provider.Track, len(Tracks))
	copy(out, Tracks)
	return out, nil
}

func (Provider) GetLibraryAlbumTracks(_ context.Context, _ string) ([]provider.Track, error) {
	out := make([]provider.Track, len(Tracks))
	copy(out, Tracks)
	return out, nil
}

func (Provider) GetCatalogPlaylistTracks(_ context.Context, _ string) ([]provider.Track, error) {
	out := make([]provider.Track, len(Tracks))
	copy(out, Tracks)
	return out, nil
}

func (Provider) CreatePlaylist(_ context.Context, name string, _ []string) (provider.Playlist, error) {
	return provider.Playlist{ID: "dp-new-" + name, Name: name}, nil
}

func (Provider) LoveSong(_ context.Context, _ string, _ bool) error { return nil }

func (Provider) GetSongRating(_ context.Context, _ string) (bool, error) { return false, nil }

func (Provider) AddToPlaylist(_ context.Context, _, _ string) error { return nil }

func (Provider) GetRecommendations(_ context.Context) ([]provider.RecommendationGroup, error) {
	return []provider.RecommendationGroup{
		{
			Title: "Recommended for You",
			Items: []provider.RecommendationItem{
				{ID: "demo-album-1", Kind: "album", Title: "A Colours Trilogy", Subtitle: "Jon Hopkins"},
				{ID: "demo-album-2", Kind: "album", Title: "Immunity", Subtitle: "Jon Hopkins"},
			},
		},
		{
			Title: "New Releases",
			Items: []provider.RecommendationItem{
				{ID: "demo-pl-1", Kind: "playlist", Title: "New Music Mix", Subtitle: "Apple Music"},
			},
		},
	}, nil
}
