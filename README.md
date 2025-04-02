# Hotel Reservation System

A complete hotel booking system built with Go, featuring user authentication, hotel and room management, and MongoDB integration.

## What is this project?

This is a web-based system that lets users:
- Create an account and log in
- Browse available hotels and rooms
- Make hotel reservations
- Manage their bookings

For hotel managers, it provides tools to:
- Add and update hotel information
- Manage room availability
- View reservations

## What you need to get started

1. **Go** version 1.16 or newer ([Download here](https://golang.org/dl/))
2. **MongoDB** running on your computer ([Download here](https://www.mongodb.com/try/download/community))
3. Basic knowledge of using the command line

## How to set up the project

1. **Get the code**

   ```bash
   git clone https://github.com/0x0Glitch/hotel-reservation.git
   cd hotel-reservation
   ```

2. **Set your secret key** 

   This is used for securing user accounts. Open a terminal and type:

   ```bash
   # On Mac/Linux
   export JWT_SECRET=your_secret_key_here

   # On Windows
   set JWT_SECRET=your_secret_key_here
   ```

3. **Put some sample data in your database**

   ```bash
   make seed
   ```

   This adds sample hotels, rooms, and a test user to get you started.

4. **Run the app**

   ```bash
   go run main.go
   ```

   Your server is now running at http://localhost:5001

## How to use the system

### Creating an account

1. Send a POST request to `/api/v1/user` with this data:
   ```json
   {
     "firstName": "Your",
     "lastName": "Name",
     "email": "your.email@example.com",
     "password": "your_password"
   }
   ```

### Logging in

1. Send a POST request to `/api/auth` with:
   ```json
   {
     "email": "your.email@example.com",
     "password": "your_password"
   }
   ```

2. You'll get back a token. Save this token and include it in all future requests in the `X-Api-Token` header.

### Finding hotels

1. Send a GET request to `/api/v1/hotel` with your token in the `X-Api-Token` header
2. You'll get a list of all available hotels

### Looking at rooms in a hotel

1. First, find the hotel ID from the previous step
2. Send a GET request to `/api/v1/hotel/{hotelID}/rooms` 
3. You'll see all rooms available for that hotel

## How to run tests

We have comprehensive tests for all parts of the system.

To run all tests:

```bash
go test ./tests/...
```

To run tests for a specific part (like just the user functionality):

```bash
go test ./tests/types -v
go test ./tests/db -v
go test ./tests/api -v
go test ./tests/middleware -v
```

## Project structure explained

Here's what each folder contains:

- **api/** - Handles web requests and responses
- **db/** - Connects to the database and manages data
- **types/** - Defines the structure of users, hotels, and rooms
- **middleware/** - Handles things like authentication
- **tests/** - Makes sure everything works correctly
- **scripts/** - Helper scripts, like adding sample data

## Common problems and solutions

1. **Can't connect to MongoDB**
   - Make sure MongoDB is running
   - Check if the connection string is correct in `db/db.go`

2. **Authentication not working**
   - Make sure your JWT_SECRET environment variable is set
   - Check that you're including the token in the X-Api-Token header

3. **Getting "unauthorized" errors**
   - Your token might have expired - try logging in again
   - Check that you're using the correct token format

## Want to improve this project?

We welcome contributions! Here's how:

1. Fork the repository
2. Create a feature branch: `git checkout -b new-feature`
3. Make your changes
4. Write tests for your changes
5. Run existing tests: `go test ./tests/...`
6. Submit a pull request

## Contact and support

If you have questions or need help, open an issue on GitHub or contact the project maintainer.

## License

This project is available under the MIT License - see the LICENSE file for details. 