
# Gopad

Welcome to **Gopad**. This collaborative text editor supports real-time document editing with multiple users. It uses operational transformation to ensure consistent edits and gapbuffer to manage text changes efficiently. Gopad also utilizes webhook hubs for client management and Docker for easy setup and deployment.

## Features

- Operational Transformation: Ensures that all changes are consistently integrated.
- GapBuffer: Optimizes the document editing process.
- Webhook Hubs: Manages client sessions effectively.
- Docker Support: Simplifies setup and deployment.

## Quick Start

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/gopad.git
   cd gopad
   ```

2. **Set up the environment**:
   Copy the `.env_example` file to `.env` and adjust the configuration as needed.
   ```bash
   cp .env_example .env
   ```

3. **Build and run with Docker**:
   ```bash
   docker-compose up --build
   ```

## How to Use

- **View a Document**:
  Access any document using its UUID:
  ```
  localhost:8080/documents/{document_uuid}
  ```
  Replace `{document_uuid}` with the actual UUID of the document you wish to view.

- **Create a New Document**:
  Generate a new document and automatically redirect to its viewing URL:
  ```
  localhost:8080/documents/new
  ```

- **Swagger Documentation**:
  Utilize Swagger to interact with the API and create documents directly:
  ```
  localhost:8080/swagger/index.html
  ```

## Dependencies

- Go (Golang)
- Docker & Docker Compose
