# Waterlogger API Documentation

This document describes the REST API endpoints available in Waterlogger.

## Base URL

All API endpoints are relative to the base URL: `http://localhost:2341/api`

## Authentication

Most endpoints require authentication. The application uses session-based authentication with cookies.

### Login

```http
POST /api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login successful"
}
```

### Logout

```http
POST /api/logout
```

**Response:**
```json
{
  "message": "Logout successful"
}
```

## Users

### List Users

```http
GET /api/users
```

**Response:**
```json
[
  {
    "id": 1,
    "username": "admin",
    "email": "admin@example.com",
    "created_at": "2024-07-14T10:30:00Z",
    "updated_at": "2024-07-14T10:30:00Z"
  }
]
```

### Create User

```http
POST /api/users
Content-Type: application/json

{
  "username": "newuser",
  "email": "newuser@example.com",
  "password": "SecurePassword123!"
}
```

### Update User

```http
PUT /api/users/{id}
Content-Type: application/json

{
  "username": "updateduser",
  "email": "updated@example.com"
}
```

### Delete User

```http
DELETE /api/users/{id}
```

## Pools

### List Pools

```http
GET /api/pools
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Main Pool",
    "volume_gallons": 20000,
    "type": "pool",
    "system_description": "Salt water system with heater",
    "created_at": "2024-07-14T10:30:00Z",
    "updated_at": "2024-07-14T10:30:00Z"
  }
]
```

### Create Pool

```http
POST /api/pools
Content-Type: application/json

{
  "name": "Hot Tub",
  "volume_gallons": 400,
  "type": "hot_tub",
  "system_description": "Heated spa with jets"
}
```

### Update Pool

```http
PUT /api/pools/{id}
Content-Type: application/json

{
  "name": "Updated Pool Name",
  "volume_gallons": 22000,
  "type": "pool",
  "system_description": "Updated system description"
}
```

### Delete Pool

```http
DELETE /api/pools/{id}
```

## Test Kits

### List Kits

```http
GET /api/kits
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Taylor K-2006",
    "description": "Complete pool test kit",
    "purchased_date": "2024-01-15T00:00:00Z",
    "replenished_date": "2024-07-01T00:00:00Z",
    "created_at": "2024-07-14T10:30:00Z",
    "updated_at": "2024-07-14T10:30:00Z"
  }
]
```

### Create Kit

```http
POST /api/kits
Content-Type: application/json

{
  "name": "Test Kit Name",
  "description": "Kit description",
  "purchased_date": "2024-01-15T00:00:00Z",
  "replenished_date": "2024-07-01T00:00:00Z"
}
```

### Update Kit

```http
PUT /api/kits/{id}
Content-Type: application/json

{
  "name": "Updated Kit Name",
  "description": "Updated description"
}
```

### Delete Kit

```http
DELETE /api/kits/{id}
```

## Samples

### List Samples

```http
GET /api/samples
```

**Response:**
```json
[
  {
    "id": 1,
    "pool_id": 1,
    "sample_datetime": "2024-07-14T14:30:00Z",
    "user_id": 1,
    "kit_id": 1,
    "pool": {
      "id": 1,
      "name": "Main Pool"
    },
    "user": {
      "id": 1,
      "username": "admin"
    },
    "kit": {
      "id": 1,
      "name": "Taylor K-2006"
    },
    "measurements": {
      "id": 1,
      "sample_id": 1,
      "fc": 2.5,
      "tc": 2.6,
      "ph": 7.4,
      "ta": 100,
      "ch": 250,
      "temperature": 78.5,
      "created_at": "2024-07-14T14:30:00Z"
    },
    "indices": {
      "id": 1,
      "sample_id": 1,
      "lsi": -0.1,
      "rsi": 7.2,
      "comment": null,
      "created_at": "2024-07-14T14:30:00Z"
    },
    "created_at": "2024-07-14T14:30:00Z",
    "updated_at": "2024-07-14T14:30:00Z"
  }
]
```

### Create Sample

```http
POST /api/samples
Content-Type: application/json

{
  "pool_id": 1,
  "sample_datetime": "2024-07-14T14:30:00Z",
  "user_id": 1,
  "kit_id": 1,
  "measurements": {
    "fc": 2.5,
    "tc": 2.6,
    "ph": 7.4,
    "ta": 100,
    "ch": 250,
    "cya": 35,
    "temperature": 78.5,
    "salinity": 3200,
    "tds": 1500,
    "appearance": "Clear and blue",
    "maintenance": "Added chlorine"
  }
}
```

### Update Sample

