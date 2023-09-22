# figure-scraping-alter-go

Project to scrape data and images from Alter's wonderful website (https://alter-web.jp/).
Measures to ensure reduced traffic to the website are made to the program.

The csv file currently has duplicate entries due to multiple runs, and because items have been released multiple years (the scraper is finds items by year, and does not check if an item is already on the list).

## Known bugs
1. For the time "Fate/EXTRA セイバーエクストラ　水着Ver.", the '/' character had to be replaced with '-' in order to satisfy the directory naming when creating the figure directory, however I noticed that when saving the images, the program went ahead and created necessary directories without complaints. 
