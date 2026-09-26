DOTSYNC - A simple dotfiles management helper

How it works:
        you create a directory in the home directory like,
        ~> mkdir dotfiles
        ~> cd dotfiles
        ~/dotfiles>

        you specify the list of directories you want to track in dirlist.txt
        ~/dotfiles> echo "/home/user/.config/sway" >> dirlist.txt
        ~/dotfiles> echo "/home/user/.config/kitty" >> dirlist.txt
        ~/dotfiles> cat dirlist.txt
        /home/user/.config/sway
        /home/user/.config/kitty

        Then you pull the entries from the dirlist.txt into the current directory
        ~/dotfiles>dotsync -action pull

list of actions:
pull:
        copy every destination directory entry from the dirlist.txt to the basename
        directories in the current directory

push:
        copy every basename entry to its respective destination directory specified
        in the dirlist.txt
        > dotsync -action push

additional flags:
-file:
     specify the file to use as directory listing   
     Default is dirlist.txt
     usage:
        > dotsync -file <filename> -action <action>

installation instructions:
        1. clone the repository
        2. go build -o dotsync main.go
        3. mv dotsync ~/.local/bin/
