# test-ride-system

Backend for test-ride registrations with:
- Team member auth (`/login`) with a default backend user
- Protected customer entry and OTP verification APIs
- Customer OTP sent via Gmail email

## Environment variables
- `PORT` (default `8080`)
- `MONGO_URI`
- `JWT_SECRET`
- `GMAIL_SENDER_EMAIL`
- `GMAIL_APP_PASSWORD` (Google App Password)

## APIs
- `POST /login` (public, returns token)
- `POST /register` (protected, enter customer details + sends email OTP)
- `POST /verify-otp` (protected, verify customer OTP using email)

## Example protected request
Set header:
`Authorization: Bearer <token>`


## UI
- Open `GET /` in browser for a minimal console UI.
- The UI calls `/login`, `/register`, and `/verify-otp` on the same backend.


## Default login
- Name: `harikeerthan`
- Email: `hariroxx47@gmail`
- Password: `Scotty@123`

No setup key/signup is required.
