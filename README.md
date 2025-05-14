# midpoint.place API

[![codecov](https://codecov.io/gh/championswimmer/api.midpoint.place/graph/badge.svg?token=PYTQV9APHD)](https://codecov.io/gh/championswimmer/api.midpoint.place)
[![Go Report Card](https://goreportcard.com/badge/github.com/championswimmer/api.midpoint.place)](https://goreportcard.com/report/github.com/championswimmer/api.midpoint.place)


[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
[![Fiber](https://img.shields.io/badge/Fiber-292E33?style=flat&logo=fiber&logoColor=white)](https://github.com/gofiber/fiber)
[![GORM](https://img.shields.io/badge/GORM-00ADD8?style=flat&logo=go&logoColor=white)](https://gorm.io)
[![SQLite](https://img.shields.io/badge/SQLite-003B57?style=flat&logo=sqlite&logoColor=white)](https://www.sqlite.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=flat&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Swagger 2.0](https://img.shields.io/badge/Swagger-2.0-85EA2D?style=flat&logo=swagger)](https://swagger.io/)


## Project Description

The midpoint.place API is a service designed to help groups of people find the optimal meeting location based on their individual positions. It provides a robust set of features for managing groups and calculating central meeting points:

### Key Features:
- **User Management**: Secure user registration and authentication system
- **Group Creation**: Users can create groups with different visibility levels (public, protected, or private)
- **Group Membership**: Users can join groups using either group IDs or unique group codes
- **Location Tracking**: Members can update their locations within groups
- **Midpoint Calculation**: Automatically calculates the geometric center (midpoint) of all group members' locations
- **Radius Setting**: Groups can specify a search radius for finding suitable meeting places
- **Privacy Controls**: Flexible group privacy settings with optional secret codes for joining

### Use Cases:
- Finding a central meeting point for social gatherings
- Organizing meetups with multiple participants from different locations
- Planning events based on attendees' locations
- Coordinating group activities in a location-aware manner

The API is built with Go using the Fiber framework, providing fast and efficient endpoints for real-time location-based group coordination.

## Source Code Overview

The codebase follows a clean architecture pattern with clear separation of concerns. Here's how the code is organized:

### Models (`src/db/models/`)
- `user.go`: Defines the User entity with authentication and location data
- `group.go`: Represents groups with properties like visibility, radius, and midpoint coordinates
- `group_user.go`: Manages the many-to-many relationship between users and groups, including member locations

### Controllers (`src/controllers/`)
- `users.go`: Handles user registration, authentication, and location updates
- `groups.go`: Manages group creation, updates, and generates unique group codes/secrets
- `group_users.go`: Controls group membership operations and calculates group midpoints

### Routes (`src/routes/`)
- `users.go`: Exposes endpoints for user management (/users/*)
- `groups.go`: Provides group-related endpoints (/groups/*)
  - POST /groups - Create new group
  - GET /groups/{id} - Get group details
  - PUT /groups/{id}/join - Join group
  - DELETE /groups/{id}/join - Leave group

### Security (`src/security/`)
- JWT-based authentication
- Password hashing and verification
- Middleware for protecting routes

### DTOs (`src/dto/`)
- Request/Response objects for API endpoints
- Input validation structures
- Data transformation layer between API and internal models

### Configuration (`src/config/`)
- Environment-based configuration
- Constants and enums for group types and user roles
- Database connection settings

### Utils (`src/utils/`)
- Logging utilities
- Common helper functions
- Error handling utilities
