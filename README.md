# NeuroMetric: Clinical-Grade Cognitive Assessment

NeuroMetric is a robust, web-based platform designed for longitudinal cognitive screening. It leverages gamified assessments and AI-powered voice processing to track cognitive health across multiple domains, including memory, attention, executive function, and visuospatial processing.

## 🚀 Project Objective
The primary goal is to provide a "Clinical-Grade" tool for researchers and healthcare providers to manage a sample size of patients (30-50), conduct periodic cognitive assessments, and analyze progress over time through data-driven reporting and graphical visualizations.

## 🛠 Tech Stack

### Backend (The Engine)
- **Language:** Golang (1.21+)
- **Framework:** [Gin Gonic](https://github.com/gin-gonic/gin) (Fast, minimal REST API framework)
- **Database:** SQLite (Zero-config, portable, and efficient for relational data)
- **ORM:** [GORM](https://gorm.io/) (Object-Relational Mapping for secure database queries)
- **Security:** 
    - Stateless **JWT (JSON Web Tokens)** for session management.
    - **Bcrypt** for secure password hashing.

### Frontend (The Interface)
- **Architecture:** Vanilla JavaScript (ES6 Modules) for zero-dependency performance.
- **Styling:** [Tailwind CSS](https://tailwindcss.com/) for a modern "Cream & Ink" clinical aesthetic.
- **Visualizations:** [Chart.js](https://www.chartjs.org/) for longitudinal progress tracking and cognitive radars.
- **State:** Token-based local storage session management.

### AI & Processing
- **Voice Transcription:** [OpenAI Whisper](https://github.com/openai/whisper) (Integrated via Python hybrid subprocess).
- **Semantic Verification:** Simulated **Gemma** logic for high-fidelity response verification.

## 📁 Project Structure
```text
neurometric-go/
├── main.go            # Unified Go server (API + Static File Serving)
├── neurometric.db     # SQLite Database file (Relational storage)
├── go.mod             # Go dependencies
└── public/            # Frontend Assets
    ├── index.html     # Patient Portal (Entry Point)
    ├── dashboard.html # Navigation Hub
    ├── merged_games.html # Core Assessment Suite (10 Modules)
    ├── interpretation.html # Clinical Narrative Report
    ├── reporting.html  # Visual Charts & Radar Profile
    ├── recommendations.html # Personalized Care Pathways
    └── js/
        └── auth.js    # JWT-based Session Manager
```

## 📋 Assessment Modules
1.  **Tile Memory:** Visual sequencing and retention.
2.  **Delayed Recall:** Memory encoding verification.
3.  **Attention Vigilance:** Sustained focus and reaction speed.
4.  **Pathfinder Network:** Executive function and task-switching.
5.  **Recognition Memory:** Longitudinal recall verification.
6.  **Spatial Reasoning:** Visuospatial organization.
7.  **Clock Drawing:** Classical cognitive deficit screening.
8.  **Kitchen Rush (3 Stages):** High-load multi-tasking and processing speed.

## ⚙️ Setup & Installation

1.  **Clone & Navigate:**
    ```bash
    cd neurometric-go
    ```

2.  **Run the Server:**
    ```bash
    go run main.go
    ```

3.  **Access the Portal:**
    Open [http://localhost:8080](http://localhost:8080) in your browser.

## 📈 Clinical Features
- **Stateless Auth:** Secure patient login without cloud dependencies.
- **Relational Data:** Every assessment is strictly linked to a patient ID, allowing for SQL-based progress calculation.
- **Local Sovereignty:** The entire platform (code and data) lives in your project folder, ensuring 100% data control and portability for clinical studies.
