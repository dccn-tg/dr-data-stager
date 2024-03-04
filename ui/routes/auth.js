const createError = require("http-errors");
var refresh = require('passport-oauth2-refresh');
/**
 * Middleware to verify session authentication status
 *
 * When using the OIDC, we expect a valid `req.user` in
 * an authenticated session.
 * 
 */
const _isAuthenticated = function(req, res, next) {

    // no user in request context, must run login
    if (!req.user) return next(createError(401, new Error("Unauthorized")));

    if (req.user && isTokenValid(req.user.token)) {
        next();
    } else {
        // try to refresh the access token
        refresh.requestNewAccessToken(
            'oidc',
            req.user.refresh_token,
            (err, accessToken, refreshToken) => {
                if (err) {
                    console.error("cannot refresh token:", err);
                    return res.status(401).send({
                        message: 'unauthorized'
                    });
                }

                // update user's access token and refresh token stored in request context
                req.session.passport.user.token = accessToken;
                req.session.passport.user.refresh_token = refreshToken;

                req.user.token = accessToken;
                req.user.refresh_token = refreshToken;

                next();
            }
        );
    }
}

// isTokenValid checks if the claims.exp in the jwt token has a valid lifetime
// longer than 30 seconds.
const isTokenValid = function(token) {
    try {
        // extract access token's expieration time
        const payload = Buffer.from(token.split('.')[1], 'base64').toString();
        const claims = JSON.parse(payload);
        return claims.exp > (Date.now()/1000) + 30;
    } catch(e) {
        return false;
    }
}

module.exports.isAuthenticated = _isAuthenticated;