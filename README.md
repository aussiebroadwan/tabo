# TAB - Tabo

As everyone in our group are Australians, we spent most our life in pubs with
the entertainment of Keno playing on every other screen. Lots of our drinking
money has converted to thin air thanks to Keno but it does put a smile on our
face.

So as a study with Websockets, Vue and making multiplayer games, we have decided
to make our own implementation of Keno which we call Tabo for no aparent reason.

## Get Started

This project has both a frontend and backend. Eventually there will be a docker
container which can be build and deployed to handle this, but at the moment you
will have to run the backend and frontend independantly.

```bash
# Tmux:0 Run the Backend
go run main.go

# Tmux:1 Run the Frontend
npm run dev 
```

> **Note:** This is currently running live at [tabo](https://tabo.tabdiscord.com)
