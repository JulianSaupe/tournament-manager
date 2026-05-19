docker_compose('./docker/docker-compose.yml')

dc_resource('engine-database', labels=["DB"])
dc_resource('identity-service-database', labels=["DB"])

local_resource(
    'API',
    serve_cmd='cd backend/engine/cmd/api && go run main.go',
    resource_deps=['engine-database'],
    labels=['Backend'],
    deps=['backend/engine']
)

local_resource(
    'identity-service',
    serve_cmd='cd backend/identity-service/src && cargo run',
    resource_deps=['identity-service-database'],
    labels=['Backend'],
    deps=[
         'backend/identity-service/src',
         'backend/identity-service/Cargo.toml',
         'backend/identity-service/Cargo.lock',
         'backend/identity-service/migrations',
    ]
)

# Define the frontend service with hot reload
local_resource(
    'frontend',
    serve_cmd='cd frontend && npm run dev',
    deps=['frontend/src'],
    resource_deps=['engine-database'],
    labels=['Frontend'],
    links=[
        link('http://localhost:5173', 'Frontend App')
    ]
)

print("""
Tournament Manager Development Environment

Services:
- PostgreSQL Database:
    - Engine: localhost:5432
    - Identity Service: localhost:5433
- API: http://localhost:3000
- Frontend: http://localhost:5173
- Identity Service: http://localhost:5000

Commands:
- tilt up: Start all services
- tilt down: Stop all services
""")
