const userMessageTypeHandShake = "userHandShake";
const userMessageTypeNotify = "userNotify";
const userMessageTypeAnswer = "userAnswer";
const userMessageTypeArbitrage = "userArbitrage";

const serverMessageTypeAdminNotify = "serverAdminNotify";
const serverMessageTypeEnterNotify = "serverEnterNotify";
const serverMessageTypeSendVideo = "serverSendVideo";
const serverMessageTypeStartPlaying = "serverStartPlaying";
const serverMessageTypeAnswer = "serverAnswer";
const serverMessageTypeArbitrage = "serverArbitrage";
const serverMessageTypeArbitrageApproved = "serverArbitrageResult";
const serverMessageTypeGameOver = "serverGameOver";

let ws;

let onAdminNotify;
let onEnterNotify;
let onSendVideo;
let onStartPlaying;
let onAnswer;
let onArbitrage;
let onArbitrageApproved;
let onGameOver;

export function initConnection(user_id, room_id) {
    let protocol = "ws://";

    if (window.location.protocol === "https:") {
        protocol = "wss://";
    }

    ws = new WebSocket(protocol + window.location.host + "/ws");

    ws.onopen = onOpen.bind({
        room_id: room_id,
        user_id: user_id
    });

    ws.onmessage = onMessage;
}

export function setCallbacks(callbacks) {
    onAdminNotify = callbacks.onAdminNotify;
    onEnterNotify = callbacks.onEnterNotify;
    onSendVideo = callbacks.onSendVideo;
    onStartPlaying = callbacks.onStartPlaying;
    onAnswer = callbacks.onAnswer;
    onArbitrage = callbacks.onArbitrage;
    onArbitrageApproved = callbacks.onArbitrageApproved;
    onGameOver = callbacks.onGameOver;
}

export function sendUserNotify(message) {
    send({
        message_type: userMessageTypeNotify,
        message: message
    });
}

export function sendAnswer(answer) {
    send({
        message_type: userMessageTypeAnswer,
        message: answer,
    });
}

export function sendArbitrage(user_id) {
    send({
        message_type: userMessageTypeArbitrage,
        message: user_id,
    });
}

function send(msg) {
    ws.send(JSON.stringify(msg));
}

function onOpen() {
    ws.send(JSON.stringify({
        room_id: this.room_id,
        message: this.user_id,
        message_type: userMessageTypeHandShake
    }));
}

function onMessage(event) {
    const data = JSON.parse(event.data);
    console.log(data);

    switch (data.message_type) {
        case serverMessageTypeAdminNotify:
            onAdminNotify();
            break;

        case serverMessageTypeSendVideo:
            onSendVideo(data.message);
            break;

        case serverMessageTypeEnterNotify:
            onEnterNotify(data.message);
            break;

        case serverMessageTypeStartPlaying:
            onStartPlaying();
            break;

        case serverMessageTypeAnswer:
            onAnswer(data.message);
            break;

        case serverMessageTypeArbitrage:
            onArbitrage(data.message);
            break;

        case serverMessageTypeArbitrageApproved:
            onArbitrageApproved();
            break;

        case serverMessageTypeGameOver:
            onGameOver(data.message);
            break;
    }
}