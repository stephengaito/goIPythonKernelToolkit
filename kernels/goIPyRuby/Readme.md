# IPyRuby kernel

## Stability 

For some reason (due either to the way I am using the Ruby ANSI-C API, or 
due to Ruby itself), the current IPyRuby kernel using Ruby 2.5.1p57 
(Ubuntu 18.04 bionic) seems to randomly return numerous "stack level too 
deep" errors, but on re-running the same code (with out restarting the 
kernel), behaves as expected. 

