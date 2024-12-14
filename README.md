# Multifinance Service

A robust and scalable microfinance management system built with **Golang**, **Fiber**, **MySQL**, **Redis**, and modern software architecture practices.

---

## 🚀 Features

- **Transaction Management**: Seamless handling of financial transactions
- **Customer Limit Tracking**: Manage customer borrowing limits
- **Secure Authentication**: JWT-based user authentication
- **High Performance**: Built with Fiber web framework
- **Scalable Architecture**: Modular design following Clean Code principles
- **Concurrent Transaction Handling**: Advanced concurrency management

---

## 🛠 Technologies

- **Programming Language**: Golang 1.22
- **Web Framework**: Fiber
- **Database**: MySQL
- **Caching**: Redis
- **Authentication**: JWT
- **Migration**: Goose
- **Testing**: Go's built-in testing with coverage

---

## 📋 System Overview

Multifinance Service is designed to support:
- Flexible loan limits (1, 2, 3, 6 months)
- Transactions across e-commerce, web, and partner dealers
- Comprehensive customer and transaction tracking

---

## 📊 Entity-Relationship Diagram (ERD)

![Multifinance Service ERD](./docs/ERD.png)

---

## 🏁 Getting Started

### Prerequisites

- Golang 1.22+
- MySQL
- Redis
- Docker (optional)

---

### 🐳 Docker Deployment

1. **Pull Docker Image**:
```bash
docker pull ikhsanhilmi/multifinance-app-service:latest

Run Container:

bashCopydocker run -d \
  --name multifinance-service \
  -p 8080:8080 \
  -e MYSQL_URL=mysql://user:password@host:port/database \
  -e REDIS_URL=redis://host:port \
  -e JWT_SECRET=your_secret_key \
  -e APP_ENV=production \
  ikhsanhilmi/multifinance-app-service:latest

Environment Variables:

MYSQL_URL: MySQL connection string
REDIS_URL: Redis connection URL
JWT_SECRET: Secret key for JWT token generation
APP_ENV: Application environment (development/production)




🔧 Local Development Setup
Clone the Repository
bashCopygit clone https://github.com/yourusername/multifinance-service.git
cd multifinance-service
Install Dependencies
bashCopygo mod tidy
Configure Environment

Copy .env.example to .env
Fill in the required configuration values

Database Migrations
bashCopy# Create new migration
make migrate-create name=create_transactions_table

# Apply migrations
make migrate-up

# Rollback migrations
make migrate-down
Running the Application
bashCopy# Run in development mode
make dev

# Build for production
make build

# Run production binary
./bin/multifinance-service

🧪 Testing
bashCopy# Run unit tests
make test

# Run tests with coverage
make test-coverage

# Run specific package tests
make test package=./internal/module/transaction

🔒 Security Features

OWASP Top 10 Protection Mechanisms

SQL Injection Prevention
Cross-Site Scripting (XSS) Protection
JWT Token Validation


Input Validation
Rate Limiting
Secure Configuration Management


📦 Database Schema
Customers Table

id: Primary Key
nik: Unique National ID
full_name: Customer Full Name
legal_name: Name on Legal Documents
birth_place: Place of Birth
birth_date: Date of Birth
salary: Customer's Monthly Income
ktp_photo_path: KTP (ID Card) Photo Path
selfie_photo_path: Selfie Photo Path

Transactions Table

id: Primary Key
customer_id: Foreign Key (Customers)
contract_number: Unique Transaction Contract Number
on_the_road_price: Total Asset Price
admin_fee: Administrative Transaction Fee
installment_amount: Monthly Installment Amount
interest_amount: Total Interest Amount
asset_name: Purchased Asset Name
tenor_months: Loan Tenure in Months
created_at: Transaction Creation Timestamp


🏗 Architectural Highlights

Clean Code Architecture
Modular Design
Dependency Injection
Repository Pattern
Service Layer Abstraction
Concurrent Transaction Handling


🤝 Contributing

Fork the repository
Create a feature branch (git checkout -b feature/AmazingFeature)
Commit changes (git commit -m 'Add some AmazingFeature')
Push to branch (git push origin feature/AmazingFeature)
Open a Pull Request

Contribution Guidelines

Follow Go coding standards
Write unit tests for new features
Update documentation
Ensure CI/CD checks pass


📄 License
Distributed under the MIT License. See LICENSE for more information.

📞 Contact
Hilmi Ikhsan

Email: [your-email@example.com]
LinkedIn: [Your LinkedIn Profile]
Project Link: https://github.com/yourusername/multifinance-service


🙏 Acknowledgements

Golang Community
Fiber Framework
MySQL
Redis
GoMock for Testing