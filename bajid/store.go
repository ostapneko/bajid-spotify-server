package bajid

type SongListStore interface {
	Read(userId UserId) (SongList, error)
	Write(userId UserId, songList SongList) error
}
