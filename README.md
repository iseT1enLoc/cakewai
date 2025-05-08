# 🎂 Cakewai – Go Backend API for AI-powered Cake E-commerce

This project is a RESTful backend API built using **Golang**, designed for an e-commerce application that sells cakes online. It incorporates **AI-driven insights and image generation** to enhance the shopping experience, offering personalized cake suggestions, trend analytics, and visual cake inspiration.

Built using a clean modular architecture with **Gin**, **MongoDB**, **Cloudinary**, and **Hugging Face**, this backend is ready for production and scalable deployments.


---

## 🚀 Features

- ✅ RESTful API structure
- 🔐 JWT Authentication (Access token and RefreshToken)
- 🗃️ MongoDB (via official Go driver)
- 🖼️ Image generation using Hugging Face APIs
- ☁️ Cloudinary for image uploading
- 🧠 Fine-grain AI prompting with Gemini
- 📁 Modular folder structure
- 🌱 Environment-based configuration
- 🛠️ Unit Testing support
- 📄 Swagger/OpenAPI support (optional)
- 🐳 Docker-ready
- 🌍 Deployable on [**DigitalOcean**](https://www.digitalocean.com/)

---

## 🧠 AI-Powered Business Analytics

Cakewai includes smart features to **boost sales and user engagement**:
- 🎨 **Visual Cake Generator**:
  - Let users type “a pastel galaxy cake” or “a birthday cake for a dog lover”
  - Hugging Face generates a cake image to inspire or personalize product listings

- 💡 **Fine-grain Prompting with Gemini**:
  - Generate product descriptions, ads, or flavor ideas
  - Auto-reply to customer queries using conversation summaries

---


## 🧰 Tech Stack

- [**Golang** `v1.21+`](https://golang.org/doc/)
- [**Gin** - Web framework](https://github.com/gin-gonic/gin)
- [**MongoDB** - NoSQL database](https://www.mongodb.com/)
- [**JWT** - JSON Web Tokens for Authentication](https://jwt.io/)
- [**Docker** - Containerization platform](https://www.docker.com/)
---

## 📁 Project Structure
cakewai/
├── api/                      # Main entry point(s) for the app
│   └── handlers               # Starts the server
│   └── middlewares               # Starts the server
│   └── routes               # Starts the server
│
├── domain/                   # data model
│
├── infras/                   # Database
│   └── mongo
│
├── internals                 
│   └── token_utils
│   └── utils
│
├── repository                #repository layer
│
├── services/                 # Business logic, separate from controllers
│
├── usecase/                  # business implement
│
├── main.go                  # business implement
│
├── .env                      # Environment variables
├── Dockerfile                      # Dockerfile
├── .gitignore
├── go.mod                    # Go module file
├── go.sum
└── README.md


## 📦 Installation & Run (Local)

```bash
# Clone the project
git clone https://github.com/iseT1enLoc/cakewai.git

# Enter the project directory
cd cakewai

# Run the application
go run main.go

## 📬 Contact

If you have questions, suggestions, or need support:

Nguyễn Võ Tiến Lộc
📧 Email: locnvt.it@gmail.com

Let me know if you’d like a `Dockerfile`, `.env.example`, or API documentation template included as well.


