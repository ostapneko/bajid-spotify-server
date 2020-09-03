export const getSongList = async (userId) => {
    const resp = await fetch(`/song_list/${userId}`)
    return await resp.json()
}
