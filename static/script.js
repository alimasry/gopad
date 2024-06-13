function getUUIDFromURL() {
    const pathname = window.location.pathname;
    const segments = pathname.split('/');
    return segments[2];
}

const uuid = getUUIDFromURL();
const editor = document.getElementById('textEditor');
const undoButton = document.getElementById('undoButton');
const redoButton = document.getElementById('redoButton');
let socket;
let isConnected = false;

function connectWebSocket() {
    socket = new WebSocket('ws://localhost:8080/documents/' + uuid + '/ws');
    socket.onopen = function (event) {
        console.log('WebSocket is connected.');
        retryInterval = 1000;
    };

    socket.onmessage = function (event) {
        const ev = JSON.parse(event.data);
        console.log(ev);
        document.title = ev.Data.Title;
        editor.value = ev.Data.Content;
        console.log('Message from server ', event.data);
    };

    socket.onclose = function (event) {
        console.log('WebSocket closed. Attempting to reconnect...', event.reason);
        setTimeout(() => {
            connectWebSocket();
        }, retryInterval);
        retryInterval = Math.min(2 * retryInterval, 30000);
    };

    socket.onerror = function (event) {
        console.error('WebSocket error observed:', event);
    };
}

connectWebSocket();

const buffer = {
    operation: '',
    position: 0,
    string: '',
    delete: 0,
    timeoutId: null
};

function sendBuffer() {
    if (!buffer.operation) return;

    const payload = JSON.stringify({
        command: buffer.operation,
        uuid: uuid,
        data: {
            position: buffer.position,
            string: buffer.string,
            delete: buffer.delete
        }
    });

    socket.send(payload);
    resetBuffer();
}

function resetBuffer() {
    buffer.operation = '';
    buffer.position = 0;
    buffer.string = '';
    buffer.delete = 0;
    if (buffer.timeoutId) {
        clearTimeout(buffer.timeoutId);
    }
    buffer.timeoutId = null;
}

editor.addEventListener('input', function (event) {
    const { inputType, data } = event;
    const currentPosition = editor.selectionStart;

    if (buffer.timeoutId) {
        clearTimeout(buffer.timeoutId);
    }
    buffer.timeoutId = null;

    if (inputType === 'insertText' || inputType === 'insertFromPaste') {
        handleInsertOperation(data, currentPosition);
    } else if (inputType === 'deleteContentBackward') {
        handleDeleteOperation(currentPosition);
    } else {
        sendBuffer();
    }

    buffer.timeoutId = setTimeout(sendBuffer, 300);
});

editor.addEventListener('keydown', function (event) {
    const currentPosition = editor.selectionStart;
    if (event.key === 'Enter') {
        handleInsertOperation("\n", currentPosition + 1);
    }
});

function handleInsertOperation(data, currentPosition) {
    if (buffer.operation === 'insert_event' && buffer.position + buffer.string.length === currentPosition - 1) {
        buffer.string += data;
    } else {
        sendBuffer();
        buffer.operation = 'insert_event';
        buffer.position = currentPosition - 1;
        buffer.string = data;
        buffer.delete = 0;
    }
}

function handleDeleteOperation(currentPosition) {
    if (buffer.operation === 'delete_event' && buffer.position === currentPosition) {
        buffer.delete++;
    } else {
        sendBuffer();
        buffer.operation = 'delete_event';
        buffer.position = currentPosition;
        buffer.string = '';
        buffer.delete = 1;
    }
}

function handleUndoRedoOperation(operation) {
    sendBuffer();
    buffer.operation = operation;
    sendBuffer();
}

function insertString(original, insert, index) {
    if (insert == undefined) return original;
    return original.substring(0, index) + insert + original.substring(index);
}

function removeFromString(original, count, index) {
    const to = Math.max(index - count + 1, 0);
    return original.substring(0, to) + original.substring(index + 1);
}

undoButton.addEventListener('click', () => handleUndoRedoOperation('undo_event'));
redoButton.addEventListener('click', () => handleUndoRedoOperation('redo_event'));
