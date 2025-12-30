# Admin Panel

## Overview
A production-grade Admin Panel platform with:
- **Backend**: Go (Golang) with chi router
- **Database**: PostgreSQL with pgx driver
- **Frontend**: Next.js 14+ App Router with Tailwind CSS

## Project Structure

```
.
├── backend/                    # Go backend
│   ├── cmd/server/            # Main entry point
│   ├── internal/
│   │   ├── config/            # Configuration management
│   │   ├── database/          # Database connection
│   │   ├── handlers/          # HTTP handlers
│   │   ├── middleware/        # Auth, logging, rate limiting
│   │   ├── models/            # Data models
│   │   ├── repository/        # Database queries
│   │   ├── services/          # Business logic
│   │   └── utils/             # Utilities
│   └── migrations/            # SQL migrations
├── frontend/                  # Next.js frontend
│   └── src/
│       ├── app/               # App Router pages
│       ├── components/        # UI components
│       └── lib/               # API client, auth, utils
└── replit.md                  # This file
```

## Features

### Authentication
- JWT-based authentication with access/refresh tokens
- Secure password hashing with Argon2id (new admin passwords) and bcrypt (backward compatible)
- Rate limiting on auth endpoints
- Session management
- Separate admin authentication table for enhanced security

### Authorization (RBAC)
- Role-based access control
- Fine-grained permissions (resource:action)
- Permission caching for performance

### Admin Features
- Dashboard with KPIs
- User management (CRUD, search, status, password reset, set-to-admin)
- Role & Permission management
- Audit logs (immutable, searchable, CSV export)
- Feature flags management (CRUD, toggle)
- Settings page with dynamic feature flags UI

### Technical Highlights
- **Backend**: chi router, pgx with connection pooling, prepared statements
- **Frontend**: React Server Components, TanStack Query, debounced search
- **Security**: CSRF/XSS protection, secure cookies, rate limiting
- **Multi-tenant ready**: tenant_id on all tables

## Default Credentials
- **Email**: admin@example.com
- **Password**: Admin123!

## API Endpoints

### Auth
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh tokens
- `GET /api/v1/auth/me` - Current user

### Users
- `GET /api/v1/users` - List users (paginated)
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `POST /api/v1/users/:id/reset-password` - Reset password

### Roles
- `GET /api/v1/roles` - List roles
- `POST /api/v1/roles` - Create role
- `GET /api/v1/roles/:id` - Get role
- `PUT /api/v1/roles/:id` - Update role
- `DELETE /api/v1/roles/:id` - Delete role

### Audit & Dashboard
- `GET /api/v1/audit-logs` - List audit logs
- `GET /api/v1/audit-logs/export` - Export audit logs to CSV
- `GET /api/v1/dashboard/stats` - Dashboard statistics

### Feature Flags
- `GET /api/v1/feature-flags` - List feature flags
- `POST /api/v1/feature-flags` - Create feature flag
- `GET /api/v1/feature-flags/:id` - Get feature flag
- `PUT /api/v1/feature-flags/:id` - Update feature flag
- `DELETE /api/v1/feature-flags/:id` - Delete feature flag
- `POST /api/v1/feature-flags/:id/toggle` - Toggle feature flag

### Admin Management
- `POST /api/v1/users/:id/set-admin` - Set/unset user as admin
- `GET /api/v1/users/:id/admin-status` - Get admin status

## Environment Variables
- `DATABASE_URL` - PostgreSQL connection string
- `SESSION_SECRET` - JWT signing secret
- `ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins (default: "*")
- `APP_ENV` - Environment (development/production)

## Recent Changes
- Dec 30, 2025: Added CSV export endpoint for audit logs
- Dec 30, 2025: Implemented feature flags API and frontend management UI
- Dec 30, 2025: Added admin authentication table with Argon2id password hashing
- Dec 30, 2025: Added set-to-admin endpoint with password setting
- Dec 30, 2025: Updated Settings page with dynamic feature flags
- Dec 30, 2025: Added RBAC permission middleware to all protected routes
- Dec 30, 2025: Improved CORS configuration with configurable origins
- Dec 30, 2025: Fixed Next.js allowedDevOrigins for Replit environment
- Initial setup with complete backend and frontend
- Database schema with indexes for performance
- JWT authentication with refresh tokens
- RBAC implementation with permission caching
