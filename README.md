# ZoneSkill

**ZoneSkill** is an AI-powered job-search assistant that helps users discover relevant employment opportunities, tailor their CVs, receive personalized career guidance, and manage job applications. It integrates Google Gemini with PostgreSQL and pgvector to support AI-assisted career services and persistent memory.

## 1. Problem Statement and Proposed Solution

### Problem Statement

Job seekers, particularly students, recent graduates, and self-taught professionals, often struggle to find jobs that match their skills and experience. They may also find it difficult to tailor their CVs to specific job descriptions, understand job requirements, and track applications across different platforms. Traditional job boards primarily focus on listing vacancies and may not provide the personalized assistance needed throughout the job-search process.

### Proposed Solution

ZoneSkill brings job discovery, resume management, AI-powered CV tailoring, conversational career assistance, and application tracking into one platform. By integrating Google Gemini and a PostgreSQL-based memory layer, ZoneSkill aims to provide a more organized and personalized job-search experience.

### Key Features

- **Job Search:** Search for and explore employment opportunities.
- **Job Recommendations:** Discover opportunities through the platform's recommendation functionality.
- **Resume Management:** Upload and manage resumes in supported formats.
- **AI-Powered CV Tailoring:** Use Google Gemini to improve CV content for specific job descriptions.
- **Conversational AI Assistant:** Receive AI-assisted career guidance and job-search support.
- **Persistent Memory:** Maintain conversational context through the backend memory layer.
- **Saved Jobs:** Save interesting opportunities for later.
- **Application Tracking:** Organize and monitor job applications.

## 2. Technology Stack

| Layer | Technology | Purpose |
|---|---|---|
| Frontend | React, Vite, TypeScript | User interface and frontend development |
| Backend | Go (`net/http`) | REST API and application logic |
| Database | PostgreSQL | Persistent application data |
| Vector Search | pgvector | Vector storage and similarity search |
| AI Model | Google Gemini | Career assistance and CV tailoring |
| AI API | Gemini `generateContent` REST API | AI model integration |
| Authentication | JWT (HS256) | Authentication and protected API access |
| File Storage | Local filesystem | Resume storage |
| Version Control | Git and GitHub | Source control and collaboration |
| Testing | Go testing tools | Unit and integration testing |

## 3. Setup and Installation Instructions

### Prerequisites

- Go 1.24 or later
- Node.js 18 or later and npm
- PostgreSQL 16 or later with the `pgvector` extension
- Git
- A Google Gemini API key for AI features

### Clone the Repository

```bash
git clone https://github.com/Ashomondi/Zone-skill.git
cd Zone-skill
```

### Step 1: Set Up the Database

Run the provided database setup script:

```bash
./scripts/setup_postgres.sh
```

Alternatively, create the database manually:

```sql
CREATE DATABASE zoneskill;
```

Ensure PostgreSQL is running and pgvector is available. The backend migration runner applies database migrations on startup and tracks them in the `schema_migrations` table. Migration `003` enables pgvector.

### Step 2: Configure and Run the Backend

```bash
cd backend
cp .env.example .env
```

Configure the environment variables in `backend/.env`:

| Variable | Required | Default | Description |
|---|---|---|---|
| `PORT` | No | `8080` | Backend HTTP port |
| `DATABASE_URL` | Yes | None | PostgreSQL connection string |
| `JWT_SECRET` | Yes | None | JWT signing secret |
| `CORS_ALLOWED_ORIGIN` | No | `http://localhost:5173` | Allowed frontend origin |
| `STORAGE_DIR` | No | `./data/resumes` | Resume storage directory |
| `GEMINI_API_KEY` | For AI features | None | Google Gemini API key |
| `GEMINI_MODEL` | No | `gemini-2.5-flash` | Gemini model used for AI requests |

Start the backend:

```bash
go run ./cmd/api
```

The backend listens on port `8080` by default. The `/health` endpoint checks database and local-storage availability.

**Note:** The `/api/chat` and `/api/tailor` endpoints are registered only when `GEMINI_API_KEY` is configured. `JWT_SECRET` must be set for the backend to start.

### Step 3: Configure and Run the Frontend

Open another terminal:

```bash
cd frontend
npm install
npm run dev
```

The frontend development server typically runs at `http://localhost:5173`.

The default API base URL is `http://localhost:8080/api`. Configure `VITE_API_BASE_URL` if your backend runs at a different address.

### Step 4: Run Tests

Run backend unit tests:

```bash
cd backend
go test ./...
```

To run integration tests against a configured test database:

```bash
export TEST_DATABASE_URL='postgres://USER:PASSWORD@localhost:5432/zoneskill_test?sslmode=disable'
go test -tags integration ./...
```

Use a dedicated test database rather than a production database.

## 4. AI Methodology and Implementation

ZoneSkill uses generative AI to support career-related conversations and CV tailoring. The Go backend connects the frontend to the Gemini API and manages the processing of AI requests.

### AI-Powered Career Assistance

Users can submit career-related questions and request assistance understanding job requirements or preparing for opportunities. Google Gemini generates responses based on the prompt and context supplied by the backend.

### AI-Powered CV Tailoring

Users can request help adapting their CV content to a specific job description. Gemini generates suggestions intended to improve the relevance and presentation of the user's existing qualifications and experience.

Users should review AI-generated content to ensure that it remains accurate and does not introduce false qualifications or achievements.

### Persistent Memory

PostgreSQL with pgvector provides the database infrastructure for persistent memory and vector similarity-search functionality.

The current implementation maintains a single backend memory stream per user. Multiple named conversations in the frontend are stored in browser `localStorage` and are not separate backend conversation threads.

