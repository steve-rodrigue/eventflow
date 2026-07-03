package https

func EventFlowRuntimeJS() string {
	return `
const socket = new WebSocket("ws://localhost:8080/api");

function trigger(eventName, payload = {}) {
	if (socket.readyState !== WebSocket.OPEN) {
		console.warn("WebSocket is not connected yet");
		return;
	}

	socket.send(JSON.stringify({
		type: "event",
		event: eventName,
		payload
	}));
}

socket.addEventListener("open", () => {
	console.log("WebSocket connected");
});

socket.addEventListener("close", () => {
	console.log("WebSocket disconnected");
});

socket.addEventListener("error", (error) => {
	console.error("WebSocket error", error);
});

socket.addEventListener("message", (message) => {
	const data = JSON.parse(message.data);

	if (data.type === "operations") {
		applyOperations(data.operations);
	}
});

function applyOperations(operations) {
	for (const op of operations) {
		const target = op.target ? document.querySelector(op.target) : null;

		if (op.type === "replace" && target) {
			target.outerHTML = op.html;
		}

		if (op.type === "remove" && target) {
			target.remove();
		}

		if (op.type === "append" && target) {
			target.insertAdjacentHTML("beforeend", op.html);
		}

		if (op.type === "prepend" && target) {
			target.insertAdjacentHTML("afterbegin", op.html);
		}

		if (op.type === "navigate") {
			window.location.href = op.url;
		}
	}
}

document.addEventListener("click", (event) => {
	const element = event.target.closest("[data-event]");

	if (!element) {
		return;
	}

	const payload = {};

	for (const [key, value] of Object.entries(element.dataset)) {
		if (key === "event") {
			continue;
		}

		payload[key] = value;
	}

	trigger(element.dataset.event, payload);
});
`
}
