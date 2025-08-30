// DOM Elements
const video = document.getElementById('video');
const canvas = document.getElementById('overlay');
const statusEl = document.getElementById('status');
const timeEl = document.getElementById('time');
const toggleBtn = document.getElementById('toggle-tracking-btn');
const playIcon = document.getElementById('play-icon');
const pauseIcon = document.getElementById('pause-icon');
const btnText = document.getElementById('btn-text');

// State Variables
const API_URL = 'https://localhost:8443/update-time';
let totalTimeWatched = 0;
let isActivelyFocused = false;
let isTrackingActive = false;
let detectionInterval;

// Heuristic function to determine if the user is looking forward.
function isLookingForward(landmarks) {
    if (!landmarks) return false;
    const nose = landmarks.getNose();
    const leftEye = landmarks.getLeftEye();
    const rightEye = landmarks.getRightEye();
    const leftEyeCenterX = leftEye.map(p => p.x).reduce((a, b) => a + b) / leftEye.length;
    const rightEyeCenterX = rightEye.map(p => p.x).reduce((a, b) => a + b) / rightEye.length;
    const noseTipX = nose[3].x;
    const tolerance = (rightEyeCenterX - leftEyeCenterX) * 0.35;
    return noseTipX > (leftEyeCenterX - tolerance) && noseTipX < (rightEyeCenterX + tolerance);
}

function updateButtonUI() {
    playIcon.classList.toggle('hidden', isTrackingActive);
    pauseIcon.classList.toggle('hidden', !isTrackingActive);
    btnText.textContent = isTrackingActive ? 'Pause Tracking' : 'Start Tracking';
    if (!isTrackingActive) {
        statusEl.textContent = '⏸️ Paused';
    }
}

async function updateTimeOnServer() {
    if (totalTimeWatched < 1) return;
    try {
        const formData = new URLSearchParams();
        formData.append('time', Math.floor(totalTimeWatched));
        await fetch(API_URL, {
            method: 'POST',
            headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
            body: formData,
        });
    } catch (error) { console.error("Failed to update server:", error); }
}

async function startApp() {
    try {
        statusEl.textContent = '🧠 Loading AI models...';
        await Promise.all([
            faceapi.nets.ssdMobilenetv1.loadFromUri('https://cdn.jsdelivr.net/npm/@vladmandic/face-api/model'),
            faceapi.nets.faceLandmark68Net.loadFromUri('https://cdn.jsdelivr.net/npm/@vladmandic/face-api/model'),
        ]);

        statusEl.textContent = '🎥 Requesting camera...';
        const stream = await navigator.mediaDevices.getUserMedia({ video: {} });
        video.srcObject = stream;
    } catch (err) {
        statusEl.textContent = '⛔ Camera or model error.';
        btnText.textContent = 'Error';
        return;
    }

    video.addEventListener('play', () => {
        toggleBtn.disabled = false;
        isTrackingActive = false;
        updateButtonUI();
        statusEl.textContent = 'Ready to track';

        setInterval(() => {
            if (isActivelyFocused && isTrackingActive) {
                totalTimeWatched += 0.1;
                timeEl.textContent = Math.floor(totalTimeWatched);
            }
        }, 100);

        setInterval(updateTimeOnServer, 5000);

        detectionInterval = setInterval(async () => {
            if (!isTrackingActive) {
                isActivelyFocused = false;
                return;
            }

            const detections = await faceapi.detectSingleFace(video, new faceapi.SsdMobilenetv1Options()).withFaceLandmarks();
            
            if (detections) {
                // The drawing line that was here has been removed.
                
                if (isLookingForward(detections.landmarks)) {
                    isActivelyFocused = true;
                    statusEl.textContent = '✅ Actively Focused';
                    statusEl.classList.add('text-green-600');
                } else {
                    isActivelyFocused = false;
                    statusEl.textContent = '👀 Look at the screen';
                    statusEl.classList.remove('text-green-600');
                }
            } else {
                isActivelyFocused = false;
                statusEl.textContent = '👀 Look at the screen';
                statusEl.classList.remove('text-green-600');
            }
        }, 300);
    });
}

toggleBtn.addEventListener('click', () => {
    isTrackingActive = !isTrackingActive;
    updateButtonUI();
});

// Start the application
startApp();