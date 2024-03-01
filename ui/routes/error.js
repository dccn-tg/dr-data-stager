var express = require('express');
var router = express.Router();

/** simple router for error code page */
router.get('/403', function(_, _, next) {
    return next({
        status: 403
    });
});

module.exports = router;