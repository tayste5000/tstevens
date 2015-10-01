var gulp = require('gulp');
var $ = require('gulp-load-plugins')();

gulp.task('sass', function(){
	gulp.src('sass/index.scss')
		.pipe($.sass())
		.pipe($.autoprefixer())
		.pipe($.rename('style.css'))
		.pipe(gulp.dest('public/css'));
});

gulp.task('watch', function(){
	gulp.watch('sass/*', ['sass']);
});