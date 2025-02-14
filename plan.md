# Get a Fucking Front Page 
## Intro

Basic first step is a front page which gives a table overview of the current crags starting with those that aren't raining, and have the longest weather window, descending with those with the most rain being at the bottom.

## How

I think we're going to generate a html file or string using templ and then maybe post it to the front-end every hour - and this will be how it updates hourly. 
Maybe we could be polling im not sure but we're just going to try and basically generate the html using templ, and then send to static server to be served. 

OR,  we serve is ourselves not sure which is simpler

TODO


[] make the templ component 
[] send post request
[] way to recieve the request on the static server?! 
[] serve the html on thge 6996 server