The exact embedding model and memory retrieval process should be documented according to the implemented code.

### AI Processing Workflow

1. A user submits a question or CV-tailoring request through the React frontend.
2. The frontend sends the request to the Go backend.
3. The backend validates the request and prepares the relevant input and applicable context.
4. The backend sends the AI request to Google Gemini.
5. Gemini generates a response.
6. The backend returns the response to the frontend for the user to review.

The backend is responsible for managing API credentials and communicating with the AI service.

## 5. AI Model Names and API Usage

### Google Gemini

- **Provider:** Google
- **Model:** `gemini-2.5-flash` by default
- **API:** Gemini `generateContent` REST API
- **Purpose:** Conversational career assistance and AI-powered CV tailoring

The model is configurable through the `GEMINI_MODEL` environment variable.

### API Endpoint

The backend uses the following endpoint format:

```http
POST https://generativelanguage.googleapis.com/v1beta/models/{MODEL}:generateContent
Content-Type: application/json
x-goog-api-key: YOUR_GEMINI_API_KEY
```

### Example Request

```json
{
  "contents": [
    {
      "parts": [
        {
          "text": "Help me tailor my CV for a backend developer position."
        }
      ]
    }
  ]
}
```

This is an illustrative request showing how text can be submitted to Gemini. The actual prompts and request-handling logic are defined by the backend implementation.

### Configuration

Add the following to the backend environment file:

```env
GEMINI_API_KEY=your_gemini_api_key
GEMINI_MODEL=gemini-2.5-flash
```

Keep the API key private. Do not commit it to GitHub or expose it in frontend code.

### PostgreSQL and pgvector

PostgreSQL stores the application's structured data, while pgvector enables vector storage and similarity searches. This provides the infrastructure for the persistent memory layer.

pgvector is a database extension, not an AI model. The embedding model and retrieval strategy should be documented separately if implemented.

## 6. Product Screenshots

Screenshots should demonstrate the actual running application and its implemented features.

Create a directory named `docs/screenshots/` and add screenshots such as:

| Screenshot | Suggested File |
|---|---|
| Main dashboard | `docs/screenshots/dashboard.png` |
| Job search interface | `docs/screenshots/job-search.png` |
| AI assistant conversation | `docs/screenshots/ai-assistant.png` |
| CV tailoring interface | `docs/screenshots/cv-tailoring.png` |
| Saved jobs and applications | `docs/screenshots/applications.png` |

Once the images have been added to the repository, display them in this section using Markdown:

```markdown
### Dashboard
![ZoneSkill dashboard](docs/screenshots/dashboard.png)

### Job Search
![ZoneSkill job search](docs/screenshots/job-search.png)

### AI Career Assistant
![ZoneSkill AI assistant](docs/screenshots/ai-assistant.png)

### CV Tailoring
![ZoneSkill CV tailoring](docs/screenshots/cv-tailoring.png)
```

Replace or remove any screenshot references for features that are not implemented. Actual screenshots are preferable to mockups for a hackathon submission.

## 7. Individual Team Contributions

ZoneSkill was developed collaboratively, with each team member contributing to different aspects of the project.

### Ashly Omondi — Project Manager / Backend Developer

- Coordinated project activities, organized tasks, and monitored team progress.
- Supported project planning and communication among team members.
- Contributed to backend development and API implementation.
- Supported the integration of application features and resolution of technical issues.

### Sospeter Kinyanjui — Backend Developer / Database Developer

- Contributed to backend functionality and application logic.
- Worked on database design and data management.
- Supported database integration with backend services.
- Contributed to reliable data storage and retrieval.

### Evans Juma — Backend Developer / Integration

- Contributed to backend development and API functionality.
- Worked on integrating application components and services.
- Supported communication between backend components and external services.
- Helped identify and resolve integration issues.

### Emmaculate Jane Akinyi — Frontend Developer / Tester

- Contributed to frontend development and user interface functionality.
- Supported integration between frontend components and backend APIs.
- Tested application features and identified bugs.
- Assisted with debugging and quality assurance.

### Ian Kimani — UI/UX Designer / Full-Stack Developer

- Contributed to user interface and user experience design.
- Worked on making the application intuitive and user-friendly.
- Contributed to frontend and backend implementation.
- Supported feature development and integration across the application.

## Repository Layout

```text
Zone-skill/
├── backend/   Go API server
├── frontend/  React + Vite + TypeScript application
├── docs/      Architecture, API, database and contributor documentation
├── scripts/   Database setup scripts
└── README.md
```

## API Documentation

See the [API Reference](https://github.com/Ashomondi/Zone-skill/blob/main/docs/API.md) for endpoint details.

The documented API functionality includes health checks, authentication, resume management, job search, recommendations, saved jobs, applications, chat, and CV tailoring.

## Additional Documentation

- [Architecture](https://github.com/Ashomondi/Zone-skill/blob/main/docs/ARCHITECTURE.md)
- [Database Documentation](https://github.com/Ashomondi/Zone-skill/blob/main/docs/DATABASE.md)
- [API Reference](https://github.com/Ashomondi/Zone-skill/blob/main/docs/API.md)
- [Contributor Guidelines](https://github.com/Ashomondi/Zone-skill/blob/main/docs/CONTRIBUTORS.md)

## Known Limitations

- The backend currently maintains one memory stream per user rather than separate backend threads for named conversations.
- Multiple named frontend conversations are stored in browser `localStorage`.
- Resume parsing supports PDF, DOCX, and TXT. Legacy DOC files are not supported by the parser.
- Resume uploads are limited to 5 MB.
- AI chat and CV tailoring require a configured Gemini API key.
- Resume files are stored on the local filesystem, which may require a different storage solution for production deployment.
