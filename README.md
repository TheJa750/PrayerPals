# PrayerPals

A simple web app for creating groups to post prayer requests.

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
  - [ ] Add endpoint and handler for fetching groups (User Feed)
  - [ ] Add endpoint and handler for fetching posts (Group Feed)
  - [ ] Add endpoint and handler for fetching comments (Post Feed)
  - [x] Add ability for group owners/admins to assign roles to users
  - [x] Add ability for group owners/admins to delete posts
  - [ ] Add ability for group owners/admins to delete groups
  - [ ] Add ability for group owners/admins to kick/ban users
  - [ ] Add ability for users to see groups they belong to
  - [x] Add ability for users to leave groups
    - [x] Add checks for leaving group as owner/admin
    - [x] Add checks for leaving group as last user
  - [ ] Add ability to invite users to group
- [ ] **Build a UI**
  - [ ] Set up project structure for static files (HTML, CSS, JS)
  - [ ] Create a simple homepage with project branding/message
  - [ ] Add user registration form
    - [ ] Connect to backend
  - [ ] Add login form
    - [ ] Connect to backend
    - [ ] Handle tokens
  - [ ] Dashboard: show user’s groups and navigation
  - [ ] Group pages: view group info and posts
  - [ ] Create post form (use API endpoint)
  - [ ] UI for joining/leaving groups
  - [ ] Admin controls for owners (assign roles, delete posts/groups)
  - [ ] Display server/API error messages and loading states
  - [ ] Add styles for basic usability and mobile-responsiveness
- [ ] **World Domination** *(if we get around to it)*
- More coming soon...
