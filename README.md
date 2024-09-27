# TAB - Tabo

As Australians, pubs have been a big part of our lives, and with pubs comes Keno flashing across screens, eating up some of our hard-earned drinking money. Despite the losses, Keno has always brought smiles to our faces.

To explore WebSockets, Vue, and multiplayer game development, we decided to create our own version of Keno, whimsically named **Tabo** - because, why not?

## Get Started

You can run the Tabo project in two ways: by using Docker or by running the services directly.

### Docker Method
To run Tabo with Docker, follow these steps:

1. **Build the Docker images**: This step creates the Docker images for both the backend and frontend services.

```bash
docker compose build
```

2. **Start the services**: Once the images are built, this command spins up the necessary containers.

```bash
docker compose up
```

 - The frontend and backend will be run in separate containers.
 - The application should be accessible after both services start.

> **Note**: If you need to stop the services, simply press CTRL+C or run docker compose down to stop and remove the containers.

### Running Services Directly

If you'd rather not use Docker, you can run the frontend and backend manually. Here's how:

1. **Backend**: In one terminal window (or a Tmux pane), navigate to the backend folder and run:

```bash
go run main.go
```

**Frontend**: In another terminal window (or a separate Tmux pane), navigate to the frontend folder and run:

```bash
npm run dev
```
> This will start the development server for the frontend.

## Environment Variables

When running the project, you will need the following environment variables for the frontend configuration:

```bash
VITE_CLIENT_ID=${DiscordAppClientID} # Your Discord app client ID
VITE_REDIRECT_URI=http://localhost:8081/redirect # Redirect URI for OAuth
VITE_TABO_BACKEND_WEBSOCKET=ws://localhost:8080/api/v1/ws # Backend WebSocket URL
```

These are used to manage OAuth for login and establish WebSocket connections with the backend.

## Database

The backend uses SQLite as the database engine. When the backend starts, it will generate a `keno.db` file in the working directory.

- **Docker**: If you're using Docker, the database is stored in a local volume, ensuring persistence between container restarts.
- **Manual Setup**: If running manually, you may want to mount or back up the `keno.db` file to ensure data is saved.

## Live Demo

The project is currently live and running at [tabo.tabdiscord.com](https://tabo.tabdiscord.com/).