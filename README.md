# 🏨 Hotel Reservation System

<div align="center">

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-4EA94B?style=for-the-badge&logo=mongodb&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=JSON%20web%20tokens&logoColor=white)

*A modern, scalable hotel reservation system built with Go, Fiber, and MongoDB*

[📚 Documentation](#documentation) • [🚀 Quick Start](#quick-start) • [📋 API Reference](#api-reference) • [🛠 Development](#development)

</div>

---

## ✨ Features

<table>
<tr>
<td width="50%">

### 🔐 **Authentication & Security**
- JWT-based authentication
- Role-based access control
- Rate limiting protection
- CORS security headers
- Request validation & sanitization

### 🏨 **Hotel Management**
- Complete CRUD operations
- Room inventory management
- Real-time availability tracking
- Multi-location support

</td>
<td width="50%">

### 📅 **Smart Booking System**
- Intelligent availability checking
- Automated email notifications
- Booking status management
- Conflict resolution

### 👥 **User Experience**
- Intuitive registration flow
- Profile management
- Booking history tracking
- Admin dashboard capabilities

</td>
</tr>
</table>

---

## 🏗 Architecture Overview

```mermaid
graph TB
    Client[👤 Client] --> API[🌐 Fiber API Server]
    API --> Auth[🔐 JWT Middleware]
    Auth --> Routes[📍 Route Handlers]
    
    Routes --> UserSvc[👥 User Service]
    Routes --> HotelSvc[🏨 Hotel Service]
    Routes --> BookingSvc[📅 Booking Service]
    
    UserSvc --> UserDB[(👤 Users Collection)]
    HotelSvc --> HotelDB[(🏨 Hotels Collection)]
    BookingSvc --> BookingDB[(📅 Bookings Collection)]
    
    UserDB --> MongoDB[(🍃 MongoDB)]
    HotelDB --> MongoDB
    BookingDB --> MongoDB
    
    BookingSvc --> Email[📧 Email Service]
    
    style API fill:#e1f5fe
    style MongoDB fill:#4caf50,color:#fff
    style Email fill:#ff9800,color:#fff
```

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.19+ 
- **MongoDB** 4.4+
- **Git**

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/hotel-reservation-system.git
cd hotel-reservation-system

# Install dependencies
go mod tidy

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Run the application
go run main.go
```

### 🐳 Docker Quick Start

```bash
# Using Docker Compose
docker-compose up -d

# The API will be available at http://localhost:8080
```

---

## 📁 Project Structure

```
hotel-reservation-system/
├── 📂 api/                 # HTTP request handlers
│   ├── auth_handler.go     # Authentication endpoints
│   ├── booking_handler.go  # Booking management
│   ├── hotel_handler.go    # Hotel operations
│   └── user_handler.go     # User management
├── 📂 cmd/                 # CLI applications
├── 📂 config/              # Configuration utilities
├── 📂 db/                  # Database layer
│   ├── models/             # Data models
│   └── stores/             # Database operations
├── 📂 middleware/          # HTTP middleware
│   ├── auth.go             # JWT authentication
│   ├── cors.go             # CORS handling
│   └── ratelimit.go        # Rate limiting
├── 📂 services/            # Business logic
├── 📂 tests/               # Test suites
├── 📂 types/               # Type definitions
├── 📂 util/                # Utility functions
└── 📄 main.go              # Application entry point
```

---

## 🛣 API Routes

### 🔓 Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/auth` | User authentication |
| `POST` | `/api/register` | User registration |

### 🔒 Protected Endpoints

<details>
<summary><strong>🏨 Hotel Operations</strong></summary>

| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| `GET` | `/api/v1/hotel` | List all hotels | User |
| `GET` | `/api/v1/hotel/:id` | Get hotel details | User |
| `POST` | `/api/v1/hotel` | Create hotel | Admin |
| `PUT` | `/api/v1/hotel/:id` | Update hotel | Admin |
| `DELETE` | `/api/v1/hotel/:id` | Delete hotel | Admin |

</details>

<details>
<summary><strong>🛏 Room Operations</strong></summary>

| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| `GET` | `/api/v1/room` | List available rooms | User |
| `GET` | `/api/v1/room/:id` | Get room details | User |
| `POST` | `/api/v1/room` | Create room | Admin |
| `PUT` | `/api/v1/room/:id` | Update room | Admin |
| `DELETE` | `/api/v1/room/:id` | Delete room | Admin |

</details>

<details>
<summary><strong>📅 Booking Operations</strong></summary>

| Method | Endpoint | Description | Role |
|--------|----------|-------------|------|
| `GET` | `/api/v1/booking` | Get user bookings | User |
| `POST` | `/api/v1/booking` | Create booking | User |
| `GET` | `/api/v1/booking/:id` | Get booking details | User |
| `PUT` | `/api/v1/booking/:id` | Update booking | User |
| `DELETE` | `/api/v1/booking/:id` | Cancel booking | User |

</details>

### 👑 Admin Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/v1/admin/users` | Manage all users |
| `GET` | `/api/v1/admin/bookings` | View all bookings |
| `GET` | `/api/v1/admin/hotels` | Manage all hotels |

---

## 🔄 Booking Flow

```mermaid
sequenceDiagram
    participant U as User
    participant API as API Server
    participant DB as MongoDB
    participant Email as Email Service
    
    U->>API: Search available rooms
    API->>DB: Query room availability
    DB-->>API: Return available rooms
    API-->>U: Display available rooms
    
    U->>API: Create booking request
    API->>DB: Validate availability
    DB-->>API: Confirm availability
    API->>DB: Create booking (Pending)
    DB-->>API: Booking created
    
    API->>Email: Send confirmation email
    Email-->>U: Booking confirmation
    API-->>U: Booking response
    
    Note over API,DB: Admin can approve/reject
    API->>DB: Update booking status
    API->>Email: Send status update
    Email-->>U: Status notification
```

---

## 🛠 Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test suite
go test ./tests/api_test.go
```

### Database Seeding

```bash
# Seed the database with sample data
go run cmd/seed/main.go
```

### Environment Variables

```env
# Server Configuration
PORT=8080
JWT_SECRET=your-super-secret-key

# Database Configuration
MONGO_URI=mongodb://localhost:27017/hotel_reservation
MONGO_DB_NAME=hotel_reservation

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
```

---

## 🔐 Authentication Flow

```mermaid
graph LR
    A[User Login] --> B[Validate Credentials]
    B --> C{Valid?}
    C -->|Yes| D[Generate JWT]
    C -->|No| E[Return Error]
    D --> F[Return Token]
    F --> G[Include in Headers]
    G --> H[Access Protected Routes]
    
    style D fill:#4caf50,color:#fff
    style E fill:#f44336,color:#fff
    style H fill:#2196f3,color:#fff
```

---

## 📊 Data Models

### User Model
```go
type User struct {
    ID       primitive.ObjectID `bson:"_id" json:"id"`
    Email    string            `bson:"email" json:"email"`
    Password string            `bson:"password" json:"-"`
    Role     string            `bson:"role" json:"role"`
    Profile  UserProfile       `bson:"profile" json:"profile"`
}
```

### Hotel Model
```go
type Hotel struct {
    ID          primitive.ObjectID `bson:"_id" json:"id"`
    Name        string            `bson:"name" json:"name"`
    Location    string            `bson:"location" json:"location"`
    Rating      float64           `bson:"rating" json:"rating"`
    Rooms       []Room            `bson:"rooms" json:"rooms"`
    Amenities   []string          `bson:"amenities" json:"amenities"`
}
```

### Booking Model
```go
type Booking struct {
    ID          primitive.ObjectID `bson:"_id" json:"id"`
    UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
    HotelID     primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
    RoomID      primitive.ObjectID `bson:"room_id" json:"room_id"`
    CheckIn     time.Time         `bson:"check_in" json:"check_in"`
    CheckOut    time.Time         `bson:"check_out" json:"check_out"`
    Status      BookingStatus     `bson:"status" json:"status"`
    TotalAmount float64           `bson:"total_amount" json:"total_amount"`
}
```

---

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 Support

<div align="center">

**Need help?** 

[![GitHub Issues](https://img.shields.io/badge/GitHub-Issues-red?style=for-the-badge&logo=github)](https://github.com/yourusername/hotel-reservation-system/issues)
[![Discord](https://img.shields.io/badge/Discord-7289DA?style=for-the-badge&logo=discord&logoColor=white)](https://discord.gg/your-discord)
[![Email](https://img.shields.io/badge/Email-D14836?style=for-the-badge&logo=gmail&logoColor=white)](mailto:support@yourhotel.com)

</div>

---

<div align="center">

</div>
