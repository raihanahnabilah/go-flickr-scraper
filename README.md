# Flickr Scraper in Golang

## Background 

In my part-time job as a Student Associate for Student Affairs Office at <a href="https://www.yale-nus.edu.sg/">Yale-NUS College</a>, I got a task to migrate some of the College's albums from Flickr to Share Point Drive. Flickr has the Download functionality where users can download 500 pictures in one go; however, for albums with pictures more than 500, this Download functionality wouldn't work. I feel challenged to put my Software Engineering skills into use, and I started looking into Flickr's Service APIs and Open Source projects. 

## Running the project



## Acknowledgement

I would like to thank:
- Halvor Haukvik, who created <a href="https://github.com/hdhauk/flickrdump">flickrdump</a> and Suhun Han, who created <a href="https://github.com/ssut/flickr-dump">flickr-dump</a>. This project is largely inspired by their projects, particularly in understanding how Flickr Service APIs works and concurrency in Golang. I could have used their projects, but there are some additional functionalities that I need (like "DownloadByAlbum" where users can download only certain albums instead of all albums). I also feel challenged to write it to make sure I understand what I am doing. 
- Ultralystics, who created <a href=https://github.com/ultralytics/flickr_scraper/tree/master> Flickr Scraper</a>. This project helps me in understanding how the flickrapi package works. 


