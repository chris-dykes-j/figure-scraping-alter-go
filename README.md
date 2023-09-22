# Figure Scraping in Go

Project to scrape data and images from Alter's wonderful website (https://alter-web.jp/).
Measures to ensure reduced traffic to the website are made to the program.

The csv file currently has duplicate entries due to multiple runs, and because items have been released multiple years (the scraper is finds items by year, and does not check if an item is already on the list).

## Known Issues
1. For the time "Fate/EXTRA セイバーエクストラ　水着Ver.", the '/' character had to be replaced with '-' in order to satisfy the directory naming when creating the figure directory, however I noticed that when saving the images, the program went ahead and created necessary directories without complaints.
2. For the 'character' column, lines on separate lines got mushed together and will need to be manually fixed as a result. Not really a big deal given the number of items that had it was small, but the bug should be fixed.
