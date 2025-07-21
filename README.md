#  GoPay Lite

GoPay Lite is a full-stack payment application built using **Golang microservices**, **Next.js frontend**, **JWT authentication**, and **Razorpay** for payment processing. 
It demonstrates a clean microservices architecture with Dockerized deployment.

## Tech Stack used:
### FrontEnd
- **Framework**: [Next.js](https://nextjs.org/)
- **Language**: JavaScript / React
- **API Calls**: REST via `fetch`
- **Styling**: CSS Modules

### Backend Microservices (Golang)
- **Auth Service**: JWT login/register, middleware
- **Payment Service**: Razorpay integration for creating orders
- **API Gateway**: Reverse proxy routing via Gorilla Mux

###  Authentication
- JWT tokens stored in `localStorage`  
- Token verification on each protected route

  ### Tools & Infra
- PostgreSQL (DB)
- Docker & Docker Compose
- Swagger Docs
- .env based config

  
##  Folder Structure

gopay-lite/
├── client/
│ └── gopay-lite-frontend/ # Next.js frontend app
├── server/
│ ├── api-gateway/ # API Gateway in Go
│ ├── auth-service/ # Auth microservice
│ └── payment-service/ # Payment microservice (Razorpay)


---

##  Setup & Run (Local Dev)

### 1. Clone Repo

```bash
git clone https://github.com/your-username/gopay-lite.git
cd gopay-lite

```

## Docker Build & Run

Ensure Docker Desktop is installed and running
```
# Auth service
cd server/auth-service
docker build -t auth-service .
docker run --rm -p 8083:8083 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=1234 \
  -e DB_NAME=gopaydb \
  -e JWT_SECRET=your-secret-key \
  auth-service

# Payment service
cd ../payment-service
docker build -t payment-service .
docker run --rm -p 8084:8084 \
  -e DATABASE_URL="postgres://postgres:1234@host.docker.internal:5432/gopaydb?sslmode=disable" \
  -e JWT_SECRET=your-secret-key \
  -e RAZORPAY_KEY=your-razorpay-key \
  -e RAZORPAY_SECRET=your-razorpay-secret \
  payment-service

# API Gateway
cd ../api-gateway
go run main.go

```

## Frontend Setup
```
cd client/gopay-lite-frontend
npm install
# Set environment
cp .env.local.example .env.local
# Edit .env.local to point to API Gateway
npm run dev
```
Frontend runs at: http://localhost:3000
API Gateway runs at: http://localhost:8080

## Features
 Register / Login with JWT
 Razorpay order creation
 Wallet & Transaction display
 Token-based API access
 Reverse proxy routing through API gateway

## API Endpoints

Method	   Endpoint	                    Description
POST	    /api/v1/auth/register       Register a user
POST	    /api/v1/auth/login	        Login a user
GET     	/api/v1/auth/me	              Get user info (protected)
POST	    /api/v1/pay	                Create Razorpay order



# UI 

<img width="1778" height="832" alt="ui3" src="https://github.com/user-attachments/assets/12619f3f-bb0d-4775-988b-38e5528f201a" />






<img width="1826" height="793" alt="ui2" src="https://github.com/user-attachments/assets/ace173b9-a00b-44c4-a6b0-bce3d0d19402" />






<img width="1873" height="801" alt="ui " src="https://github.com/user-attachments/assets/57518d03-a4a4-439f-9599-56c739f7c7b9" />





<img width="1881" height="923" alt="ui5" src="https://github.com/user-attachments/assets/0db8ce63-ed5e-4886-92ca-75e698e40874" />





