const getSpotifyToken = () => {
    return document.cookie
        .split(';')
        .find(r => r.startsWith('bajid-spotify-token'))
        .split('=')[1];
};

const token = getSpotifyToken();
const auth_headers = {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
};

export const initSpotify = () => {
    return new Promise(((resolve, reject) => {
        window.onSpotifyWebPlaybackSDKReady = () => resolve()
    }))
}

export const getUserID = async () => {
    const resp = await fetch('https://api.spotify.com/v1/me', {headers: auth_headers});
    const json = await resp.json()
    return json.id
}

export const getTracks = async (uris) => {
    const ids = uris.map(uri => uri.split(':')[2])
    const idsStr = ids.join(',')
    const resp = await fetch(`https://api.spotify.com/v1/tracks?ids=${idsStr}`, {headers: auth_headers})
    const doc = await resp.json();
    return doc.tracks;
}

export const setRepeatTrack = async () => {
    await fetch('https://api.spotify.com/v1/me/player/repeat?state=track', {method: 'PUT', headers: auth_headers});
}


export const createPlayer = async () => {
    return new Promise((resolve, reject) => {
        const player = new Spotify.Player({
            name: 'Bajid Box Player',

            getOAuthToken: cb => {
                cb(token);
            }
        });

        // Error handling
        player.addListener('initialization_error', ({message}) => {
            reject(message);
        });
        player.addListener('authentication_error', ({message}) => {
            reject(message);
        });
        player.addListener('account_error', ({message}) => {
            reject(message);
        });
        player.addListener('playback_error', ({message}) => {
            console.log(message);
        });

        // Ready
        player.addListener('ready', ({device_id}) => {
            console.log('Ready with Device ID', device_id);
            resolve(player)
        });

        // Not Ready
        player.addListener('not_ready', ({device_id}) => {
            console.log('Device ID has gone offline', device_id);
        });

        // Connect to the player!
        player.connect();
    })
}

export const play = ({
                         spotify_uri,
                         playerInstance: {
                             _options: {
                                 getOAuthToken,
                                 id
                             }
                         }
                     }) => {
    getOAuthToken(access_token => {
        fetch(`https://api.spotify.com/v1/me/player/play?device_id=${id}`, {
            method: 'PUT',
            body: JSON.stringify({uris: [spotify_uri]}),
            headers: auth_headers,
        }).catch(err => `error sending request to play song: ${err}`);
    });
};

export const buildLetterToTrack = (songList, tracks) => {
    const uriToTrack = {};
    for (const track of tracks) {
        uriToTrack[track.uri] = track;
    }

    const letterToTrack = {};
    for (const [letter, uri] of Object.entries(songList)) {
        const track = uriToTrack[uri];
        if (track) {
            letterToTrack[letter] = track;
        }
    }

    return letterToTrack;
}
