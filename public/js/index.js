'use strict';

import {buildLetterToTrack, createPlayer, getTracks, getUserID, initSpotify, play} from "./modules/spotify.js";
import {getSongList} from "./modules/bajid.js";

const e = React.createElement;

(async function () {
    let player;
    try {
        await initSpotify();
        player = await createPlayer();
    } catch (e) {
        if (e === 'Authentication failed') {
            window.location.replace('/login');
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

    ReactDOM.render(e(App, {letterToTrack, player}, null), app)
})()

function App({letterToTrack, player}) {
    const [state, setState] = React.useState({
        letter: 'a',
        track: letterToTrack['a']
    });

    player.addListener('player_state_changed', playerState => {
        console.log('player state change', playerState)
        let trackIsFinished = playerState.paused && (playerState.position > 0)

        if (trackIsFinished) {
            const randomCharCode = Math.floor(Math.random() * (123 - 97) + 97);
            let letter = String.fromCharCode(randomCharCode);
            const track = letterToTrack[letter];
            setState({letter, track});
        }
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
        console.log(state);
        if (state.track) {
            play({spotify_uri: state.track.uri, playerInstance: player});
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
