var express = require('express');
var router = express.Router();
var StatusCodes = require('http-status-codes').StatusCodes;
var ReasonPhrases = require('http-status-codes').ReasonPhrases;

/** simple router for error code page */
router.get('/403', function(_, _, next) {
    return next({
        status: StatusCodes.FORBIDDEN,
        message: ReasonPhrases.FORBIDDEN
    });
});

module.exports = router;