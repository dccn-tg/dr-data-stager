var config = require('config');
var path = require('path');
var RestClient = require('node-rest-client').Client;
var util = require('../lib/utility');
var express = require('express');
var router = express.Router();

const auth = require("./auth");

/* retrieve job detail by the given id */
var _getJobDetail = function(id, user, pass, cb) {
    var url = config.get('stager.restfulEndpoint') + '/job/' + id;
    var c = new RestClient({user: user, password: pass});
    var args = { headers: { "Accept": "application/json" } };
    c.get(url, args, function(j, resp) {
        console.log('stager response status: ' + resp.statusCode);
        if ( resp.statusCode == 200 ) {
            if ((typeof j.data !== 'undefined') &&
                (typeof j.data.stagerUser !== 'undefined') &&
                (j.data.stagerUser == user)) {
                cb(j, '');
            } else {
                cb(null, "job not found or user doesn't own the job: " + id)
            }
        } else {
            cb(null, "fail retrieving job detail: " + id + " code: " + resp.statusCode);
        }
    }).on('error', function(e) {
        cb(null, e);
    });
}

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

/* Get transfer-job counts from Stager */
router.get('/job/state', function(request, response) {

    var args = { headers: { "Accept": "application/json" } };
    var sess = request.session;
    var c = new RestClient({user: sess.user.stager,
                            password: sess.pass.stager});

    c.get(config.get('stager.restfulEndpoint') + '/stats', args, function(data, resp) {

        try {
            console.log('stager response status: ' + resp.statusCode);
            if ( resp.statusCode == 200 ) {
                response.status(200);
                response.json(data);
            } else {
                response.status(404);
                response.json({});
            }
        } catch(e) {
            console.error(e);
            util.responseOnError('json', {}, response);
        }
    }).on('error', function(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    });
});

/* Get transfer jobs from stager and show only those belongs to the same user */
router.get('/jobs/:state/:from-:to', function(request, response) {

    var args = { headers: { "Accept": "application/json" } };
    var sess = request.session;
    var c = new RestClient({user: sess.user.stager,
                            password: sess.pass.stager});

    var jobs = {};
    var state  = request.params.state;
    var idx_f = request.params.from;
    var idx_t = request.params.to;
    var url = config.get('stager.restfulEndpoint') + '/jobs/' + state + '/' + idx_f + '..' + idx_t + '/desc';
    c.get(url, args, function(data, resp) {
        try {
            console.log('stager response status: ' + resp.statusCode);
            if ( resp.statusCode == 200 ) {
                response.status(200);
                jobs = data.filter( function(j) {
                    return (typeof j.data !== 'undefined') &&
                           (typeof j.data.stagerUser !== 'undefined') &&
                           (j.data.stagerUser == sess.user.stager);
                });
                response.json(jobs);
            } else {
                response.status(404);
                response.json({});
            }
        } catch(e) {
            console.error(e);
            util.responseOnError('json', {}, response);
        }
    }).on('error', function(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    });
});

/* Get single transfer job by id */
router.get('/job/:id', function(request, response) {
    var sess = request.session;
    try {
        _getJobDetail(request.params.id, sess.user.stager, sess.pass.stager, function(job, error) {
            if (error) {
                console.error(error);
                util.responseOnError('json', {}, response);
            }
            response.status(200);
            response.json(job);
        });
    } catch(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    }
});

/* delete an existing job */
router.delete('job/:id', function(request, response) {
    var sess = request.session;
    var c = new RestClient({user: sess.user.stager,
                            password: sess.pass.stager});                        
    // get the job to check if the job is owned by the stager users
    try {
        _getJobDetail(request.params.id, sess.user.stager, sess.pass.stager, function(job, error) {
            if (error) {
                console.error(e);
                util.responseOnError('json', {}, response);
            }
            var url = config.get('stager.restfulEndpoint') + '/job/' + request.params.id;
            var args = { headers: { "Accept": "application/json" } };
            var req = c.delete(url, args, function(msg, resp) {
                console.log('stager response status: ' + resp.statusCode);
                if ( resp.statusCode == 200 ) {
                    response.json(msg);
                } else {
                    response.status(404);
                    response.json({});
                }            
            });
        });
    } catch(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    }
});

/* Start or restart a stopped job */
router.put('/job/:id/state/inactive', function(request, response) {
    var sess = request.session;
    var c = new RestClient({user: sess.user.stager,
                            password: sess.pass.stager});                        
    // get the job to check if the job is owned by the stager users
    try {
        _getJobDetail(request.params.id, sess.user.stager, sess.pass.stager, function(job, error) {
            if (error) {
                console.error(e);
                util.responseOnError('json', {}, response);
            }
            var url = config.get('stager.restfulEndpoint') + '/job/' + request.params.id + '/state/inactive';
            var args = { headers: { "Accept": "application/json" } };
            var req = c.put(url, args, function(msg, resp) {
                console.log('stager response status: ' + resp.statusCode);
                if ( resp.statusCode == 200 ) {
                    response.json(msg);
                } else {
                    response.status(404);
                    response.json({});
                }            
            });
        });
    } catch(e) {
        console.error(e);
        util.responseOnError('json', {}, response);
    }    
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
            try {
                console.log('stager service response status: ' + resp.statusCode);
                if ( resp.statusCode == 200 ) {
                    response.status(200);
                    response.json(data.jobs);
                } else {
                    response.status(404);
                    response.json([]);
                }
            } catch(e) {
                console.error(e);
                util.responseOnError('json', [], response);
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
        try {
            console.log('stager response status: ' + resp.statusCode);
            if ( resp.statusCode == 200 ) {
                response.status(200);
                response.json(data.jobs);
            } else {
                response.status(404);
                response.json([]);
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

module.exports = router;