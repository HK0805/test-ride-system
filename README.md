# test-ride-system

Backend for test-ride registrations with:
- Team member auth (`/signup`, `/login`)
- Protected customer entry and OTP verification APIs
- Customer OTP sent via Gmail email

## Environment variables
- `PORT` (default `8080`)
- `MONGO_URI`
- `JWT_SECRET`
- `TEAM_SETUP_KEY` (used in `/signup` to create teammate accounts)
- `GMAIL_SENDER_EMAIL`
- `GMAIL_APP_PASSWORD` (Google App Password)

## APIs
- `POST /signup` (public, requires setup key in body)
- `POST /login` (public, returns JWT)
- `POST /register` (protected, enter customer details + sends email OTP)
- `POST /verify-otp` (protected, verify customer OTP using email)

## Example protected request
Set header:
`Authorization: Bearer <token>`


## UI
- Open `GET /` in browser for a minimal console UI.
- The UI calls `/signup`, `/login`, `/register`, and `/verify-otp` on the same backend.
