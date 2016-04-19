const gulp = require("gulp");
const rename = require("gulp-rename");

gulp.task("default", function(){
    gulp.src("node_modules/jquery/dist/jquery.min.js").pipe(rename("jquery.js")).pipe(gulp.dest("../public/js"));
    gulp.src([
        "node_modules/bootstrap/dist/css/bootstrap.css",
        "node_modules/bootstrap/dist/css/bootstrap.css.map",
        "node_modules/bootstrap/dist/css/bootstrap-theme.css",
        "node_modules/bootstrap/dist/css/bootstrap-theme.css.map"
    ]).pipe(gulp.dest("../public/css"));
    gulp.src("node_modules/bootstrap/dist/js/bootstrap.js").pipe(gulp.dest("../public/js"));
    gulp.src("node_modules/bootstrap/dist/fonts/**.*").pipe(gulp.dest("../public/fonts"));
    gulp.src('node_modules/font-awesome/css/**.*' ).pipe(gulp.dest("../public/css"));
    gulp.src('node_modules/font-awesome/fonts/**.*').pipe(gulp.dest("../public/fonts"));
});