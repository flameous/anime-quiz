export const contentLoader = document.getElementById("content-loader");
export const contentJoinForm = document.getElementById("content-join-form");
export const contentAdminPanel = document.getElementById("content-admin-panel");
export const numRounds = document.getElementById("num-rounds");
export const startButton = document.getElementById("start-button");
export const contentQuizForm = document.getElementById("content-quiz-form");
export const contentHeaderVolume = document.getElementById("content-header-volume");
export const contentLeaderboardList = document.getElementById("content-leaderboard-list");
export const userId = document.getElementById("user-id");
export const contentQuiz = document.getElementById("content-quiz");
export const contentHeaderCounter = document.getElementById("content-header-counter");
export const contentQuizProgress = document.getElementById("content-quiz-progress");
export const contentQuizStatus = document.getElementById("content-quiz-status");
export const contentQuizArbitrage = document.getElementById("content-quiz-arbitrage");
export const contentQuizInput = document.getElementById("content-quiz-input");
export const messages = document.getElementById("messages");

export function hide(el) {
    el.style.display = "none";
}

export function show(el) {
    el.style.display = "initial";
}

export function updateProgressBar(current, all) {
    contentQuizProgress.innerText = `${current} / ${all}`;
}

export function showMessage(text) {
    let message = document.createElement("div");
    message.innerText = text;

    messages.appendChild(message);
    setTimeout(() => {
        message.remove();
    }, 30 * 1000);
}