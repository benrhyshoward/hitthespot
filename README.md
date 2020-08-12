# Hit The Spot

Personalised music quiz based on your favourite music from Spotify.

See it deployed at [hitthespot.app](https://hitthespot.app)

## Tech
 - Go backend
 - React + Redux frontend
 - MongoDB
 - Docker
 - Docker Compose
 - Spotify API for user listening history
 - Musixmatch API for lyrics

## Running locally
Steps
- Get some Spotify API credentials from a [Spotify developer application](https://developer.spotify.com/dashboard/applications)
- Add `http://localhost/auth/callback` as a redirect uri in the Spotify application dashboard
- Get a valid [Musixmatch API key](https://developer.musixmatch.com/) 
- Fill these credentials into `.env`
- Run `docker-compose up`
- Go to `http://localhost:3000/`

## Running in a prod environment  
 A bit more complicated than locally, the main differences are
 - We don't want to use the React dev server, so require bundling front end code and serving everything through Go
 - We don't want to serve over http, so require a certificate
 - We don't want to use a temporary Mongo container, so require a connection to a running Mongo instance

Steps
- Generate a certificate for the host you are deploying to
    - Place the cert into `/server/certs/certificate.pem`
    - Place the key into `/server/certs/key.pem` 
- Go to `/client` and run `yarn` then `yarn build`
- Copy the `build` folder into `/server/build`
- Inside `/server` run `docker build -t hitthespot .`
- Run `docker run -p 443:8080 --{{environment-variables}} hitthespot`
    With the following environment variables
    - `HTTPS=true`
    - `STATIC_FILE_PATH=/go/src/github.com/benrhyshoward/hitthespot/server/build`
    - `SPOTIFY_ID=<spotify api id>`
    - `SPOTIFY_SECRET=<spotify secret>`
    - `MUSIXMATCH_API_KEY=<musixmatch api key>`
    - `MONGO_CONNECTION_STRING=<mongo db connection string>`
    - `GO_SERVER_EXTERNAL_URL=https://<hostname>`
    - `FRONTEND_SERVER_EXTERNAL_URL=https://<hostname>`
    where `<hostname>` is the hostname of the server you are deploying to
- Go to `https://<hostname>`

## Future improvements

- Add tests
- More types of questions
- Timers for questions
- Variable scores based on how quickly answers are given
- Error handling for front end fetch requests
- Maybe use shorter ids to make urls a bit nicer
- Could use browser routing rather than hash routing
- Remove ability to directly browse static files and directories through Go server
- Caching when doing Spotify API calls, currently each question type fetches the same information separately
- Use paging to get more than 50 results from Spotify API calls
- Improve Mongo DB schema for faster quering and add indexes
