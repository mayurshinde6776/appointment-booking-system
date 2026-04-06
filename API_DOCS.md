# Appointment Booking API Documentation

The REST API implements clean architecture and provides robust endpoints for Coaches to declare availability, and for Users to query dynamically calculated 30-minute time slots and securely make bookings.

Base URL: `http://localhost:8080/api/v1`

---

## 1. Set Coach Availability

Creates a recurring weekly shift bounding the hours a Coach works on a specific day of the week.

- **URL**: `/coaches/availability`
- **Method**: `POST`

### Request Body
```json
{
  "coach_id": 1,
  "day_of_week": "Monday",
  "start_time": "09:00",
  "end_time": "15:00"
}
```
*Note: `day_of_week` must match standard capitalized day names (e.g., `Monday`).*

### Success Response (201 Created)
```json
{
  "success": true,
  "message": "availability set successfully",
  "data": {
    "id": 1,
    "coach_id": 1,
    "day_of_week": "Monday",
    "start_time": "09:00",
    "end_time": "15:00"
  }
}
```

### Error Responses
- **400 Bad Request**: Invalid JSON or missing required fields.
- **500 Internal Server Error**: Database failures (e.g., passing a `coach_id` that does not exist in the database, resulting in a foreign key constraint violation).

---

## 2. Fetch Available Slots

Calculates an array of available 30-minute boundary ISO timestamps for a particular Coach on a specified exact Date. Automatically masks out times that already have an active Booking.

- **URL**: `/users/slots`
- **Method**: `GET`

### Query Parameters
| Parameter | Type | Required | Description |
|---|---|---|---|
| `coach_id` | `int` | Yes | The ID of the Coach offering slots. |
| `date` | `string` | Yes | Target query date formatted as `YYYY-MM-DD`. |

### Example Request
`GET /users/slots?coach_id=1&date=2025-10-27`

### Success Response (200 OK)
Returns a direct array of available ISO 8601 formatting strings:
```json
[
  "2025-10-27T09:00:00Z",
  "2025-10-27T09:30:00Z",
  "2025-10-27T10:00:00Z",
  "2025-10-27T14:30:00Z"
]
```

### Error Responses
- **400 Bad Request**: Missing/invalid `coach_id` or `date` formatting.

---

## 3. Create a Booking

Secures an appointment between a User and a Coach at a specific time. Defended by database-level transactions enforcing strict double-booking protections.

- **URL**: `/users/bookings`
- **Method**: `POST`

### Request Body
```json
{
  "user_id": 101,
  "coach_id": 1,
  "datetime": "2025-10-27T09:30:00Z"
}
```

### Success Response (201 Created)
```json
{
  "id": 1,
  "user_id": 101,
  "coach_id": 1,
  "slot_time": "2025-10-27T09:30:00Z",
  "status": "confirmed",
  "created_at": "2026-04-06T12:00:00Z"
}
```

### Error Responses
- **409 Conflict**: `{"error": "slot already booked"}` - Dual requests hit the same exact slot window and the unique DB constraint caught the collision.
- **400 Bad Request**: Mathematical validation failed (i.e. the chosen time does not correctly align to an openly available 30-min window for that coach).

---

## 4. Get User Bookings

Retrieves all active bookings that a specific User has scheduled.

- **URL**: `/users/bookings`
- **Method**: `GET`

### Query Parameters
| Parameter | Type | Required | Description |
|---|---|---|---|
| `user_id` | `int` | Yes | The ID of the querying User. |

### Example Request
`GET /users/bookings?user_id=101`

### Success Response (200 OK)
```json
[
  {
    "id": 1,
    "user_id": 101,
    "coach_id": 1,
    "slot_time": "2025-10-27T09:30:00Z",
    "status": "confirmed",
    "created_at": "2026-04-06T12:00:00Z"
  }
]
```

### Error Responses
- **400 Bad Request**: The `user_id` query parameter is missing or numerically invalid.
- **500 Internal Server Error**: Connectivity failures when retrieving records.

---

## 5. Cancel / Delete Booking

Releases a secured appointment slot, allowing it to become visually identifiable in `/users/slots` lists once again.

- **URL**: `/users/bookings/{id}`
- **Method**: `DELETE`

### Path Parameters
| Parameter | Type | Required | Description |
|---|---|---|---|
| `id` | `int` | Yes | The unique database ID of the specific Booking to delete. |

### Example Request
`DELETE /users/bookings/1`

### Success Response (200 OK)
```json
{
  "message": "booking deleted successfully"
}
```

### Error Responses
- **400 Bad Request**: Invalid Booking ID format.
- **404 Not Found**: A booking matching the supplied ID was not found inside the database.
