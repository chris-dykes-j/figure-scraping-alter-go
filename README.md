# Figure Scraping in Go

Project to scrape data and images from Alter's wonderful website (https://alter-web.jp/).
Measures to ensure reduced traffic to the website are made to the program.

The scraped data csv file is a bit messy, but the file in the transform directory is cleaned of duplicate entries. The spacing issue was manually fixed, but has been fixed in the program.

## Known Issues
1. For the time "Fate/EXTRA セイバーエクストラ　水着Ver.", the '/' character had to be replaced with '-' in order to satisfy the directory naming when creating the figure directory, however I noticed that when saving the images, the program went ahead and created necessary directories without complaints.
