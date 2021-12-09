reset 
set style fill solid
set title 'performance comparison'
set ylabel 'ns/op'
set term png enhanced font 'Verdana,10'
set output 'perf.png'

plot 'go-work-final.dat' u 2:xtic(1) w histograms title 'go-work', \
'for-select-for-select-final.dat' u 2:xtic(1) w histograms title 'for-s-for-s', \
'for-select-for-range-final.dat' u 2:xtic(1) w histograms title 'for-s-for-r', \
'no-jobChan-for-range-worker-final.dat' u 2:xtic(1) w histograms title 'no-jobChan-for-r';

