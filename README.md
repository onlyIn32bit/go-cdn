# Go CDN

This is a lightweight file upload and static file server built with Go and the Fiber web framework. It allows you to upload files securely using an API key and serves them via a public URL.

## Features

- Secure file uploads with API key authentication.
- Automatically generates unique filenames for uploaded files or using provided file name.
- Serves uploaded files as static assets.
- Configurable via environment variables and flags.

## Requirements

- Go 1.24 or later
- `.env` file with the required configuration (see below).

## Installation

1. Install the release

2. Create the `.env` with your secret based on the `.env.example` file

    ```env
    API_KEY="your-api-key" # API key set for interacting with the server
    DOMAIN="https://cdn.your.domain.com" # Your domain
    ```

3. Run the executable

    ```bash
    ./go-cdn
    ```

    The server will run by default on port `8090` and serve content from `./uploads` folder.

    **Output**

    ```txt
    Started server on port: 8090
    ```

    You can specify the port and uploads folder by using flags

    - `--port`: `./go-cdn --port=3001`
    - `--upload-dir`: `./go-cdn --upload-dir="./content"`

## Run from source

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
    DOMAIN="https://cdn.your.domain.com"
    ```

4. Run it!

    ```bash
    go run main.go
    ```
