let player;
let currentVolume = 100;

export function init(player_id, onReady) {
    const tag = document.createElement('script');
    tag.src = "https://www.youtube.com/iframe_api";
    const firstScriptTag = document.getElementsByTagName('script')[0];
    firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

    window.onYouTubeIframeAPIReady = onYouTubeIframeAPIReady.bind({
        player_id: player_id,
        onReady: onReady
    });
}

export function loadVideo(video_id, start, onBuffered) {
    player.loadVideoById(video_id, start, "small");
    player.mute();
    setTimeout(() => {
        player.pauseVideo();
        player.seekTo(start);
        player.unMute();
        player.setVolume(currentVolume);

        onBuffered();
    }, 5000);
}

export function playVideo() {
    player.playVideo();
}

export function stopVideo() {
    player.stopVideo();
}

export function setVolume(volume) {
    currentVolume = volume;
    player.setVolume(volume);
}

function onYouTubeIframeAPIReady() {
    player = new window.YT.Player(this.player_id, {
        height: '200',
        width: '200',
        playerVars: {
            'autoplay': 0,
            'controls': 0,
            'disablekb': 0,
            'fs': 0,
            'iv_load_policy': 3,
            'loop': 0,
            'rel': 0,
            'showinfo': 0,
        },
        events: {
            'onReady': this.onReady
        }
    });
}