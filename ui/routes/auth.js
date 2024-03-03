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

    console.log("user:", JSON.stringify(req.user));

    // no user in request context, must run login
    if (!req.user) return next(createError(401, new Error("Unauthorized")));

    if (req.user && req.user.validUntil > (Date.now()/1000)) {
        next();
    } else {

        console.log("valid unti:", req.user.validUntil);

        refresh.requestNewAccessToken(
            'oidc',
            req.user.refreshToken,
            (err, accessToken, refreshToken) => {
                if (err) {
                    console.error("refresh error:", err);
                    return res.status(401).send({
                        message: 'unauthorized'
                    });
                }

                // update user's access token and refresh token stored in request context
                req.user.accessToken = accessToken;
                req.user.refreshToken = refreshToken;
                next();
            }
        );
    }
}

module.exports.isAuthenticated = _isAuthenticated;