```http
PUT /api/samples/{id}
Content-Type: application/json

{
  "sample_datetime": "2024-07-14T15:00:00Z",
  "measurements": {
    "fc": 3.0,
    "tc": 3.1,
    "ph": 7.5,
    "ta": 105,
    "ch": 250,
    "temperature": 79.0
  }
}
```

### Delete Sample

```http
DELETE /api/samples/{id}
```

## Charts

### Get Chart Data

```http
GET /api/charts/data?pool_id=1&start_date=2024-07-01&end_date=2024-07-14&parameters=ph,fc,tc
```

**Query Parameters:**
- `pool_id` (optional): Filter by pool ID
- `start_date` (optional): Start date in ISO format
- `end_date` (optional): End date in ISO format
- `parameters` (optional): Comma-separated list of parameters to include

**Response:**
```json
{
  "labels": ["2024-07-01", "2024-07-02", "2024-07-03"],
  "datasets": [
    {
      "label": "pH",
      "data": [7.4, 7.5, 7.3],
      "borderColor": "#3b82f6",
      "backgroundColor": "rgba(59, 130, 246, 0.1)"
    },
    {
      "label": "Free Chlorine",
      "data": [2.5, 2.8, 2.3],
      "borderColor": "#10b981",
      "backgroundColor": "rgba(16, 185, 129, 0.1)"
    }
  ]
}
```

## Export

### Export to Excel

```http
GET /api/export/excel
```

**Response:**
- Content-Type: `application/vnd.openxmlformats-officedocument.spreadsheetml.sheet`
- Content-Disposition: `attachment; filename="WL20240714_143022.xlsx"`

### Export to Markdown

```http
GET /api/export/markdown
```

**Response:**
- Content-Type: `text/markdown`
- Content-Disposition: `attachment; filename="WL20240714_143022.md"`

## Settings

### Get Settings

```http
GET /api/settings
```

**Response:**
```json
{
  "unit_system": "imperial",
  "default_pool_id": 1,
  "chart_parameters": ["ph", "fc", "tc", "ta", "ch", "temperature"],
  "export_format": "excel"
}
```

### Update Settings

```http
POST /api/settings
Content-Type: application/json

{
  "unit_system": "metric",
  "default_pool_id": 2,
  "chart_parameters": ["ph", "fc", "tc"],
  "export_format": "markdown"
}
```

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request

```json
{
  "error": "Validation error message"
}
```

### 401 Unauthorized

```json
{
  "error": "Authentication required"
}
```

### 403 Forbidden

```json
{
  "error": "Insufficient permissions"
}
```

### 404 Not Found

```json
{
  "error": "Resource not found"
}
```

### 500 Internal Server Error

```json
{
  "error": "Internal server error"
}
```

## Water Chemistry Parameters

### Parameter Definitions

- **FC (Free Chlorine)**: Available chlorine for sanitization (ppm)
- **TC (Total Chlorine)**: Free chlorine + combined chlorine (ppm)
- **pH**: Acidity/alkalinity level (0-14 scale)
- **TA (Total Alkalinity)**: pH buffering capacity (ppm)
- **CH (Calcium Hardness)**: Dissolved calcium concentration (ppm)
- **CYA (Cyanuric Acid)**: Chlorine stabilizer (ppm)
- **Temperature**: Water temperature (°F or °C)
- **Salinity**: Salt content for saltwater pools (ppm)
- **TDS (Total Dissolved Solids)**: Total dissolved substances (mg/L)

### Ideal Ranges

- **FC**: 1.0 - 4.0 ppm
- **TC**: Should match FC (minimize combined chlorine)
- **pH**: 7.4 - 7.6
- **TA**: 80 - 120 ppm
- **CH**: 200 - 400 ppm
- **CYA**: 30 - 50 ppm
- **Temperature**: 78 - 82°F (25 - 28°C)
- **Salinity**: 2,700 - 3,400 ppm (optimal: 3,200 ppm)
- **TDS**: < 1,500 ppm

### Calculated Indices

- **LSI (Langelier Saturation Index)**: -0.3 to +0.3 (balanced water)
- **RSI (Ryznar Stability Index)**: 6.0 - 7.0 (stable water)

## Rate Limiting

The API implements basic rate limiting to prevent abuse:

- **General endpoints**: 100 requests per minute per IP
- **Authentication endpoints**: 10 requests per minute per IP
- **Export endpoints**: 5 requests per minute per IP

## Data Formats

### Date/Time Format

All date/time values use ISO 8601 format: `2024-07-14T14:30:00Z`

### Numeric Values

- Floating-point numbers for measurements
- Integers for IDs and counts
- Null values for optional parameters

### String Values

- UTF-8 encoding
- Maximum lengths enforced on input
- HTML content is escaped for security