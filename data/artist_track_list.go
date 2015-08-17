package data

type ArtistTrackList []ArtistTrack
type ArtistTrackList2 *[]ArtistTrack

func (p ArtistTrackList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p ArtistTrackList) Len() int           { return len(p) }
func (p ArtistTrackList) Less(i, j int) bool { return len(p[i].Tracks) < len(p[j].Tracks) }
