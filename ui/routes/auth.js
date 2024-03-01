const createError = require("http-errors");

/**
 * Middleware to verify session authentication status
 *
 * When using the OIDC, we expect a valid `req.user` in
 * an authenticated session.
 * 
 */
const _isAuthenticated = function(req, _, next) {
    if (req.user && req.user.validUntil > (Date.now()/1000)) {
        next();
    } else {
        next(createError(401, new Error("Unauthorized")));
    }
}

module.exports.isAuthenticated = _isAuthenticated;