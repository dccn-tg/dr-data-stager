const path = require('path');
const fs = require('fs');
const crypto = require('crypto');

// general error handler to send response to the client
var _responseOnError = function(c_type, c_data, resp) {
    resp.status(500);
    if (c_type === 'json') {
        resp.json(c_data);
    } else {
        resp.send(c_data);
    }
}

var _encryptStringWithRsaPublicKey = function(toEncrypt, relativeOrAbsolutePathToPublicKey) {
    var absolutePath = path.resolve(relativeOrAbsolutePathToPublicKey);
    var publicKey = fs.readFileSync(absolutePath, "utf-8");
    var buffer = Buffer.from(toEncrypt);
    
    // encrypt data using the public key, the padding is needed for the service (Go code) to
    // encrypt it correctly.
    var encrypted = crypto.publicEncrypt({
        key: publicKey,
		padding: crypto.constants.RSA_PKCS1_PADDING
    }, buffer);
    return encrypted.toString("base64");
};

module.exports.responseOnError = _responseOnError;
module.exports.encryptStringWithRsaPublicKey = _encryptStringWithRsaPublicKey;
