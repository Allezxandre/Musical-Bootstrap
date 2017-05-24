# Musical-Bootstrap

## Usage

### Register the App on Spotify

1. Go to [https://developer.spotify.com/my-applications/](https://developer.spotify.com/my-applications/) and create a new app.
2. In the `Redirect URI` field, add `http://localhost:19331/callback`, and save. This will whitelist the callback URL used by this app.
3. Create a `.env` file that looks like this:

```
SPOTIFY_ID=################################
SPOTIFY_SECRET=################################
```

using the `Client ID` and `Client Secret` values provided by Spotify.

### Run the app

In your terminal, run this:

```bash
./run.sh -artist "Tchami"
```

This script will fetch Go packages, build the project, source the variables in the `.env` file, and run.