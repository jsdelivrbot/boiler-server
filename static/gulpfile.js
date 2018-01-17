'use strict';

// sass compile
var gulp = require('gulp');
var minifyCss = require("gulp-clean-css");
var rename = require("gulp-rename");
var uglify = require("gulp-uglify");
var concat = require("gulp-concat");



//*** CSS  minify task
gulp.task('minify', function () {
    // css minify 
    gulp.src(['./assets/apps/css/*.css', '!./assets/apps/css/*.min.css']).pipe(minifyCss()).pipe(rename({suffix: '.min'})).pipe(gulp.dest('./assets/apps/css/'));

    gulp.src(['./assets/global/css/*.css','!./assets/global/css/*.min.css']).pipe(minifyCss()).pipe(rename({suffix: '.min'})).pipe(gulp.dest('./assets/global/css/'));
    gulp.src(['./assets/pages/css/*.css','!./assets/pages/css/*.min.css']).pipe(minifyCss()).pipe(rename({suffix: '.min'})).pipe(gulp.dest('./assets/pages/css/'));    
    
    gulp.src(['./assets/layouts/**/css/*.css','!./assets/layouts/**/css/*.min.css']).pipe(rename({suffix: '.min'})).pipe(minifyCss()).pipe(gulp.dest('./assets/layouts/'));
    gulp.src(['./assets/layouts/**/css/**/*.css','!./assets/layouts/**/css/**/*.min.css']).pipe(rename({suffix: '.min'})).pipe(minifyCss()).pipe(gulp.dest('./assets/layouts/'));

    gulp.src(['./assets/global/plugins/bootstrap/css/*.css','!./assets/global/plugins/bootstrap/css/*.min.css']).pipe(minifyCss()).pipe(rename({suffix: '.min'})).pipe(gulp.dest('./assets/global/plugins/bootstrap/css/'));

   
});

//合并压缩css
gulp.task('minCss',function(){
    gulp.src(['assets/global/plugins/datatables/datatables.css','assets/global/plugins/datatables/plugins/bootstrap/datatables.bootstrap.css',
        'assets/global/plugins/bootstrap/css/bootstrap.min.css',
        'assets/global/plugins/bootstrap-switch/css/bootstrap-switch.min.css','assets/global/plugins/bootstrap-daterangepicker/daterangepicker.min.css',
        'assets/global/plugins/fullcalendar/fullcalendar.min.css',
        'bower_components/sweetalert2/dist/sweetalert2.css',
        'bower_components/angular-ui-select/dist/select.css',
        '/js/scripts/ladda/ladda-themeless.min.css',
        'assets/global/css/components-rounded.min.css',
        'assets/global/css/plugins.min.css',
        'assets/layouts/boiler/css/layout.min.css',
        'assets/layouts/boiler/css/themes/darkblue.min.css'
        ])
        .pipe(concat('base.min.css'))
        .pipe(minifyCss())
        .pipe(gulp.dest('assets'));
})



//合并压缩js插件
gulp.task('jsBase', function () {
    gulp.src(['./assets/global/plugins/jquery.min.js','./assets/global/plugins/bootstrap/js/bootstrap.min.js',
        './assets/global/plugins/bootstrap-hover-dropdown/bootstrap-hover-dropdown.min.js',

        './assets/global/plugins/js.cookie.min.js',
        './assets/global/plugins/bootstrap-switch/js/bootstrap-switch.min.js',
        './assets/global/plugins/angularjs/angular.min.js',
        './assets/global/plugins/angularjs/i18n/angular-locale_zh-cn.js',
        './assets/global/plugins/angularjs/angular-sanitize.min.js',
        './assets/global/plugins/angularjs/angular-touch.min.js',
        './assets/global/plugins/angularjs/angular-cookies.min.js',
        './assets/global/plugins/angular-ui-router.min.js',
        './assets/global/plugins/angularjs/plugins/ocLazyLoad.min.js',
        'js/lib/angular-ui-bootstrap/ui-bootstrap-tpls-2.5.0.min.js'
        ])
        .pipe(concat('base.min.js'))
        .pipe(uglify())
        .pipe(gulp.dest('js'));
});

gulp.task('jsPlugins', function () {
    gulp.src(['bower_components/moment/moment.js','bower_components/moment/locale/zh-cn.js',
        'bower_components/angular-moment/angular-moment.min.js',
        'bower_components/angular-ui-select/dist/select.js',

        './assets/global/plugins/datatables/datatables.min.js',
        'bower_components/sweetalert2/dist/sweetalert2.min.js',
        'bower_components/angular-bootstrap-switch/dist/angular-bootstrap-switch.min.js',
        'js/lib/angular-datatables/angular-datatables.min.js',
        'js/lib/amcharts/amcharts.js',
        'js/lib/amcharts/serial.js',
        'js/lib/amcharts/pie.js',
        'js/lib/amcharts/themes/light.js',
        'js/lib/amcharts/lang/zh.js'
    ])
        .pipe(concat('plugins.min.js'))
        .pipe(uglify())
        .pipe(gulp.dest('js'));
});

gulp.task('jsAmcharts', function () {
    gulp.src([

    ])
        .pipe(concat('amcharts.min.js'))
        .pipe(uglify())
        .pipe(gulp.dest('js'));
});


gulp.task('jsLayouts', function () {
    gulp.src(['js/app.js','./assets/layouts/boiler/scripts/layout.min.js',
        './assets/layouts/boiler/scripts/demo.min.js'
    ])
        .pipe(concat('layouts.min.js'))
        .pipe(uglify())
        .pipe(gulp.dest('js'));
});

gulp.task('jsLadda', function () {
    gulp.src(['js/scripts/ladda/spin.min.js','js/scripts/ladda/ladda.min.js',
        'js/scripts/ladda/core.js'
    ])
        .pipe(concat('ladda.min.js'))
        .pipe(uglify())
        .pipe(gulp.dest('js'));
});