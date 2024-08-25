# 🎸 ChordHub API

API for a platform where users can upload, share, and discover chord breakdowns for guitar songs. The API is built with Go, using the Gin web framework and GORM for database management. It features robust user authentication, a role-based access system, and is designed with a layered architecture to ensure scalability and maintainability.

# ✨ Features
**User Authentication**: User authentication using JWT access and refresh tokens.

**Role-Based Access Control**: Granular permissions for regular users and admin users.

**CRUD Operations for Songs & Artists**: Manage songs and associated artists.

**Future Plans**:

- Search Service: Search by song lyrics, titles, and artists.

- Frontend Development: A user-friendly interface for managing and discovering chord breakdowns.

# 🔒 Authentication & Authorization
JWT Tokens
This API uses JWT tokens for secure authentication. After registering or logging in, the user will receive an access token and a refresh token. These tokens must be included in the Authorization header for protected routes.

User Roles
User: Can upload and manage their own songs.
Admin: Has full access to manage all artists, songs, and users.

# 📚 API Endpoints
**Public Routes**
- Register: POST /api/v1/register
- Login: POST /api/v1/login
- Refresh Token: POST /api/v1/refresh
- Get Artists: GET /api/v1/artists
- Get Artist Information: GET /api/v1/artists/:id
- Get Song Information: GET /api/v1/songs/:id

**Protected Routes (Requires Authentication)**
- Get User Info: GET /api/v1/users/me
- Upload Song: POST /api/v1/songs
- Update Song: PUT /api/v1/songs/:id

**Admin Routes (Requires Admin Role)**
- Create Artist: POST /api/v1/artists
- Update Artist: PUT /api/v1/artists/:id
- Delete Artist: DELETE /api/v1/artists/:id
- Create New User: POST /api/v1/users/create
