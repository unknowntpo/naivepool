reset 
set style fill solid
set title 'performance comparison' font 'Verdana,10'
set ylabel 'ns/op (log10)' font 'Verdana, 8'
set style histogram clustered gap 1 title offset 0.1,0.25
set term png enhanced font 'Verdana,5'
set output 'perf.png'
set logscale y 10

set multiplot

set title 'simple task'
set size 0.3,0.5
set origin 0.0,0
plot [*:3.5] 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2 w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2 w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2 w histograms title 'no-jobChan-for-r';


set title 'long-running task'
set size 0.3,0.5
set origin 0.3,0
plot [3.6:7.3] 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2 w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2 w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2 w histograms title 'no-jobChan-for-r';


set title 'print'
set size 0.3,0.5
set origin 0.6,0
plot [7.6:11.3] 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2 w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2 w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2 w histograms title 'no-jobChan-for-r';

unset multiplot
