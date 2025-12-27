# Load Docker Compose for the database
docker_compose('./docker/docker-compose.yml')

dc_resource('postgres', labels=["Backend"])
dc_resource('identity-service-database', labels=["Identity-Service"])

# Define the backend service
local_resource(
    'API',
    serve_cmd='cd backend/cmd/api && go run main.go',
    resource_deps=['postgres'],  # Ensure database is running before frontend
    labels=['Backend'],
    deps=['backend/']
)

# Define the frontend service with hot reload
local_resource(
    'frontend',
    serve_cmd='cd frontend && npm run dev',
    deps=['frontend/src'],
    resource_deps=['postgres'],  # Ensure database is running before frontend
    labels=['Frontend'],
    links=[
        link('http://localhost:5173', 'Frontend App')
    ]
)

# Display helpful information
print("""
Tournament Manager Development Environment

Services:
- PostgreSQL Database:
    - backend: localhost:5432
    - identity service: localhost:5433
- API: http://localhost:3000
- Frontend: http://localhost:5173
- Identity Service: http://localhost:5000

Commands:
- tilt up: Start all services
- tilt down: Stop all services
""")
