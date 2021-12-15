reset 
set style fill solid
set format y "10^ { %L }"
set ylabel 'ns/op (log base 10)' font 'Verdana, 8'
set key left top
set style histogram clustered gap 1 title offset 0.1,0.25
set term png enhanced font 'Verdana,5'
set output 'perf.png'
set logscale y 10

set multiplot layout 1, 3 title 'performance comparison' font 'Verdana,15' offset 0.5,-10

set title 'simple task'
set size 0.3,0.7
set origin 0.0,0
plot [*:3.5] 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2 w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2 w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2 w histograms title 'no-jobChan-for-r', \
'normal-goroutine-final.dat' u 2 w histograms title 'normal-go', \
'pond-final.dat' u 2 w histograms title 'pond';

set title 'long-running task'
set size 0.3,0.7
set origin 0.3,0
plot [3.6:7.5] 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2 w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2 w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2 w histograms title 'no-jobChan-for-r', \
'normal-goroutine-final.dat' u 2 w histograms title 'normal-go', \
'pond-final.dat' u 2 w histograms title 'pond';

set title 'print'
set size 0.3,0.7
set origin 0.6,0
plot [7.55:11.5] 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2 w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2 w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2 w histograms title 'no-jobChan-for-r', \
'normal-goroutine-final.dat' u 2 w histograms title 'normal-go', \
'pond-final.dat' u 2 w histograms title 'pond';

unset multiplot
