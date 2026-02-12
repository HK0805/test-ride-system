# test-ride-system

Backend for test-ride registrations with:
- Team member auth (`/login`)
- Protected customer entry and OTP verification APIs
- Customer OTP sent via Gmail email

## Environment variables
- `PORT` (default `8080`)
- `MONGO_URI`
- `JWT_SECRET`
- `GMAIL_SENDER_EMAIL`
- `GMAIL_APP_PASSWORD` (Google App Password)
- `OTP_DEV_MODE` (`true` to bypass SMTP and print OTP in server logs for local testing)

## APIs
- `POST /login` (public, returns JWT)
- `POST /register` (protected, enter customer details + sends email OTP)
- `POST /verify-otp` (protected, verify customer OTP using OTP only)

## Team login credentials
- Name: `harikeerthan`
- Email: `hariroxx47@gmail.com`
- Password: `Scotty@123`

## Example protected request
Set header:
`Authorization: Bearer <token>`

## UI
- Open `GET /` in browser for a minimal console UI.
- The UI calls `/login`, `/register`, and `/verify-otp` on the same backend.
