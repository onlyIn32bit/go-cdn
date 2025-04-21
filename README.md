# Go CDN

This is a lightweight file upload and static file server built with Go and the Fiber web framework. It allows you to upload files securely using an API key and serves them via a public URL.

## Features

- Secure file uploads with API key authentication.
- Automatically generates unique filenames for uploaded files.
- Serves uploaded files as static assets.
- Configurable via environment variables.

## Requirements

- Go 1.24 or later
- A `.env` file with the required configuration (see below).

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/onlyin32bit/go-cdn.git
   cd go-cdn
   ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Create a `.env` file based on `.env.example` and change the configuration:

    ```env
    API_KEY="your-api-key"
    UPLOAD_DIR="./uploads"
    DOMAIN="https://cdn.your.domain.com"

    ```
