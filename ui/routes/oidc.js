var express = require('express');

const auth = require("./auth");

// openid-client
var Issuer = require('openid-client').Issuer;

// setup passport-openidconnect
var passport = require('passport');
var refresh = require('passport-oauth2-refresh');
var OidcStrategy = require('passport-openidconnect').Strategy;

var authServer = process.env.AUTH_SERVER;

var oidcStrategy = new OidcStrategy({
    issuer: authServer,
    authorizationURL: authServer + '/connect/authorize',
    tokenURL: authServer + '/connect/token',
    userInfoURL: authServer + '/connect/userinfo',
    clientID: process.env.AUTH_CLIENT_ID,
    clientSecret: process.env.AUTH_CLIENT_SECRET,
    callbackURL: '/oidc/callback',
    skipUserProfile: true,  // we are going to get user profile ourselves.
    proxy: true,
    scope: ["openid", "profile", "offline_access", "urn:dccn:identity:uid", "urn:dccn:data-stager-api:*"],
}, (_issuer, _profile, _context, idToken, accessToken, refreshToken, verified) => {

    // getting user profile, and check if the profile contains 'urn:dccn:uid' attribute.
    Issuer.discover(authServer).then((issuer) => {
        new issuer.Client({
            client_id: process.env.AUTH_CLIENT_ID,
        }).userinfo(accessToken).then(profile => {
            // only user with DCCN account is authorized??
            if ( ! profile['urn:dccn:uid'] ) throw new Error("missing DCCN account in profile");
            return verified(null, {
                id_token: idToken,
                token: accessToken,
                refresh_token: refreshToken,
                username: profile['urn:dccn:uid'],
                displayName: profile.name,
                email: profile.email
            });
        }).catch(err => {
            console.error(err);
            return verified(null, false);
        });
    });
})

passport.use('oidc', oidcStrategy);
refresh.use('oidc', oidcStrategy);

// serialize functions are needed to store user object (returned from the `verified` callback) in session.
passport.serializeUser(function(user, cb) {
    cb(null, user);
});
  
passport.deserializeUser(function(user, cb) {
    return cb(null, user);
});

/**
 * Reconstructs the original URL of the request.
 * 
 * This code is inspired by https://github.com/jaredhanson/passport-openidconnect/blob/fee0639a75235e8cce4597d6a87c9f1bcb3cdb8e/lib/utils.js#L17
 */
const logoutRedirectUrl = function(req) {
    const host = req.headers['x-forwarded-host'] || req.get('host');
    const tls  = req.connection.encrypted || ('https' == (req.headers['x-forwarded-proto'] || "").toLowerCase().split(/\s*,\s*/)[0]);
    return (tls ? 'https':'http') + "://" + host;
};

var router = express.Router();

// endpoint to trigger OIDC login workflow
router.get('/login', passport.authenticate('oidc'));

// callback endpoint after authentication at OIDC provider
router.get('/callback', passport.authenticate('oidc', {
    successReturnToOrRedirect: '/',
    failureRedirect: '/error/403'
}));

// endpoint to trigger logout workflow
// TODO: should better use DELETE or POST method? (CORS issue)
router.get('/logout',
    //auth.isAuthenticated,
    (req, res) => {
        if (req.user && req.user.id_token) {
            const id_token_hint = req.user.id_token;
            req.logout(err => {
                if (err) {
                    console.log("service logout error: ", err);
                }
                // redirect browser to the end_session_endpoint of the OIDC provider
                res.redirect(authServer + 
                    "/connect/endsession?id_token_hint=" + id_token_hint + 
                    "&post_logout_redirect_uri=" + logoutRedirectUrl(req));
            });
        } else {
            // redirect browser to the end_session_endpoint without id_token_hint
            res.redirect(authServer + 
                "/connect/endsession?post_logout_redirect_uri=" + logoutRedirectUrl(req));
        }
});

// endpoint go get user profile fetched from OIDC provider and stored in the session.
router.get("/profile",
    auth.isAuthenticated,
    (req, res) => {
        res.status(200).json({
            data: {
                id: req.user.username,
                name: req.user.displayName
            },
            error: null
        });
    }
);

module.exports = router;