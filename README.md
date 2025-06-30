# WasaText

Connect with your friends effortlessly using WasaText! Send and receive messages, whether one-on-one
or in groups, all from the convenience of your PC. Enjoy seamless conversations with text or images and
easily stay in touch through your private chats or group discussions.

## Features

- Messaging
- Photo sharing
- User authentication
- Group conversations

## Project Structure

```
wasaProject/
├── webui/           # Vue frontend application
├── cmd/             # Backend API server
├── doc/             # OpenAPI specifications
├── internal/             # Backend application logic
├── compose.yml                # Docker composition file
├── Dockerfile.backend         # Backend container configuration
├── Dockerfile.frontend         # Frontend container configuration
├── Dockerfile.db         # Database container configuration
└── README.md         # This file
```

## Getting Started

### Prerequisites

- Docker
- Docker Compose

### Running the Application

1. Clone the repository:
   ```bash
   git clone https://github.com/FallenKarma/wasaProject.git
   cd wasaProject
   ```

2. Start the application using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. Access the application:
   - Frontend: http://localhost:4173
   - Backend API: http://localhost:8080

### Stopping the Application

```bash
docker-compose down
```

## Technology Stack

- **Frontend**: Modern web framework (Vue.js)
- **Backend**: RESTful API server (Golang)
- **Database**: Relational database (PostgreSQL)
- **Containerization**: Docker & Docker Compose

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request
