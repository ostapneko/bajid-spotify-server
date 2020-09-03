export const initPlayer = () => {
    window.onSpotifyWebPlaybackSDKReady = () => {
        createPlayer()
            .then((player) => {
                play({
                    playerInstance: player,
                    spotify_uri: 'spotify:track:0thSyp6VHdgcQsSuoYEdOJ',
                });
            })
            .catch((err) => console.error(err))
    };
}

const createPlayer = async () => {
    return new Promise((resolve, reject) => {
        const token = getSpotifyToken();
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

        // Playback status updates
        player.addListener('player_state_changed', state => {
            console.log(state);
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

const getSpotifyToken = () => {
    return document.cookie
        .split(';')
        .find(r => r.startsWith('bajid-spotify-token'))
        .split('=')[1];
};

const play = ({
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
            body: JSON.stringify({ uris: [spotify_uri] }),
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${access_token}`
            },
        }).catch(err => `error sending request to play song: ${err}`);
    });
};
