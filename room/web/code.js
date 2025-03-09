let me = true;
let monacoReady;

const runButton = document.getElementById('playButton');
const terminalContainer = document.getElementById('terminalContainer');

const term = new Terminal({
    cursorBlink: false,
    disableStdin: true,
    theme: {
        background: '#1e1e2e',
        foreground: '#e5e7eb',
        cursor: '#00d4ff',
        selection: 'rgba(59,130,246,0.3)'
    },
    fontSize: 14,
    fontFamily: 'JetBrains Mono, Menlo, Monaco, Consolas, monospace',
    letterSpacing: 0.8,
    lineHeight: 1.2,
    cols: 80,
    rows: 15,
    scrollback: 1000,
    screenReaderMode: true
});

term.open(document.getElementById('terminal'));

runButton.addEventListener('click', async () => {
    await monacoReady;
    const code = window.editor.getValue();
    const terminalContainer = document.getElementById('terminalContainer');
    const language = document.getElementById('language').value;
    const username = "vikash";

    terminalContainer.classList.add('active');
    term.clear();

    try {
        if (language === 'javascript') {
            const originalLog = console.log;
            const originalError = console.error;
            let hasOutput = false;

            console.log = (...args) => {
                hasOutput = true;
                term.writeln(args.join(' '));
            };

            console.error = (...args) => {
                hasOutput = true;
                term.writeln(`\x1b[31m${args.join(' ')}\x1b[0m`);
            };

            eval(code);
            
            if (!hasOutput) {
                term.writeln('\x1b[33mProgram executed successfully - no output generated\x1b[0m');
            }

            console.log = originalLog;
            console.error = originalError;
        } else {
            term.write('\x1b[33mExecuting code...\x1b[0m\r\n');

            const response = await fetch(`http://localhost:6969/output?language=${language}&id=${roomId}&username=${username}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ code: code })
            });

            const data = await response.json();
            
            term.write('\x1b[32mExecution completed\x1b[0m\r\n');
            
            if (response.ok) {
                const output = data.Output.split('\n');
                output.forEach(line => {
                    term.writeln(`\x1b[37m${line}\x1b[0m`);
                });
            } else {
                term.writeln(`\x1b[31mError: ${data.error || 'Unknown error'}\x1b[0m`);
            }
        }
    } catch (err) {
        term.writeln(`\x1b[31mError: ${err.message}\x1b[0m`);
    }
});

const urlParams = new URLSearchParams(window.location.search);
const roomId = urlParams.get("room");

const socket = new WebSocket(`ws://192.168.29.119:6969/ws?roomid=${roomId}`);

document.getElementById("generateCode").addEventListener("click", () => {
    document.getElementById("modal").classList.add("active");
});

document.getElementById("closeModal").addEventListener("click", () => {
    document.getElementById("modal").classList.remove("active");
});

document.getElementById("insertCode").addEventListener("click", async () => {
    await monacoReady;

    const prompt = document.getElementById("generatedCode").value;
    const language = document.getElementById("language").value;

    try {
        const response = await fetch(`http://localhost:6969/generate?prompt=${encodeURIComponent(prompt)}&language=${encodeURIComponent(language)}&id=${encodeURIComponent(roomId)}`);
        const data = await response.json();

        if (data.code) {
            window.editor.setValue(data.code);
            document.getElementById("modal").classList.remove("active");
        } else {
            alert("Failed to generate code. Please try again.");
        }
    } catch (error) {
        console.error("Error fetching generated code:", error);
        alert("Error connecting to the server.");
    }
});

document.getElementById("debugCode").addEventListener("click", async () => {
    await monacoReady;

    const language = document.getElementById("language").value;

    try {
        const response = await fetch(`http://localhost:6969/debug?id=${encodeURIComponent(roomId)}&language=${encodeURIComponent(language)}`, {
            method: "POST",
            body: JSON.stringify({
                code: window.editor.getValue(),
            }),
            headers: {
                "Content-type": "application/json",
            },
        });
        const data = await response.json();

        if (data.code) {
            window.editor.setValue(data.code);
        } else {
            alert("Failed to debug code. Please try again.");
        }
    } catch (error) {
        console.error("Error fetching debugged code:", error);
        alert("Error connecting to the server.");
    }
});

socket.onopen = async () => {
    console.log("Connected to WebSocket");

    await monacoReady;
    socket.send(`${roomId}~load`);
};

socket.onerror = (error) => {
    console.error("WebSocket Error:", error);
};

socket.onmessage = async (event) => {
    console.log("Received:", event.data);

    let rawData;
    try {
        rawData = JSON.parse(event.data);
    } catch (error) {
        console.error("Error parsing JSON:", error, "Data received:", event.data);
        return;
    }

    me = false;
    await monacoReady;

    if (Array.isArray(rawData)) {
        let changes = [];
        for (let item of rawData) {
            if (typeof item === "string") {
                try {
                    let parsedItem = JSON.parse(item);
                    if (Array.isArray(parsedItem)) {
                        changes.push(...parsedItem);
                    } else {
                        console.error("Unexpected inner format:", parsedItem);
                    }
                } catch (err) {
                    console.error("Error parsing inner JSON:", err, "Data:", item);
                }
            } else if (typeof item === "object") {
                changes.push(item);
            }
        }
        applyChanges(changes);
    } else if (typeof rawData === "object") {
        applyChanges(rawData);
    }
};

function sendChanges(changes) {
    socket.send(`${roomId}~${JSON.stringify(changes)}`);
}

async function applyChanges(changes) {
    await monacoReady;

    console.log("Applying changes:", changes);

    window.editor.executeEdits("remote", changes.map(change => ({
        range: new monaco.Range(
            change.range.startLineNumber,
            change.range.startColumn,
            change.range.endLineNumber,
            change.range.endColumn
        ),
        text: change.text,
        forceMoveMarkers: true
    })));

    window.editor.pushUndoStop();
}

require.config({ paths: { 'vs': 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.44.0/min/vs' } });

monacoReady = new Promise((resolve) => {
    require(["vs/editor/editor.main"], function () {
        window.editor = monaco.editor.create(document.getElementById("editor"), {
            value: "// Start coding...\nprint('Hello world')",
            language: "python",
            theme: "vs-dark",
            fontSize: 16,
            minimap: { enabled: false },
            automaticLayout: true
        });

        document.getElementById("language").addEventListener("change", function (event) {
            let newLang = event.target.value;
            monaco.editor.setModelLanguage(window.editor.getModel(), newLang);
        });

        window.editor.onDidChangeModelContent((event) => {
            if (me) {
                let changes = event.changes.map(change => ({
                    range: {
                        startLineNumber: change.range.startLineNumber,
                        startColumn: change.range.startColumn,
                        endLineNumber: change.range.endLineNumber,
                        endColumn: change.range.endColumn
                    },
                    text: change.text
                }));

                sendChanges(changes);
            } else {
                me = true;
            }
        });

        console.log("Monaco Editor is ready!");
        resolve();
    });
});
