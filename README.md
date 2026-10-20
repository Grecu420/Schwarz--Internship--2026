# 🏡 StayNest

**StayNest** is a scalable full-stack booking platform built with a **Vue 3** frontend and a **Go-based microservices** architecture. It provides a seamless, intuitive experience for both hosts managing properties and guests booking their next stay, powered by real-time messaging, interactive maps, and asynchronous background processing.

---

## ✨ Features & Walkthrough

### 👤 1. Account & Profile Management
> 🎥 

https://github.com/user-attachments/assets/5e4ff89e-43ae-4b88-8d30-78fa90690595



- 🔒 **Secure Authentication:** Register and log in securely using JWT and Bcrypt password hashing.
- ⚙️ **Profile Customization:** Update personal details, manage account settings, and switch seamlessly between guest and host interactions.

### 🗺️ 2. Discover & Connect
> 🎥 

https://github.com/user-attachments/assets/606c98f6-0488-4d19-8dba-aaa5e8d1892a





- 📍 **Interactive Map Search:** Browse available locations via an integrated **Leaflet** map with **Nominatim** reverse-geocoding for precise address lookups.
- 💬 **Responsive Messaging:** Communicate directly with hosts or guests using a smooth chat module powered by HTTP polling and Pinia state management.

### ✈️ 3. Booking Management *(Guest Experience)*
> 🎥 












https://github.com/user-attachments/assets/4c05df2f-8c0b-4035-8456-a8b651fa6833



- 📅 **Streamlined Reservations:** Select dates and book properties through a frictionless, step-by-step flow.
- 🧳 **Trip Dashboard:** Track upcoming and past reservations in a dedicated view.
- 📧 **Automated Notifications:** Receive asynchronous email confirmations handled smoothly in the background by **RabbitMQ**.

### 🏠 4. Property Management *(Host Experience)*
> 🎥 

https://github.com/user-attachments/assets/fa8ebbf6-14b8-4523-b998-4b7751035fff







- 📝 **List a Property:** Create new listings easily with automated address completion.
- 🖼️ **Optimized Media:** Upload property photos with automatic client-side optimization and cloud storage integration via **Cloudinary**.
- 🎛️ **Listing Control:** Edit property details, update availability, and manage listing statuses directly from the host dashboard.

---

## 🛠️ Tech Stack

- 💻 **Frontend:** Vue 3, TypeScript, Pinia, Axios, Sit-Onyx, Leaflet
- ⚙️ **Backend:** Go (Microservices), gRPC & Protobuf, PostgreSQL, RabbitMQ
- 🐳 **Infrastructure:** Docker
