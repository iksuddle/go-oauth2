# go oauth2 github example

## setup

Clone the repo.

To run the example you'll need to [create an OAuth App in GitHub](https://github.com/settings/developers).
Use the following URLs:
- Homepage URL: `http://localhost:3000`
- Authorization Callback URL: `http://localhost:3000/login/callback`

Once you get your client ID and client secret, create an .env file:

```
PORT=3000
CLIENT_ID=your_client_id
CLIENT_SECRET=your_client_secret
SESION_KEY=session_key
```

`PORT` should be the same as the one used in your OAuth App URLs
and the `SESSION_KEY` is used by gorilla/sessions for session authentication.
In an actual server make sure it is randomized.

## usage

Run the project using
```
make run
```
and visit `http://localhost:3000` in your browser.
