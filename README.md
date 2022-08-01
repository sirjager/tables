
# Tables

[![Test](https://github.com/SirJager/tables/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/SirJager/tables/actions/workflows/test.yml)

Tables allows you to manage you database tables in a more easier way .

- [@SirJager](https://www.github.com/SirJager)

# Endpoint

## UnAuthenticated Endpoints

| Title         | Descrption     | Method   |  Endpoint             |
|---------------|----------------|----------|-----------------------|
| Signup        | Create account | POST     | `/users/signup`       |
| Signin        | Login  account | POST     | `/users/signin`       |
| Refresh Tokne | Refresh Token  | POST     | `/users/renew-access` |


## Authenticated Endpoints

### Users Routes

| Title         | Descrption     | Method   |  Endpoint       |
|---------------|----------------|----------|-----------------|
| List Users    | List All Users | GET      | `/users`        |
| Current user  | Get  User      | GET      | `/users/me`     |
| Remove User   | Delete my acc  | DELETE   | `/users/me`     |


### Tables Routes

| Title         | Descrption     | Method   |  Endpoint                |
|---------------|----------------|----------|--------------------------|
| Create Table  | Create   Table | POST     | `/tables`                |
| List My Tables| Tables by me   | GET      | `/tables`                |
| Table Schema  | Table columns  | GET      | `/tables/{tablename}`    |
| Remove Table  | Delete my table| DELETE   | `/tables/{tablename}`    |

### Tables Routes

| Title         | Descrption     | Method   |  Endpoint                |
|---------------|----------------|----------|--------------------------|
| Create Table  | Create   Table | POST     | `/tables`                |
| List My Tables| Tables by me   | GET      | `/tables`                |
| Table Schema  | Table columns  | GET      | `/tables/{tablename}`    |
| Remove Table  | Delete my table| DELETE   | `/tables/{tablename}`    |

### Manage Columns Routes

| Title              | Descrption          | Method   |  Endpoint                         |
|--------------------|---------------------|----------|-----------------------------------|
| Add Column         | Add Column          | POST     | `/tables/{tablename}/columns`     |
| Remove Column      | Delete Column       | DELETE   | `/tables/{tablename}/columns`     |
| Add Primary key    | Add Primary Key     | POST     | `to be implemented`               |
| Delete Primary key | Delete Primary Key  | DELETE   | `to be implemented`               |
| Update Primary key | Update Primary Key  | PATCH    | `to be implemented`               |


### Manage Rows

| Title         | Descrption     | Method   |  Endpoint                      |
|---------------|----------------|----------|--------------------------------|
| Get Rows      | Get rows       | GET      | `/tables/{tablename}/rows`     |
| Insert Rows   | Insert rows    | POST     | `/tables/{tablename}/rows`     |
| Delete Rows   | Delete Rows    | DELETE   | `/tables/{tablename}/rows`     |
| Update Row    | Update Rows    | PATCH    | `to be implemented` 


### 🚀 What Next

- Make tests
- Add, Remove Update Primary Keys
- Update Rows
- Build a front end


## Tech Stack

**Language:** Go

**Database:** Postgres

**Web framework:** Gin

**Authentication:** JWT, Paseto
