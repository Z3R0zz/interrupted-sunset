# Interrupted Sunset

The long awaited download your files and run update for interrupted.me as we sunset the project.

---

## 🏗️ Project Structure

```
backend/    # Go backend API and services
frontend/   # Next.js frontend application
```

- **Backend:** RESTful API built with Go, organized into modules for config, controllers, database, middleware, models, routes, and more.
- **Frontend:** Next.js app with TypeScript, Tailwind CSS, custom API wrapper, and UI enhancements.

---

## 🚀 Getting Started

### 1. Clone the Repository

```sh
git clone https://github.com/Z3R0zz/interrupted-sunset.git
cd interrupted-sunset
```

### 2. Setup the Backend

```sh
cd backend
cp ../frontend/.env.example .env
go mod tidy
go run src/main.go
```

### 3. Setup the Frontend

```sh
cd frontend
cp .env.example .env
npm install
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) to view the frontend.

---

## ✨ Features

- **Go Backend:** Modular, scalable REST API
- **Next.js Frontend:** TypeScript, Tailwind CSS, and modern conventions
- **Custom API Wrapper:** Simplified API requests with authentication support
- **UI Enhancements:** Pre-installed fonts, animations, and UI libraries
- **Authentication Middleware:** Ready for secure endpoints

---

## 📂 Folder Overview

- [`backend/`](backend/) - Go backend source code
- [`frontend/`](frontend/) - Next.js frontend source code

---

## 📝 License

See [`LICENSE`](LICENSE) for license details.
