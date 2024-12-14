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

- Install [Golang](https://golang.org/dl/)
- MySQL [MYSQL](https://dev.mysql.com/downloads/installer/)
- Redis Install [Redis](https://redis.io/download)
- Docker Install [Docker](https://docs.docker.com/engine/install/)

---

### 🐳 Docker Deployment

1. **Pull Docker Image**:

   ```bash
   docker pull ikhsanhilmi/multifinance-app-service:latest
   ```

2. **Run Container**:

   ```bash
   docker run -d \
     --name multifinance-service \
     -p 8080:8080 \
     -e MYSQL_URL=mysql://user:password@host:port/database \
     -e REDIS_URL=redis://host:port \
     -e JWT_SECRET=your_secret_key \
     -e APP_ENV=production \
     ikhsanhilmi/multifinance-app-service:latest
   ```

### Environment Variables

- **MYSQL_URL**: MySQL connection string
- **REDIS_URL**: Redis connection URL
- **JWT_SECRET**: Secret key for JWT token generation
- **APP_ENV**: Application environment (development/production)

---

### 🔧 Local Development Setup

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/multifinance-service.git
   cd multifinance-service
   ```

2. **Install Dependencies**:

   ```bash
   go mod tidy
   ```

3. **Configure Environment**:

   - Copy `.env.example` to `.env`
   - Fill in the required configuration values

4. **Database Migrations**:

   ```bash
   Manage database migrations using `goose`:

- Create a new migration:
    ```bash
    make goose-create name=create_users_table
    ```
- Apply migrations:
    ```bash
    make goose-up
    ```
- Rollback migrations:
    ```bash
    make goose-down
    ```
- Check migration status:
    ```bash
    make goose-status
    ```
    
For more details on `goose`, visit the [official documentation](https://github.com/pressly/goose).
   ```

5. **Running the Application**:

   ```bash
   # Run in development mode
   make run
   ```

---

## Development

### Developer Guide

#### General Rules

- We use conventional commits: <https://www.conventionalcommits.org>
  - `feat: commit message` - For new features.
  - `refactor: commit message` - For code refactoring.
  - `fix: commit message` - For bug fixes.
  - `test: commit message` - For tests.
  - `docs: commit message` - For documentation updates.
  - `style: commit message` - For code style changes.

- Use `git-chglog` to generate changelogs before merging to the release branch:  
  <https://github.com/git-chglog/git-chglog>

#### Branching Strategy

1. Use **feature branches** for all new features and bug fixes.
2. Merge **feature branches** into the main branch using pull requests.
3. Keep the main branch up-to-date and high quality.

**Feature Branch Naming Convention:**

- `feature/<feature-name>`: For new features.
- `bugfix/<bug-description>`: For bug fixes.
- `hotfix/<fix-description>`: For urgent fixes.

#### Folder Structure

- `cmd/bin`: Contains `main.go`, which runs the API server or seeds the database.
- `internal`:
  - `adapter`: Holds driving and driven adapters:
    - **Driving Adapters**: Interfaces for the API handler (e.g., REST, CLI).
    - **Driven Adapters**: Interfaces for database/repository interactions.
  - `infrastructure`: Configuration and logger setup.
  - `module`: Contains the core business logic and entities.
    - `entity`: Defines data models (DTOs/DAOs).
    - `repository`: Handles database operations.
    - `service`: Implements business logic.
    - `handler`: Manages API request handling.
  - `route`: Stores route definitions.
  - `pkg`: Shared utilities or common functions.
  - `logs`: Stores application log files.

---

## How to Create a New Module

1. Copy the `internal/module/z_template` folder and rename it to the new module name.
2. Update the following subfolders within the module:
   - `entity`: Define the data models (DTO/DAO).
   - `ports`: Define interfaces for the module.
   - `repository`: Implement database interactions.
   - `service`: Implement business logic.
   - `handler`: Implement API handlers.
3. Update route definitions in the `route` folder.

---

### Development Workflow

1. Read the requirements and define entities in the `entity` folder.
2. Define contracts (interfaces) in the `ports` folder.
3. Implement repository and service layers.
4. Create and integrate API handlers.
5. Test the feature thoroughly.

---

### 🧪 Testing

1. **Run Unit Tests**:

   ```bash
   make test
   ```

---

## 🔒 Security Features

### OWASP Top 10 Protection Mechanisms

- **SQL Injection Prevention**
- **Cross-Site Scripting (XSS) Protection**
- **JWT Token Validation**
- **Input Validation**
- **Rate Limiting**
- **Secure Configuration Management**

---

## 📦 Database Schema

### Customers Table

- **id**: Primary Key
- **nik**: Unique National ID
- **full_name**: Customer Full Name
- **legal_name**: Name on Legal Documents
- **birth_place**: Place of Birth
- **birth_date**: Date of Birth
- **salary**: Customer's Monthly Income
- **ktp_photo_path**: KTP (ID Card) Photo Path
- **selfie_photo_path**: Selfie Photo Path

### Transactions Table

- **id**: Primary Key
- **customer_id**: Foreign Key (Customers)
- **contract_number**: Unique Transaction Contract Number
- **on_the_road_price**: Total Asset Price
- **admin_fee**: Administrative Transaction Fee
- **installment_amount**: Monthly Installment Amount
- **interest_amount**: Total Interest Amount
- **asset_name**: Purchased Asset Name
- **tenor_months**: Loan Tenure in Months
- **created_at**: Transaction Creation Timestamp

---

## 🏗 Architectural Highlights

- **Clean Code Architecture**
- **Modular Design**
- **Dependency Injection**
- **Repository Pattern**
- **Service Layer Abstraction**
- **Concurrent Transaction Handling**

---

## 🤝 Contributing

1. **Fork the Repository**
2. **Create a Feature Branch**:

   ```bash
   git checkout -b feature/AmazingFeature
   ```

3. **Commit Changes**:

   ```bash
   git commit -m 'Add some AmazingFeature'
   ```

4. **Push to Branch**:

   ```bash
   git push origin feature/AmazingFeature
   ```

5. **Open a Pull Request**

### Contribution Guidelines

- Follow Go coding standards
- Write unit tests for new features
- Update documentation
- Ensure CI/CD checks pass

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

---

## 📞 Contact

- **Hilmi Ikhsan**
  - **Email**: [your-email@example.com]
  - **LinkedIn**: [Your LinkedIn Profile]
  - **Project Link**: [https://github.com/yourusername/multifinance-service](https://github.com/yourusername/multifinance-service)

---

## 🙏 Acknowledgements

- Golang Community
- Fiber Framework
- MySQL
- Redis
- GoMock for Testing
