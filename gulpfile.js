var gulp = require('gulp');
var concat = require('gulp-concat');
var sass = require('gulp-sass');
var browserSync = require('browser-sync').create();
var del = require('del');

var paths = {
    scripts: {
        src: [
            'bower_components/angular/angular.min.js',
            'bower_components/angular-resource/angular-resource.min.js',
            'bower_components/jquery/dist/jquery.min.js',
            'bower_components/uikit/js/uikit.min.js',
            'bower_components/uikit/js/core/**/*.min.js',
            'app/scripts/**/app.js',
            'app/scripts/**/*.js'
        ],
        concat: 'all.js',
        dest: 'public/js'
    },
    styles: {
        src: [
            'app/styles/main.scss'
        ],
        concat: 'all.css',
        dest: 'public/css'
    },
    views: [
        'templates/**/*'
    ],
    fonts: {
        src: [
            'bower_components/uikit/fonts/**/*'
        ],
        dest: 'public/fonts'
    }
};

gulp.task('fonts', function () {
    gulp.src(paths.fonts.src)
        .pipe(gulp.dest(paths.fonts.dest));
});

gulp.task('views', function () {
    gulp.src(paths.views);
});

gulp.task('styles', function () {
    gulp.src(paths.styles.src)
        .pipe(sass())
        .pipe(concat(paths.styles.concat))
        .pipe(gulp.dest(paths.styles.dest));
});

gulp.task('scripts', function () {
    return gulp.src(paths.scripts.src)
        .pipe(concat(paths.scripts.concat))
        .pipe(gulp.dest(paths.scripts.dest));
});

gulp.task('watch', function () {

    browserSync.init({
        proxy: "localhost:9000",
        logPrefix: "BrowserSync",
        logConnections: false,
        reloadOnRestart: false,
        notify: false,
        open: false,
        tunnel: true
    });

    gulp.watch(paths.scripts, ['scripts', 'styles', 'views']).on("change", browserSync.reload);
    gulp.watch(paths.styles, ['scripts', 'styles', 'views']).on("change", browserSync.reload);
    gulp.watch(paths.views, ['scripts', 'styles', 'views']).on("change", browserSync.reload);
});

// The default task (called when you run `gulp` from cli)
gulp.task('default', ['scripts', 'styles', 'fonts']);
