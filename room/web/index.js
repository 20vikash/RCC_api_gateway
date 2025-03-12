const host = "192.168.29.27:6969" 

const API_ENDPOINT_CREATE = `http://${host}/createroom`;
const API_ENDPOINT_JOIN = `http://${host}/join`;
const loadingOverlay = document.querySelector('.loading-overlay');
const errorMessage = document.getElementById('errorMessage');

document.getElementById("roomForm").addEventListener("submit", async function(event) {
    event.preventDefault();
    const username = document.getElementById("username").value;
    const isCreateMode = document.getElementById("createTab").classList.contains("active");
    
    if (isCreateMode) {
        await handleRoomCreation(username);
    } else {
        const roomId = document.getElementById("room").value;
        await handleRoomJoin(username, roomId);
    }
});

async function handleRoomCreation(username) {
    try {
        showLoading();
        errorMessage.style.display = 'none';

        const response = await fetch(API_ENDPOINT_CREATE);
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const roomId = await response.text();
        window.location.href = `/code.html?username=${encodeURIComponent(username)}&room=${encodeURIComponent(roomId)}`;
        
    } catch (error) {
        console.error('Error:', error);
        errorMessage.textContent = 'Failed to create room. Please try again.';
        errorMessage.style.display = 'block';
    } finally {
        hideLoading();
    }
}

async function handleRoomJoin(username, roomId) {
    try {
        showLoading();
        errorMessage.style.display = 'none';

        const response = await fetch(`${API_ENDPOINT_JOIN}?room=${encodeURIComponent(roomId)}`);

        if (!response.ok) {
            errorMessage.textContent = 'Room not found.';
            errorMessage.style.display = 'block';
            return;
        }

        window.location.href = `/code.html?username=${encodeURIComponent(username)}&room=${encodeURIComponent(roomId)}`;
    } catch (error) {
        console.error('Error:', error);
        errorMessage.textContent = 'An error occurred. Please try again.';
        errorMessage.style.display = 'block';
    } finally {
        hideLoading();
    }
}

function showLoading() {
    loadingOverlay.style.display = 'flex';
    document.body.style.pointerEvents = 'none';
}

function hideLoading() {
    loadingOverlay.style.display = 'none';
    document.body.style.pointerEvents = 'auto';
}

document.getElementById("createTab").addEventListener("click", function() {
    document.getElementById("createTab").classList.add("active");
    document.getElementById("joinTab").classList.remove("active");
    
    const roomInput = document.getElementById("roomInput");
    roomInput.style.display = "none";
    roomInput.querySelector("input").removeAttribute("required");

    document.getElementById("form-title").innerText = "Create Room";
});

document.getElementById("joinTab").addEventListener("click", function() {
    document.getElementById("joinTab").classList.add("active");
    document.getElementById("createTab").classList.remove("active");
    
    const roomInput = document.getElementById("roomInput");
    roomInput.style.display = "block";
    roomInput.querySelector("input").setAttribute("required", "true");

    document.getElementById("form-title").innerText = "Join Room";
});
