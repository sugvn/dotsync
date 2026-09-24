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
        ~/dotfiles> sudo dotsync pull

important note:
        all the files in the dotfiles directory should be write protected,so you wont modify the
        files by accident and only modify and pull from the original directory entries.
        so dotsync pull need root priviliges to modify.

        This ensures that you dont modify both the file in the dotfiles directory as well as the
        original directory entried files and cause merge conflicts

list of commands:

pull <basename>:
        copy changes from the destination directory to the current directory
        eg: cp -r ~/.config/sway ./sway

        By default copy every destination directory entry from the dirlist.txt to the basename
        directories in the current directory
        > dotsync pull

push <basename>:
        copy changes of the current directory to the destination directory
        eg: cp -r ./sway ~/.config/sway

        By default copy every basename entry to its respective destination directory specified
        in the dirlist.txt
        > dotsync push

commit <msg>:
        commit the changes to git
