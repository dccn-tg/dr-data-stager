var config = require('config');
var util = require('../lib/utility');
var express = require('express');
var router = express.Router();
const { createAdapter } = require("webdav-fs");

/* Authenticate user to the RDM service via the webdav interface */
router.post('/login', function(request, response, next) {

    var sess = request.session;

    var wfs = createAdapter(
        config.get('rdm.irodsWebDavEndpoint'), {
            "username": request.body.username,
            "password": request.body.password
        }
    );

    wfs.readdir('/.login', function(err) {
        if (!err) {
            wfs.readdir('/', function(err, contents) {
                if (!err) {
                    // set session data
                    if (typeof sess.user === "undefined" ||
                        typeof sess.pass === "undefined" ) {
                        sess.user = {rdm: request.body.username};
                        sess.pass = {rdm: request.body.password};
                    } else {
                        sess.user.rdm = request.body.username;
                        sess.pass.rdm = request.body.password;
                    }

                    response.status(200);
                    response.json(contents);
                } else {
                    console.error('login error: ' + err);
                    response.status(404);
                    response.json({'error': err.message});
                }
            });
        } else {
            console.error('login error: ' + err);
            response.status(404);
            response.json({'error': err.message});
        }
    });
});

/* logout user by removing corresponding session data */
router.post('/logout', function(request, response, next) {
    var sess = request.session;
    delete sess.user.rdm;
    delete sess.pass.rdm;
    response.json({'logout': true});
});

/* Get directory content for jsTree, using the WebDAV interface */
router.get('/dir', function(request, response, next) {

    var files = [];

    var sess = request.session;

    var dir = request.query.dir;
    var isRoot = request.query.isRoot;

    var wfs = createAdapter(
        config.get('rdm.irodsWebDavEndpoint'), {
            "username": sess.user['rdm'],
            "password": sess.pass['rdm']
        }
    );

    wfs.readdir(dir, "stat", function(err, contents) {
        if (!err) {
            contents.forEach( function(f) {
                if ( f.isFile() ) {
                    files.push({
                        id: dir.replace(/\/$/,'') + '/' + f.name,
                        type: 'f',
                        parent: isRoot === 'true' ? '#':dir,
                        text: f.name,
                        icon: 'fa fa-file-o',
                        li_attr: {'title':''+f.size+' bytes'},
                        children: false
                    });
                } else {
                    files.push({
                        id: dir.replace(/\/$/,'') + '/' + f.name + '/',
                        type: 'd',
                        parent: isRoot === 'true' ? '#':dir,
                        text: f.name,
                        icon: 'fa fa-folder',
                        li_attr: {},
                        children: true
                    });
                }
            });
            response.json(files);
        } else {
            console.log("Error:", err.message);
            util.responseOnError('json',[],response);
        }
    }, 'stat');
});

/* create new directory via WebDAV */
router.post('/mkdir', function(request, response, next) {
    var sess = request.session;

    var dir = request.body.dir;

    var wfs = createAdapter(
        config.get('rdm.irodsWebDavEndpoint'), {
            "username": sess.user['rdm'],
            "password": sess.pass['rdm']
        }
    );

    wfs.mkdir(dir, function(err) {
        if (!err) {
            response.status(200);
            response.json(['OK']);
        } else {
            console.log("WebDAV error:", err.message);
            util.responseOnError('json',[err.message],response);
        }
    });
});

module.exports = router;