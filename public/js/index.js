'use strict';

import {
    buildLetterToTrack,
    createPlayer,
    getTracks,
    getUserID,
    initSpotify,
    play,
} from "./modules/spotify.js";
import {getSongList} from "./modules/bajid.js";

const e = React.createElement;

(async function () {
    let player;
    let device_id;
    try {
        await initSpotify();
        let res = await createPlayer();
        player = res.player;
        device_id = res.device_id;
    } catch (e) {
        if (e === 'Authentication failed') {
            window.location.replace('/login');
        } else {
            throw e
        }
    }
    const userId = await getUserID();
    console.log("userId", userId);
    const songList = await getSongList(userId);
    console.log("songList", songList);
    const uris = Object.values(songList);
    console.log("uris", uris);
    const tracks = await getTracks(uris);
    console.log("tracks", tracks);
    const letterToTrack = buildLetterToTrack(songList, tracks);
    console.log("letterToTrack", letterToTrack);

    const app = document.querySelector('#app');

    ReactDOM.render(e(App, {letterToTrack, player, device_id}, null), app)
})()

function App({letterToTrack, player, device_id}) {
    const [state, setState] = React.useState({
        letter: 'a',
        track: letterToTrack['a']
    });

    onkeypress = ({key}) => {
        if (key.charCodeAt(0) > 96 && key.charCodeAt(0) < 123) {
            const letter = key;
            console.log(letter);
            const track = letterToTrack[letter];
            setState({letter, track});
        }
    };

    React.useEffect(() => {
        if (state.track) {
            play({spotify_uri: state.track.uri, playerInstance: player, device_id});
        }
    }, [state.track && state.track.uri]);

    return e(Track, state, null);
}

function Track(props) {
    if (!props.track) {
        return e('div', null, null)
    }

    return e(
        'div',
        {className: 'track'},
        e(
            'div',
            {className: 'artist'},
            e(
                'p',
                {className: 'letter'},
                props.letter.toUpperCase()
            ),
            e(
                'p',
                {className: 'artistName'},
                props.track.artists[0].name + '!'
            )
        ),
        e(
            'div',
            {className: 'song'},
            e(
                'img',
                {className: 'trackImage', src: props.track.album.images[0].url}
            ),
            e(
                'p',
                {className: 'songName'},
                props.track.name
            )
        )
    )
}
