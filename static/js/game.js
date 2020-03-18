import * as player from "./player.js";
import * as ws from "./ws.js";
import * as UI from "./UI.js";

const state = {
    "room_id": "",
    "count_quizzes": 0,
    "current_quiz": 0
};

export function initGame(room_id) {
    state.room_id = room_id;

    player.init("player", () => {
        UI.hide(UI.contentLoader);
        UI.show(UI.contentJoinForm);
    });

    ws.setCallbacks({
        onAdminNotify: onAdminNotify,
        onEnterNotify: onEnterNotify,
        onStartGame: onStartGame,
        onSendVideo: onSendVideo,
        onStartPlaying: onStartPlaying,
        onAnswer: onAnswer,
        onArbitrage: onArbitrage,
        onArbitrageApproved: onArbitrageApproved,
        onGameOver: onGameOver,
    });

    UI.contentJoinForm.onsubmit = joinGame;
    UI.startButton.onclick = startGame;
    UI.contentQuizForm.onsubmit = sendAnswer;
    UI.contentHeaderVolume.oninput = function () {
        player.setVolume(parseInt(this.value));
    };
}

export function showLeaderboard(leaderboard) {
    for (let user_id in leaderboard) {
        let li = document.createElement("li");
        li.innerText = user_id + ": " + leaderboard[user_id];

        UI.contentLeaderboardList.appendChild(li);
    }
}

function joinGame(event) {
    event.preventDefault();

    let user_id = document.getElementById("user-id").value;
    if (user_id === "") {
        return
    }

    ws.initConnection(user_id, state.room_id);

    UI.hide(UI.contentJoinForm);
    UI.show(UI.contentQuiz);
}

function onEnterNotify(data) {
    UI.contentHeaderCounter.innerText = `(в игре: ${data["count"]})`;

    UI.showMessage(`${data["user_id"]} подключился`);
}

function onAdminNotify() {
    UI.show(UI.startButton);
}

function startGame() {
    UI.hide(UI.startButton);

    ws.sendUserNotify("startGame");
}

function onStartGame(count_quizzes) {
    state.count_quizzes = count_quizzes;
    UI.updateProgressBar(state.current_quiz, state.count_quizzes)
}

function onSendVideo(data) {
    UI.contentQuizStatus.innerText = "Загружаем звук";
    UI.contentQuizArbitrage.innerHTML = "";

    state.current_quiz++;
    UI.updateProgressBar(state.current_quiz, state.count_quizzes)

    player.loadVideo(data["video_id"], data["start"], () => {
        ws.sendUserNotify("")
    })
}

function onStartPlaying() {
    UI.contentQuizStatus.innerText = "Угадываем";
    UI.show(UI.contentQuizForm);
    UI.contentQuizInput.focus();

    player.playVideo();
}

function sendAnswer(event) {
    event.preventDefault();

    let answer = UI.contentQuizInput.value;
    if (answer === "") {
        return;
    }

    UI.contentQuizInput.value = "";
    UI.hide(UI.contentQuizForm);

    ws.sendAnswer(answer);
}

function onAnswer(answer) {
    UI.contentQuizStatus.innerText = "Правильный ответ: " + answer;
    UI.hide(UI.contentQuizForm);
}

function onArbitrage(data) {
    let li = document.createElement("li");
    let button = document.createElement("button");

    button.innerText = data["user_id"] + ": " + data["answer"];
    button.dataset.userID = data["user_id"];
    button.onclick = sendArbitrage;

    li.appendChild(button);

    UI.contentQuizArbitrage.appendChild(li);
}

function sendArbitrage() {
    let user_id = this.dataset.userID;
    this.parentElement.remove();

    ws.sendArbitrage(user_id)
}

function onArbitrageApproved() {
    UI.showMessage("Ваш ответ принят арбитражом");
}

function onGameOver(leaderboard) {
    UI.hide(UI.contentQuiz);
    player.stopVideo();

    showLeaderboard(leaderboard)
}
