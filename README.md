# bajid-spotify-server

# Configuration

The application is built to run within GCP Cloud Run and is backed by Firestore.

It needs several parameters:
- SPOTIFY_CLIENT_ID is passed at built time in the Cloud Run
- SPOTIFY_CLIENT_SECRET is fetched from Secrets Manager
