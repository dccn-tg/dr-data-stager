var express = require('express');
var session = require('express-session');
var path = require('path');
var favicon = require('serve-favicon');
var logger = require('morgan');
var cookieParser = require('cookie-parser');
var bodyParser = require('body-parser');

var passport = require('passport');

var authRouter = require('./routes/oidc');
var appRoutes = require('./routes/index');
var apiRoutesFS = require('./routes/mod_fs');
var apiRoutesRdm = require('./routes/mod_rdm');
var apiRoutesStager = require('./routes/mod_stager');

var app = express();

// view engine setup
app.set('views', path.join(__dirname, 'views'));
app.set('view engine', 'jade');

app.use(favicon(path.join(__dirname, 'public', 'favicon.ico')));
app.use(logger('dev'));
app.use(bodyParser.json({limit:'50mb'}));
app.use(bodyParser.urlencoded({limit: '50mb', extended: false}));
app.use(cookieParser());
app.use(express.static(path.join(__dirname, 'public')));

/* session property
   - rolling expiration upon access
   - save newly initiated session right into the store
   - delete session from story when unset
   - cookie age: 4 hours (w/ rolling expiration)
   - session data store: memory on the server
*/
app.use( session({
    secret: 'planet Tatooine',
    resave: true,
    rolling: true,
    saveUninitialized: true,
    unset: 'destroy',
    name: 'stager-ui.sid',
    cookie: {
        httpOnly: false,
        maxAge: 4 * 3600 * 1000
    }
}));


app.use(cookieParser());

app.use(passport.initialize());
app.use(passport.session());

// OIDC auth router
app.use('/oidc', authRouter);

// AJAX posts
app.use('/fs', apiRoutesFS)
app.use('/rdm', apiRoutesRdm)
app.use('/stager', apiRoutesStager)

// main webapp page
app.use('/', appRoutes);

// error handlers
app.use(function(err, _, res, _) {
  if (err.status == 401) {
    // redirect client to the login endpoint
    res.redirect('/oidc/login');
  } else {
    res.status(err.status || 500);
    res.render('error', {
      message: err.message,
      error: app.get('env') === 'development' ? err : {}
    });
  }
});


module.exports = app;
