# Bookmarkscli
A bookmark manager tool for the Linux command line.

What this program can do?
- On a first run, the program will ask about the default web browser to open the bookmarks and the desired location for the sqlite3 database with the bookmarks.
  Please provide the full path to the web browser executable and the full path and name for the db.

- When adding a bookmark into the database, the program will check if the url already exist and will warn about it.

- When the url text field contains a url, you can go to the name and tags text fields and press ctrl+u to populate those fields with words from the url.
  Just do the minor corrections to remove/add words as needed.

- In the listbox, you will be able to open a bookmark in your web browser using the intro key.
  Just use the up/down arrow keys to select a line and the ctrl+e to populate the text fields name, url and tags. Pressing the intro key will perform an update.
  Pressing ctrl+d will remove the listbox selected item from the db.

- To change the theme edit the file inside the themes folder.

<img src="bookmarkscli_demo.gif">
