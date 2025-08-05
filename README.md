# PrayerPals

A simple web app for creating groups to share prayer requests.

## Current To-do/Roadmap

- [x] Get functional server base
- [ ] **Build an API**
  - [x] Add endpoint and handler for checking server health
  - [x] Add endpoint and handler for creating user
    - [x] Add auth package for password hashing
  - [x] Add endpoint and handler for logging in
    - [x] Add token logic to auth package
  - [x] Add endpoint and handler for creating groups
  - [x] Add endpoint and handler for joining groups
  - [x] Add endpoint and handler for creating posts for specific groups
    - [x] Add endpoint and handler for comments on specific posts
  - [x] Add endpoint and handler for fetching groups (User Feed)
  - [x] Add endpoint and handler for fetching posts (Group Feed)
  - [x] Add endpoint and handler for fetching comments (Post Feed)
  - [x] Add ability for group owners/admins to assign roles to users
  - [x] Add ability for group owners/admins to delete posts
  - [x] Add ability for group owners/admins to delete groups
  - [x] Add ability for group owners/admins to kick/ban users
  - [x] Add ability for users to leave groups
    - [x] Add checks for leaving group as owner/admin
    - [x] Add checks for leaving group as last user
  - [ ] Add access/refresh token handling
    - [x] Add ability to refresh access token if refresh token is valid
    - [x] Add endpoint/handler for revoking refresh token
    - [x] Add logic for revoking token on logout
    - [ ] Add logic for revoking token on password change
  - [ ] Add ability to invite users to group
    - [x] Add custom invite code logic
    - [x] Add query for looking up group by invite code/edit join group logic to accomodate
    - [ ] Add method for group admins to send invites to specific users (email is unique)
  - [ ] Add Account features
    - [ ] Add endpoint for changing username
    - [ ] Add endpoint for changing password
- [ ] **Build a UI**
  - [x] Set up project structure for static files (HTML, CSS, JS)
  - [x] Create a simple homepage with project branding/message
  - [x] Add user registration form
    - [x] Connect to backend
  - [x] Add login form
    - [x] Connect to backend
    - [x] Handle tokens
  - [x] Dashboard: show user’s groups and navigation
  - [x] Group pages: view group info and posts
  - [x] Create post form (use API endpoint)
  - [x] UI for joining/leaving groups
  - [ ] Admin controls for owners (assign roles, delete posts/groups)
  - [x] Display server/API error messages and loading states
  - [x] Add styles for basic usability and mobile-responsiveness
- [ ] **World Domination** *(if we get around to it)*
- More coming soon...
