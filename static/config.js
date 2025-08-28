const functions = require("firebase-functions");

exports.getConfig = functions.https.onRequest((req, res) => {
  const BACKEND_URL = functions.config().backend.url;  // get backend url from firebase config
  res.json({ BACKEND_URL });
});
