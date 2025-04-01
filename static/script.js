class CollaborativeEditor {
    constructor() {
        this.uuid = this.getUUIDFromURL();
        this.editor = document.getElementById('textEditor');
        this.undoButton = document.getElementById('undoButton');
        this.redoButton = document.getElementById('redoButton');
        this.socket = null;
        this.retryInterval = 1000;
        this.isProcessingServerUpdate = false;
        
        // Buffer for batching operations
        this.buffer = {
            operation: '',
            position: 0,
            string: '',
            delete: 0,
            timeoutId: null
        };

        this.EDIT_DEBOUNCE_MS = 300;

        this.setupEventListeners();
        this.connectWebSocket();
    }

    getUUIDFromURL() {
        const pathname = window.location.pathname;
        const segments = pathname.split('/');
        return segments[2];
    }

    setupEventListeners() {
        this.editor.addEventListener('input', (e) => {
            if (!this.isProcessingServerUpdate) {
                this.handleInput(e);
            }
        });
        this.editor.addEventListener('keydown', this.handleKeydown.bind(this));
        this.undoButton.addEventListener('click', () => this.handleUndoRedoOperation('undo_event'));
        this.redoButton.addEventListener('click', () => this.handleUndoRedoOperation('redo_event'));
    }

    connectWebSocket() {
        this.socket = new WebSocket(`ws://localhost:8080/documents/${this.uuid}/ws`);
        
        this.socket.onopen = () => {
            console.log('WebSocket connected');
            this.retryInterval = 1000;
        };

        this.socket.onmessage = (event) => {
            const ev = JSON.parse(event.data);
            this.handleServerUpdate(ev);
        };

        this.socket.onclose = (event) => {
            console.log('WebSocket closed. Attempting to reconnect...', event.reason);
            setTimeout(() => this.connectWebSocket(), this.retryInterval);
            this.retryInterval = Math.min(2 * this.retryInterval, 30000);
        };

        this.socket.onerror = (error) => {
            console.error('WebSocket error:', error);
        };
    }

    handleServerUpdate(ev) {
        if (!ev.Data) return;

        this.isProcessingServerUpdate = true;
        
        // Update title if it's changed
        if (ev.Data.Title) {
            document.title = ev.Data.Title;
        }

        // Only update content if it's a sync event
        if (ev.Command === 'sync_event' && ev.Data.Content !== undefined) {
            const currentPosition = this.editor.selectionStart;
            this.editor.value = ev.Data.Content;
            this.editor.selectionStart = this.editor.selectionEnd = currentPosition;
        }
        
        this.isProcessingServerUpdate = false;
    }

    handleInput(event) {
        const { inputType, data } = event;
        const currentPosition = this.editor.selectionStart;

        if (this.buffer.timeoutId) {
            clearTimeout(this.buffer.timeoutId);
        }

        if (inputType === 'insertText' || inputType === 'insertFromPaste') {
            this.handleInsertOperation(data, currentPosition);
        } else if (inputType === 'deleteContentBackward') {
            this.handleDeleteOperation(currentPosition);
        } else {
            this.sendBuffer();
        }

        this.buffer.timeoutId = setTimeout(() => this.sendBuffer(), this.EDIT_DEBOUNCE_MS);
    }

    handleInsertOperation(data, currentPosition) {
        if (this.buffer.operation === 'insert_event' && 
            this.buffer.position + this.buffer.string.length === currentPosition - 1) {
            // Combine with previous insert
            this.buffer.string += data;
        } else {
            this.sendBuffer();
            this.buffer.operation = 'insert_event';
            this.buffer.position = currentPosition - 1;
            this.buffer.string = data;
            this.buffer.delete = 0;
        }
    }

    handleDeleteOperation(currentPosition) {
        if (this.buffer.operation === 'delete_event' && 
            this.buffer.position === currentPosition) {
            // Combine with previous delete
            this.buffer.delete++;
        } else {
            this.sendBuffer();
            this.buffer.operation = 'delete_event';
            this.buffer.position = currentPosition;
            this.buffer.string = '';
            this.buffer.delete = 1;
        }
    }

    handleKeydown(event) {
        if (event.key === 'Enter') {
            const currentPosition = this.editor.selectionStart;
            this.handleInsertOperation("\n", currentPosition + 1);
        }
    }

    handleUndoRedoOperation(operation) {
        this.sendBuffer();
        this.buffer.operation = operation;
        this.sendBuffer();
    }

    sendBuffer() {
        if (!this.buffer.operation) return;

        const payload = JSON.stringify({
            command: this.buffer.operation,
            data: {
                position: this.buffer.position,
                string: this.buffer.string,
                delete: this.buffer.delete
            }
        });

        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            this.socket.send(payload);
        }

        this.resetBuffer();
    }

    resetBuffer() {
        this.buffer.operation = '';
        this.buffer.position = 0;
        this.buffer.string = '';
        this.buffer.delete = 0;
        if (this.buffer.timeoutId) {
            clearTimeout(this.buffer.timeoutId);
        }
        this.buffer.timeoutId = null;
    }
}

// Initialize the editor
const collaborativeEditor = new CollaborativeEditor();
