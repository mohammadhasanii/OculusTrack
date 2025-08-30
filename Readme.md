# OculusTrack 
An AI-powered web application to track user screen time by analyzing eye gaze with a webcam, built with Go and `face-api.js`.

---

## Live Demo
This demo shows the real-time eye tracking and focus detection in action.

<p align="center">
  <video src="./demo.mp4" width="700" controls autoplay loop muted>
     Your browser does not support the video tag.
  </video>
</p>

---

## ✨ Features
* **Real-time Eye Tracking:** Uses `face-api.js` to detect facial landmarks and determine if the user is looking at the screen.
* **Accurate Time Logging:** A Go backend logs the total time the user is actively focused.
* **Secure Local Server:** Runs on a local HTTPS server with an auto-generated SSL certificate.
* **Modern UI:** A clean, minimalist interface built with Tailwind CSS, featuring a start/pause control.

---

## 🛠️ Tech Stack
* **Backend:** Go (Golang)
* **Frontend:** HTML, Tailwind CSS
* **AI/ML:** `face-api.js` (TensorFlow.js)

---

## 🚀 Getting Started
Follow these steps to get the project running on your local machine.

### 1. Clone the Repository
```bash
git clone https://github.com/your-username/OculusTrack.git
cd OculusTrack
```

### 2. Create Project Files
Create the necessary files and folders with the source code provided below:
- `main.go`
- `static/index.html`
- `static/script.js`

### 3. Run the Go Server
This command will start the server and automatically generate the required SSL certificate files (`localhost.crt` & `localhost.key`) on the first run.

```bash
go run main.go
```

The server will be available at `https://localhost:8443`.

---

## 🔒 Important: Enabling HTTPS for Webcam Access
Modern browsers require a secure `https://` connection to access your webcam for privacy reasons. Since this project uses a self-signed certificate, you must manually instruct your operating system to trust it.

**You only need to do this once.**

### On macOS
1. After running the server for the first time, find the generated `localhost.crt` file in your project folder.
2. Open the **Keychain Access** application.
3. Drag and drop the `localhost.crt` file into the **System** keychain.
4. Find the "localhost" certificate in the list, double-click it.
5. Expand the **"Trust"** section.
6. Change the **"When using this certificate"** dropdown to **"Always Trust"**.
7. Close the window (you may need to enter your password).
8. Restart your browser completely.

### On Windows
1. After running the server for the first time, find the generated `localhost.crt` file.
2. Double-click the `localhost.crt` file.
3. Click the **"Install Certificate..."** button.
4. Select **"Current User"** and click **Next**.
5. Choose **"Place all certificates in the following store"** and click **"Browse..."**.
6. Select the **"Trusted Root Certification Authorities"** store and click **OK**.
7. Click **Next**, then **Finish**. Acknowledge the security warning by clicking **Yes**.
8. Restart your browser completely.

---


---

## 🧠 How It Works
1. **Face Detection:** The frontend uses `face-api.js` to detect faces and facial landmarks from the webcam feed.
2. **Gaze Analysis:** The application analyzes eye positions and orientations to determine if the user is looking at the screen.
3. **Time Tracking:** When eyes are detected as "focused," the system logs active screen time via API calls to the Go backend.
4. **Real-time Feedback:** The UI provides instant visual feedback showing current focus status and accumulated time.

---

## 💡 Usage
1. Open your browser and navigate to `https://localhost:8443`
2. Allow camera permissions when prompted
3. Click "Start Tracking" to begin eye tracking
4. The system will monitor your focus and display real-time statistics
5. Use "Pause/Resume" to control tracking as needed
