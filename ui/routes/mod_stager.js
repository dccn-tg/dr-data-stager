var config = require('config');
var path = require('path');
var RestClient = require('node-rest-client').Client;
var util = require('../lib/utility');
var express = require('express');
var router = express.Router();

const auth = require("./auth");

/* Get directory content for jsTree */
router.get('/dir', auth.isAuthenticated, function(request, response) {

    var files = [];
    var dir = request.query.dir;
    var isRoot = request.query.isRoot;

    var args = {
        data: { path: dir },
        headers: {
            "Accept": "application/json",
            "Content-Type": "application/json"
        }
    };

    var c = new RestClient({
        user: request.user.username,
        password: request.user.token
    });

    var apiURL = config.get('stager.restfulEndpoint') + '/dir'

    c.get(apiURL, args, function(data, resp) {
        try {
            if ( resp.statusCode == 200 ) {
                data.entries.forEach( function(f) {
                    switch (f.type) {
                        case 'dir':
                            files.push({
                                id: path.join(dir, f.name) + '/',
                                parent: isRoot === 'true' ? '#':dir,
                                text: f.name,
                                icon: 'fa fa-folder',
                                li_attr: {},
                                children: true
                            });
                            break;
                        case 'regular':
                            files.push({
                                id: path.join(dir, f.name),
                                parent: isRoot === 'true' ? '#':dir,
                                text: f.name,
                                icon: 'fa fa-file-o',
                                li_attr: {'title':''+f.size+' bytes'},
                                children: false
                            });
                            break;
                        case 'symlink':

                            var names = f.name.split(' -> ');

                            files.push({
                                id: path.join(dir, names[0]),
                                parent: isRoot === 'true' ? '#':dir,
                                text: names[0],
                                icon: 'fa fa-file-o',
                                li_attr: {'title': names.length === 2 ? 'symlink to ' + names[1] : 'symlink'},
                                state: {
                                    disabled: true
                                },
                                children: false
                            });
                            break;
                        default:
                            console.warn('ignore file with unknown type: ', path.join(dir, f.name));
                    }
                });
                response.json(files);
            } else {
                console.error('GET', apiURL, 'args: ', JSON.stringify(args), 'status:', resp.statusCode);
                util.responseOnError('json', [], response);
            }
        } catch(e) {
            console.error(e);
            util.responseOnError('json', [], response);
        }
    }).on('error', function(e) {
        console.error(e);
        util.responseOnError('json', [], response);
    });
});

/* Start or restart a stopped job */
router.put('/job/:id/state/inactive', function(request, response) {
    var sess = request.session;
    var c = new RestClient({user: sess.user.stager,
                            password: sess.pass.stager});                        
    // get the job to check if the job is owned by the stager users
    try {

    } catch(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    }    
});

/* Get single transfer job by id */
router.get('/job/:id', auth.isAuthenticated, function(request, response) {

    var args = { headers: { "Accept": "application/json" } };
    var c = new RestClient({
        user: request.user.username,
        password: request.user.token
    });

    var url = config.get('stager.restfulEndpoint') + '/job/' + request.params.id;

    c.get(url, args, function(data, resp) {
        if ( resp.statusCode == 200 ) {
            response.status(200);
            response.json(data);
        } else {
            console.log('api-server response status: ' + resp.statusCode);
            response.status(resp.statusCode);
            response.json({
                message: resp.statusMessage
            });
        }
    }).on('error', function(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    });
});

/* Submit transfer jobs to stager */
router.post('/jobs', auth.isAuthenticated, function(request, response) {

    var sess = request.session;

    var jobs = [];
    if ( typeof request.body.jobs !== 'undefined' ) {
        jobs = JSON.parse(request.body.jobs);
    }

    var stagerJobs = jobs.map( j => {
        return {
            ...j,
            stagerUser: request.user.username,
            stagerUserEmail: request.user.email,
            drUser: sess.user.rdm,
            drPass: util.encryptStringWithRsaPublicKey(sess.pass.rdm, '/opt/stager-ui/ssl/public.pem'),
            timeout: 86400,
            timeout_noprogress: 3600,
            title: 'sync to ' + j.dstURL
        };
    });

    if ( stagerJobs.length > 0 ) {

        var c = new RestClient({
            user: request.user.username,
            password: request.user.token
        });

        var args = {
            headers: {
                "Accept": "application/json",
                "Content-Type": "application/json"
            },
            data: {
                jobs: stagerJobs
            }
        };

        var url = config.get('stager.restfulEndpoint') + '/jobs';
        c.post(url, args, function(data, resp) {
            if ( resp.statusCode == 200 ) {
                response.status(200);
                response.json(data.jobs);
            } else {
                console.error('api-server response status: ' + resp.statusCode);
                response.status(resp.statusCode);
                response.json([]);
            }
        }).on('error', function(e) {
            console.error(e);
            util.responseOnError('json', [], response);
        });
    } else {
        console.log('No stager job to submit');
        response.status(200);
        response.json([]);
    }
});

/* Get all transfer jobs from stager and show only those belongs to the same user */
router.get('/jobs', auth.isAuthenticated, function(request, response) {

    var args = { headers: { "Accept": "application/json" } };
    var c = new RestClient({
        user: request.user.username,
        password: request.user.token
    });

    var url = config.get('stager.restfulEndpoint') + '/jobs';
    c.get(url, args, function(data, resp) {
        if ( resp.statusCode == 200 ) {
            response.status(200);
            response.json(data.jobs);
        } else {
            console.error('api-server response status: ' + resp.statusCode);
            response.status(resp.statusCode);
            response.json([]);
        }
    }).on('error', function(e) {
        console.error(e);
        util.responseOnError('json', [], response);
    });
});

/* delete an existing job */
router.delete('/job/:id', auth.isAuthenticated, function(request, response) {

    var args = { headers: { "Accept": "application/json" } };
    var c = new RestClient({
        user: request.user.username,
        password: request.user.token
    });

    var url = config.get('stager.restfulEndpoint') + '/job/' + request.params.id;

    c.delete(url, args, function(data, resp) {
        if ( resp.statusCode == 200 ) {
            response.status(200);
            response.json(data);
        } else {
            console.log('api-server response status: ' + resp.statusCode);
            response.status(resp.statusCode);
            response.json({
                message: resp.statusMessage
            });
        }
    }).on('error', function(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    });
});

module.exports = router;