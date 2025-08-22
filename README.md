# Prayer Pals

Prayer Pals is a web application designed for church groups and faith communities to stay connected throughout the week. Users can create or join private groups, share prayer requests, comment with encouragement or updates, and build supportive fellowship—no matter where they are.

---

## Table of Contents

1. [Screenshots](#screenshots)
2. [Features](#features)
3. [Tech Stack](#tech-stack)
4. [Getting Started](#getting-started)
5. [API Overview](#api-overview)
6. [License](#license)
7. [Contact](#contact)
8. [Security Considerations](#security-considerations)

---

## Screenshots

### User Dashboard

![User dashboard showing groups and relevant buttons](docs/dashboard.png)

### Group Page

![Group page showing prayer request feed](docs/group-feed.png)

### Group Members Modal

![Group members modal showing admin view with moderation/promotion buttons](docs/members-modal.png)

### Post Page

![Post view with comments and comment form](docs/post-with-comments.png)

### Dark Mode

![User dashboard in dark mode](docs/dark-mode.png)

---

## Features

- User authentication with registration and secure JWT/session management
- Create, join, and manage private groups using unique invitation codes
- Post new prayer requests for group members to see and discuss
- Comment on requests with updates, encouragement, or prayers
- Simple, mobile-friendly design using Svelte and Vite
- Group admin tools (kick/ban/promote, change invite code, manage group rules)
- Markdown support for group rules/info
- Modern Go backend (Gorilla Mux, PostgreSQL, SQLC) for reliability and security

---

## Tech Stack

- **Frontend:** Svelte, Vite, JavaScript, CSS (responsive + dark mode)
- **Backend:** Go, Gorilla Mux, SQLC, PostgreSQL
- **Authentication:** JWT, Refresh Tokens, Secure cookies
- **Other:** Docker (planned), RESTful API

---

## Getting Started

> **Local and deployment instructions will be added soon.**

**Backend Requirements:**

- Go 1.21+
- PostgreSQL database  
- Environment variables:  
  - `DB_URL`: Database connection string  
  - `JWT_SECRET`: Secret for signing JWTs  
  - `PLATFORM`: e.g. `"dev"` for access to admin/test routes

**Frontend Requirements:**

- Node.js (LTS Recommended)
- Package manager (npm, yarn, or pnpm)

### Quick Start

```bash
# Backend
cd PrayerPals   # Adjust path as needed
go run main.go  # Or use your preferred build/run process

# (Setup PostgreSQL and apply migrations before running)
# See TODO: detailed instructions below for database setup and env

# Frontend
cd frontend     # Adjust path as needed
npm install
npm run dev     # or: yarn dev / pnpm dev
```

**To Do:**

- Add Docker support or scripts as needed
- Provide step-by-step environment and database setup instructions

---

## API Overview

High-level summary of REST endpoints (all routes are prefixed by `/api`):

| Endpoint                                      | Method | Purpose                                      | Auth  |
| --------------------------------------------- | ------ | -------------------------------------------- | ----- |
| /users                                        | POST   | Register a new user                          | No    |
| /login                                        | POST   | User login                                   | No    |
| /refresh                                      | POST   | Refresh access token                         | Yes   |
| /logout                                       | POST   | Log out user, clear cookies                  | Yes   |
| /users/update                                 | PUT    | Change username or password                  | Yes   |
| /groups                                       | POST   | Create group                                 | Yes   |
| /groups                                       | GET    | List user's groups                           | Yes   |
| /groups/invite/{invite_code}/join             | POST   | Join group with invite code                  | Yes   |
| /groups/invite/{invite_code}                  | GET    | Get group details by invite code             | Yes   |
| /groups/{group_id}                            | GET    | Get group info                               | Yes   |
| /groups/{group_id}/leave                      | DELETE | Leave group                                  | Yes   |
| /groups/{group_id}/posts                      | GET    | List group posts (pagination: limit/offset)  | Yes   |
| /groups/{group_id}/posts                      | POST   | Create post in group                         | Yes   |
| /groups/{group_id}/posts/{post_id}            | DELETE | Delete post (group admin only)                 | Yes   |
| /groups/{group_id}/posts/{post_id}/comments   | GET    | List comments on a post                      | Yes   |
| /groups/{group_id}/posts/{post_id}/comments   | POST   | Add comment to a post                        | Yes   |
| /groups/{group_id}/members                    | GET    | List group members                           | Yes   |
| /groups/{group_id}/members/{user_id}/promote  | PUT    | Promote member to admin (group admin only)   | Yes   |
| /groups/{group_id}/members/{user_id}/moderate | PUT    | Kick/ban member (group admin only)                 | Yes   |
| /groups/{group_id}/invite-code                | PUT    | Change invite code (group admin only)              | Yes   |
| /groups/{group_id}/rules                      | PUT    | Change group rules (group admin only)              | Yes   |
| /groups/{group_id}/description                | PUT    | Change group description (group admin only)        | Yes   |
| /groups/{group_id}/members/{user_id}/remove-content | PUT | Remove all posts by user (group admin only)       | Yes   |

See code for more details on request bodies and expected responses.

---

## License

See [LICENSE.txt](LICENSE.txt) for details.

---

## Contact

Questions? Feedback?  
Please open an issue.

---

**To Do:**

- Add deployment instructions and .env.example
- Add screenshots/demo images above
- Document open TODOs or future security improvements

---

## Security Considerations

This section outlines the key security features currently implemented in Prayer Pals as well as known limitations/future improvements.  
While great care has been taken to protect user data and privacy, this project is still under active development and may not yet meet all best practices required for public deployment.

**Implemented:**

- Passwords are securely hashed with bcrypt
- JWT-based authentication, stored as HttpOnly cookies
- Server-side input validation on user input and forms

**Planned / Not Yet Implemented:**

- Enforce HTTPS and `Secure`/`SameSite` cookie flags in production
- Frontend and backend CSRF and XSS protections
- Rate limiting and abuse-prevention measures
- Fine-grained logging/auditing of moderation actions
- Improved error handling/user feedback for sensitive operations
- Password recovery method
- Automated security audits for dependencies
- Email verification during signup
