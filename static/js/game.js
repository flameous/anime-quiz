import * as player from "./player.js";
import * as ws from "./ws.js";

const state = {
    "room_id": ""
};

export function initGame(room_id) {
    state.room_id = room_id;

    player.init("player", () => {
        document.getElementById("content-loader").style.display = "none";
        document.getElementById("content-join-form").style.display = "block";
    });

    ws.setCallbacks({
        onAdminNotify: onAdminNotify,
        onEnterNotify: onEnterNotify,
        onSendVideo: onSendVideo,
        onStartPlaying: onStartPlaying,
        onAnswer: onAnswer,
        onArbitrage: onArbitrage,
        onArbitrageApproved: onArbitrageApproved,
        onGameOver: onGameOver,
    });

    document.getElementById("content-join-form").onsubmit = joinGame;
    document.getElementById("start-button").onclick = startGame;
    document.getElementById("content-quiz-form").onsubmit = sendAnswer;
    document.getElementById("content-header-volume").oninput = function () {
        player.setVolume(parseInt(this.value));
    };
}

export function showLeaderboard(leaderboard) {
    for (let user_id in leaderboard) {
        let li = document.createElement("li");
        li.innerText = user_id + ": " + leaderboard[user_id];

        document.getElementById("content-leaderboard-list").appendChild(li);
    }
}

function joinGame(event) {
    event.preventDefault();

    let user_id = document.getElementById("user-id").value;
    if (user_id === "") {
        return
    }

    ws.initConnection(user_id, state.room_id);

    document.getElementById("content-join-form").style.display = "none";
    document.getElementById("content-quiz").style.display = "block";
}

function onEnterNotify(data) {
    document.getElementById("content-header-counter").innerText = "(в игре: " + data["count"] + ")";

    let message = document.createElement("div");
    message.innerText = data["user_id"] + " подключился";

    document.getElementById("messages").appendChild(message);
    setTimeout(() => {
        message.remove();
    }, 30 * 1000);
}

function onAdminNotify() {
    document.getElementById("start-button").style.display = "inline";
}

function startGame() {
    document.getElementById("start-button").style.display = "none";

    ws.sendUserNotify("startGame");
}

function onSendVideo(data) {
    document.getElementById("content-quiz-status").innerText = "Загружаем звук";
    document.getElementById("content-quiz-arbitrage").innerHTML = "";

    player.loadVideo(data["video_id"], data["start"], () => {
        ws.sendUserNotify("")
    })
}

function onStartPlaying() {
    document.getElementById("content-quiz-status").innerText = "Угадываем";
    document.getElementById("content-quiz-form").style.display = "block";
    document.getElementById("content-quiz-input").focus();

    player.playVideo();
}

function sendAnswer(event) {
    event.preventDefault();

    let answer = document.getElementById("content-quiz-input").value;
    if (answer === "") {
        return;
    }
    document.getElementById("content-quiz-input").value = "";

    document.getElementById("content-quiz-form").style.display = "none";

    ws.sendAnswer(answer);
}

function onAnswer(answer) {
    document.getElementById("content-quiz-status").innerText = "Правильный ответ: " + answer;
    document.getElementById("content-quiz-form").style.display = "none";
}

function onArbitrage(data) {
    let li = document.createElement("li");
    let button = document.createElement("button");

    button.innerText = data["user_id"] + ": " + data["answer"];
    button.dataset.userID = data["user_id"];
    button.onclick = sendArbitrage;

    li.appendChild(button);

    document.getElementById("content-quiz-arbitrage").appendChild(li);
}

function sendArbitrage() {
    let user_id = this.dataset.userID;
    this.parentElement.remove();

    ws.sendArbitrage(user_id)
}

function onArbitrageApproved() {
    let message = document.createElement("div");
    message.innerText = "Ваш ответ принят арбитражом";

    document.getElementById("messages").appendChild(message);
    setTimeout(() => {
        message.remove();
    }, 30 * 1000);
}

function onGameOver(leaderboard) {
    document.getElementById("content-quiz").style.display = "none";
    player.stopVideo();

    showLeaderboard(leaderboard)
}
